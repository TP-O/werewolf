import {
  Injectable,
  UseFilters,
  UseInterceptors,
  UsePipes,
  ValidationPipe,
} from '@nestjs/common';
import {
  ConnectedSocket,
  GatewayMetadata,
  MessageBody,
  OnGatewayConnection,
  OnGatewayDisconnect,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
  WsException,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { CORSConfig, ValidationConfig } from 'src/config';
import { ActiveStatus, EmitEvent, ListenEvent, RoomEvent } from 'src/enum';
import { AllExceptionFilter, WsExceptionsFilter } from 'src/common/filter';
import {
  EventNameBindingInterceptor,
  SocketUserIdBindingInterceptor,
} from 'src/common/interceptor';
import { EmitEvents } from 'src/type';
import {
  BookRoomDto,
  InviteToRoomDto,
  JoinRoomDto,
  KickOutOfRoomDto,
  LeaveRoomDto,
  RespondRoomInvitationDto,
  TransferOwnershipDto,
} from '../room/dto';
import { UserService } from '../user/user.service';
import { CommunicationService } from './communication.service';
import { MessageService } from '../message/message.service';
import { RoomService } from '../room/room.service';
import { SendRoomMessageDto, SendPrivateMessageDto } from '../message/dto';

@Injectable()
@UseFilters(new AllExceptionFilter(), new WsExceptionsFilter())
@UsePipes(new ValidationPipe(ValidationConfig))
@WebSocketGateway<GatewayMetadata>({
  namespace: '/',
  cors: CORSConfig,
})
export class CommunicationGateway
  implements OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  readonly server: Server<null, EmitEvents>;

  constructor(
    private userService: UserService,
    private connectionService: CommunicationService,
    private messageService: MessageService,
    private roomService: RoomService,
  ) {}

  /**
   * Store user state before connection.
   *
   * @param client socket client.
   */
  async handleConnection(client: Socket<null, EmitEvents>) {
    try {
      const user = await this.connectionService.validateConnection(
        this.server,
        client,
      );
      await this.userService.connect(user, client.id);

      const friendSIds = await this.userService.getOnlineFriendsSocketIds(
        user.id,
      );
      this.server.to(friendSIds).emit(EmitEvent.UpdateFriendStatus, {
        id: user.id,
        status: ActiveStatus.Online,
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
   * Remove user state after disconnection.
   *
   * @param client socket client.
   */
  async handleDisconnect(client: Socket<null, EmitEvents>) {
    try {
      const user = await this.userService.getBySocketId(client.id);

      if (user != null) {
        const friendSIds = await this.userService.getOnlineFriendsSocketIds(
          user.id,
        );
        const { leftRooms } = await this.userService.disconnect(user);

        this.server.to(friendSIds).emit(EmitEvent.UpdateFriendStatus, {
          id: user.id,
          status: ActiveStatus.Offline,
        });

        leftRooms.forEach((room) => {
          if (room.memberIds.length > 0) {
            this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
              event: RoomEvent.Leave,
              actorId: client.userId,
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
   * Send private message.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.SendPrivateMessage),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.SendPrivateMessage)
  async sendPrivateMessage(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: SendPrivateMessageDto,
  ) {
    await this.messageService.createPrivateMessage(client.userId, payload);
    const receiverSId = await this.userService.getSocketIdByUserId(
      payload.receiverId,
    );

    if (receiverSId != null) {
      this.server.to(receiverSId).emit(EmitEvent.ReceivePrivateMessage, {
        ...payload,
        senderId: client.userId,
      });
    }
  }

  /**
   * Create a new room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.BookRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.BookRoom)
  async handleCreateRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: BookRoomDto,
  ) {
    const room = await this.roomService.book(client.userId, payload.isPublic);
    client.join(room.id);

    client.emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Create,
      actorId: client.userId,
      room,
    });
  }

  /**
   * Join to a new room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.JoinRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.JoinRoom)
  async handleJoinRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: JoinRoomDto,
  ) {
    const room = await this.roomService.join(client.userId, payload.roomId);
    client.join(room.id);

    this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Join,
      actorId: client.userId,
      room,
    });
  }

  /**
   * Leave the room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.LeaveRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.LeaveRoom)
  async handleLeaveRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: LeaveRoomDto,
  ) {
    const room = await this.roomService.leave(client.userId, payload.roomId);
    client.leave(room.id);

    this.server.to(client.id).to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Leave,
      actorId: client.userId,
      room,
    });
  }

  /**
   * Kick member out of room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.KickOutOfRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.KickOutOfRoom)
  async handleKickOutOfRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: KickOutOfRoomDto,
  ) {
    const { room, kickedMemberSocketId } = await this.roomService.kick(
      client.userId,
      payload.memberId,
      payload.roomId,
    );

    this.server
      .to(kickedMemberSocketId)
      .to(room.id)
      .emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Kick,
        actorId: client.userId,
        room,
      });

    this.server.to(kickedMemberSocketId).socketsLeave(room.id);
  }

  /**
   * Send room message.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.SendRoomMessage),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.SendRoomMessage)
  async handleSendRoomMesage(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: SendRoomMessageDto,
  ) {
    const room = await this.roomService.get(payload.roomId);

    if (!room.memberIds.includes(client.userId)) {
      throw new WsException('You are not in this room!');
    }

    if (room.isPersistent) {
      await this.messageService.createRoomMessage(client.userId, payload);
    }

    this.server.to(payload.roomId).emit(EmitEvent.ReceiveRoomMessage, {
      ...payload,
      senderId: client.userId,
    });
  }

  /**
   * Transfer ownership to another member in room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.TranserOwnership),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.TranserOwnership)
  async handleTransferOwnership(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: TransferOwnershipDto,
  ) {
    const room = await this.roomService.transferOwnership(
      client.userId,
      payload.candidateId,
      payload.roomId,
    );

    this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Owner,
      actorId: client.userId,
      room,
    });
  }

  /**
   * Invite a guest into room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.InviteToRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.InviteToRoom)
  async handleInviteToRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: InviteToRoomDto,
  ) {
    const { room, guestSocketId } = await this.roomService.invite(
      client.userId,
      payload.guestId,
      payload.roomId,
    );

    this.server.to(guestSocketId).emit(EmitEvent.ReceiveRoomInvitation, {
      roomId: room.id,
      inviterId: client.userId,
    });

    this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      event: RoomEvent.Invite,
      actorId: client.userId,
      room,
    });
  }

  /**
   * Respond to room invitation.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.RespondRoomInvitation),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.RespondRoomInvitation)
  async handleRespondInvitation(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    @MessageBody() payload: RespondRoomInvitationDto,
  ) {
    const { room, leftRooms } = await this.roomService.respondInvitation(
      client.userId,
      payload.isAccpeted,
      payload.roomId,
    );

    if (payload.isAccpeted) {
      client.join(room.id);

      this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
        event: RoomEvent.Join,
        actorId: client.userId,
        room,
      });

      leftRooms.forEach((room) => {
        if (room.memberIds.length > 0) {
          this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
            event: RoomEvent.Leave,
            actorId: client.userId,
            room,
          });
        }
      });
    }
  }
}
