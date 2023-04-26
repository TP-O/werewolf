import {
  ForbiddenException,
  Inject,
  Injectable,
  forwardRef,
} from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { SendPrivateMessageDto, SendRoomMessageDto } from './dto';
import {
  BookRoomDto,
  InviteToRoomDto,
  JoinRoomDto,
  KickOutOfRoomDto,
  LeaveRoomDto,
  RespondRoomInvitationDto,
  TransferOwnershipDto,
} from '../room/dto';
import { Player } from '@prisma/client';
import { AuthService } from '../auth';
import { PlayerService, PlayerStatus } from '../player';
import { RoomService } from '../room';
import { EmitEventFunc } from './chat.type';
import { EmitEvent, ListenEvent, RoomEvent } from './chat.enum';

@Injectable()
export class ChatService {
  constructor(
    private authService: AuthService,
    @Inject(forwardRef(() => PlayerService))
    private playerService: PlayerService,
    @Inject(forwardRef(() => RoomService))
    private roomService: RoomService,
  ) {}

  /**
   * Verify token.
   *
   * @param headerAuthorization
   * @returns player record.
   */
  private async validateAuthorization(headerAuthorization: string) {
    const token = String(headerAuthorization).replace('Bearer ', '');
    const player = await this.authService.getPlayer(token);

    return player;
  }

  /**
   * Solve conflict if multiple people connect to the
   * same account.
   *
   * @param server websocket server.
   * @param player
   */
  private async handleConflict(
    server: Server<null, EmitEventFunc>,
    player: Player,
  ) {
    const { disconnectedId, leftRooms } = await this.playerService.disconnect(
      player,
    );

    server.to(disconnectedId).emit(EmitEvent.Error, {
      event: ListenEvent.Connect,
      message: 'This account is being connected by someone else!',
    });
    server.to(disconnectedId).disconnectSockets();

    leftRooms.forEach((room) => {
      if (room.memberIds.length > 0) {
        server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
          event: RoomEvent.Leave,
          actorIds: [player.id],
          room,
        });
      }
    });
  }

  /**
   * Check if the connection satisfies some sepecific conditions
   * before allowing the connection.
   *
   * @param server websocket server.
   * @param client socket client.
   * @returns updated player.
   */
  private async validateConnection(server: Server, client: Socket) {
    const player = await this.validateAuthorization(
      client.handshake.headers.authorization,
    );

    if (player.statusId != null) {
      await this.handleConflict(server, player);
    }

    return player;
  }

  /**
   * Connect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async connect(server: Server<null, EmitEventFunc>, client: Socket) {
    try {
      const player = await this.validateConnection(server, client);
      await this.playerService.connect(player, client.id);

      const friendSIds = await this.playerService.getOnlineFriendsSocketIds(
        player.id,
      );
      server.to(friendSIds).emit(EmitEvent.UpdateFriendStatus, {
        id: player.id,
        status: PlayerStatus.Online,
      });
    } catch (error: any) {
      client.emit(EmitEvent.Error, {
        event: ListenEvent.Connect,
        message: error.message,
      });
      client.disconnect();
    }
  }

  /**
   * Disconnect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async disconnect(server: Server<null, EmitEventFunc>, client: Socket) {
    try {
      const player = await this.playerService.getBySocketId(client.id);

      if (player != null) {
        const friendSIds = await this.playerService.getOnlineFriendsSocketIds(
          player.id,
        );
        const { leftRooms } = await this.playerService.disconnect(player);

        server.to(friendSIds).emit(EmitEvent.UpdateFriendStatus, {
          id: player.id,
          status: PlayerStatus.Offline,
        });

        leftRooms.forEach((room) => {
          if (room.memberIds.length > 0) {
            server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
              event: RoomEvent.Leave,
              actorIds: [client.playerId],
              room,
            });
          }
        });
      }
    } catch (error) {
      //
    }
  }

  /**
   * Send a private message to friend.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async sendPrivateMessage(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: SendPrivateMessageDto,
  ) {
    const receiverSId = await this.playerService.getSocketIdByplayerId(
      payload.receiverId,
    );

    if (receiverSId != null) {
      server.to(receiverSId).emit(EmitEvent.ReceivePrivateMessage, {
        ...payload,
        senderId: client.playerId,
      });
    }
  }

  /**
   * Send a message to joined room.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async sendRoomMessage(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: SendRoomMessageDto,
  ) {
    const room = await this.roomService.get(payload.roomId);

    if (room.isMuted) {
      throw new ForbiddenException('Unable to chat at this time!');
    }

    if (!room.memberIds.includes(client.playerId)) {
      throw new ForbiddenException('You are not in this room!');
    }

    server.to(payload.roomId).emit(EmitEvent.ReceiveRoomMessage, {
      ...payload,
      senderId: client.playerId,
    });
  }

  /**
   * Create a room that is deleted after all members leave.
   *
   * @param client socket client.
   * @param payload
   */
  async createTemporaryRoom(client: Socket, payload: BookRoomDto) {
    const room = await this.roomService.book(client.playerId, payload.isPublic);
    client.join(room.id);

    client.emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Create,
      actorIds: [client.playerId],
      room,
    });
  }

  /**
   * Join to new room.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async joinRoom(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: JoinRoomDto,
  ) {
    const room = await this.roomService.join(client.playerId, payload.roomId);
    client.join(room.id);

    server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Join,
      actorIds: [client.playerId],
      room,
    });
  }

  /**
   * Leave the joined room.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async leaveRoom(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: LeaveRoomDto,
  ) {
    const room = await this.roomService.leave(client.playerId, payload.roomId);
    client.leave(room.id);

    server
      .to(client.id)
      .to(room.id)
      .emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Leave,
        actorIds: [client.playerId],
        room,
      });
  }

  /**
   * Kick a member out of the room.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async kickOutOfRoom(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: KickOutOfRoomDto,
  ) {
    const { room, kickedMemberSocketId } = await this.roomService.kick(
      client.playerId,
      payload.memberId,
      payload.roomId,
    );

    server
      .to(kickedMemberSocketId)
      .to(room.id)
      .emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Kick,
        actorIds: [client.playerId],
        room,
      });

    server.to(kickedMemberSocketId).socketsLeave(room.id);
  }

  /**
   * Transfer room ownership to a member in that room.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async transferRoomOwnership(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: TransferOwnershipDto,
  ) {
    const room = await this.roomService.transferOwnership(
      client.playerId,
      payload.candidateId,
      payload.roomId,
    );

    server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Owner,
      actorIds: [client.playerId],
      room,
    });
  }

  /**
   * Send a room invitation to player.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async inviteToRoom(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: InviteToRoomDto,
  ) {
    const { room, guestSocketId } = await this.roomService.invite(
      client.playerId,
      payload.guestId,
      payload.roomId,
    );

    server.to(guestSocketId).emit(EmitEvent.ReceiveRoomInvitation, {
      roomId: room.id,
      inviterId: client.playerId,
    });
    server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Invite,
      actorIds: [client.playerId],
      room,
    });
  }

  /**
   * Respond to a room invitation.
   *
   * @param server websocket server.
   * @param client socket client.
   * @param payload
   */
  async respondRoomInvitation(
    server: Server<null, EmitEventFunc>,
    client: Socket,
    payload: RespondRoomInvitationDto,
  ) {
    const { room, leftRooms } = await this.roomService.respondInvitation(
      client.playerId,
      payload.isAccpeted,
      payload.roomId,
    );

    if (payload.isAccpeted) {
      client.join(room.id);

      server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Join,
        actorIds: [client.playerId],
        room,
      });

      leftRooms.forEach((room) => {
        if (room.memberIds.length > 0) {
          server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
            event: RoomEvent.Leave,
            actorIds: [client.playerId],
            room,
          });
        }
      });
    }
  }
}
