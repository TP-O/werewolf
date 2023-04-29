import {
  BadRequestException,
  ForbiddenException,
  Inject,
  Injectable,
  NotFoundException,
  UnauthorizedException,
  forwardRef,
} from '@nestjs/common';
import { Server, Socket } from 'socket.io';
import { SendPrivateMessageDto, SendRoomMessageDto } from './dto';
import { AuthService } from '../auth';
import { PlayerId, PlayerService, PlayerStatus, SocketId } from '../player';
import { RoomService } from '../room';
import { EmitEvent, ListenEvent, RoomChangeType } from './chat.enum';
import { EmitEventMap } from './chat.type';
import { LoggerService } from '../common';
import { Player } from '@prisma/client';

@Injectable()
export class ChatService {
  constructor(
    private authService: AuthService,
    @Inject(forwardRef(() => PlayerService))
    private playerService: PlayerService,
    @Inject(forwardRef(() => RoomService))
    private roomService: RoomService,
    private logger: LoggerService,
  ) {
    this.logger.setContext(ChatService.name);
  }

  /**
   * Connect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async connect(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: PlayerId }>,
  ): Promise<void> {
    try {
      const player = await this._validateConnection(server, client);
      await this.playerService.connect(player.id, client.id);
      client.data.id = player.id;

      // Notify online friends
      const friendSids = (
        await this.playerService.getFriendsSocketIds(player.id)
      ).filter((sid) => !!sid) as SocketId[];
      server.to(friendSids).emit(EmitEvent.FriendStatus, {
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
   * Check if the connection satisfies some sepecific conditions
   * before allowing the connection.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  private async _validateConnection(
    server: Server<EmitEventMap>,
    client: Socket,
  ): Promise<Player> {
    const player = await this._validateAuthorization(
      client.handshake.headers.authorization ?? '',
    );

    const sid = await this.playerService.getSocketId(player.id);
    if (sid) {
      server.to(sid).emit(EmitEvent.Error, {
        event: ListenEvent.Connect,
        message: 'This account is being connected by someone else!',
      });
      server.to(sid).disconnectSockets();
      this.disconnect(server, client);
    }

    return player;
  }

  /**
   * Verify token.
   *
   * @param headerAuthorization
   * @returns player record.
   */
  private async _validateAuthorization(
    headerAuthorization: string,
  ): Promise<Player> {
    const token = String(headerAuthorization).replace('Bearer ', '');
    if (!token) {
      throw new UnauthorizedException('Invalid token!');
    }

    const player = await this.authService.getPlayer(token);
    return player;
  }

  /**
   * Disconnect the client.
   *
   * @param server websocket server.
   * @param client socket client.
   */
  async disconnect(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: PlayerId }>,
  ): Promise<void> {
    try {
      if (client.data.id) {
        // Notify online friends
        const friendSids = (
          await this.playerService.getFriendsSocketIds(client.data.id)
        ).filter((sid) => !!sid) as SocketId[];
        server.to(friendSids).emit(EmitEvent.FriendStatus, {
          id: client.data.id,
          status: PlayerStatus.Online,
        });

        const rooms = await this.roomService.removeFromRooms(
          [], // Leave all rooms
          client.data.id,
        );
        rooms.forEach((room) => {
          server.to(room.id).emit(EmitEvent.RoomChange, {
            changeType: RoomChangeType.Leave,
            room: {
              id: room.id,
              memberIds: room.memberIds,
            },
          });
          client.leave(room.id);
        });
      }
    } catch (error) {
      this.logger.error(
        `[SID: ${client.id} - ID: ${client.data.id}] Disconnect failed`,
      );
    }
  }

  /**
   * Send a private message to friend.
   *
   * @param server
   * @param client
   * @param payload
   */
  async sendPrivateMessage(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: PlayerId }>,
    payload: SendPrivateMessageDto,
  ): Promise<void> {
    if (!client.data.id) {
      return;
    }

    const sid = await this.playerService.getSocketId(payload.receiverId);
    if (!sid) {
      throw new BadRequestException('This player is offline!');
    }

    server.to(sid).emit(EmitEvent.PrivateMessage, {
      ...payload,
      senderId: client.data.id,
    });
  }

  /**
   * Send a message to joined room.
   *
   * @param server
   * @param client
   * @param payload
   */
  async sendRoomMessage(
    server: Server<EmitEventMap>,
    client: Socket<EmitEventMap, EmitEventMap, EmitEventMap, { id: PlayerId }>,
    payload: SendRoomMessageDto,
  ): Promise<void> {
    if (!client.data.id) {
      return;
    }

    const room = await this.roomService.get(payload.roomId);
    if (!room) {
      throw new NotFoundException("Room doesn't exist!");
    }

    if (!room.memberIds.includes(client.data.id)) {
      throw new ForbiddenException('You are not in this room!');
    }

    if (room.isMuted) {
      throw new ForbiddenException('Unable to chat at this time!');
    }

    server.to(payload.roomId).emit(EmitEvent.RoomMessage, {
      ...payload,
      senderId: client.data.id,
    });
  }
}
