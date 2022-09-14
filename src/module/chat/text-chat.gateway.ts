import {
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
import { ValidationConfig } from 'src/config/validation.config';
import { AllExceptionsFilter } from 'src/filter/all-exceptions.filter';
import { UserService } from 'src/module/common/user.service';
import { ConnectionService } from './connection.service';
import { MessageService } from './message.service';
import { SendPrivateMessageDto } from './dto/send-private-message.dto';
import { SocketUserIdBindingInterceptor } from 'src/interceptor/socket-user-id-binding.interceptor';
import { RoomService } from './room.service';
import { JoinRoomDto } from './dto/join-room.dto';
import { LeaveRoomDto } from './dto/leave-room.dto';
import { EventNameBindingInterceptor } from 'src/interceptor/event-name-binding.interceptor';
import { EmitEvent, ListenEvent } from 'src/enum/event.enum';
import { EmitEvents } from 'src/type/event.type';
import { ActiveStatus } from 'src/enum/user.enum';
import { KickOutOfRoomDto } from './dto/kick-out-of-room.dto';
import { RoomChange } from 'src/enum/room.enum';

@UseFilters(new AllExceptionsFilter())
@UsePipes(new ValidationPipe(ValidationConfig))
@WebSocketGateway<GatewayMetadata>({
  namespace: 'text',
  cors: {
    origin: '*',
    methods: ['GET', 'POST'],
    credentials: true,
  },
})
export class TextChatGateway
  implements OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  private readonly server: Server<null, EmitEvents>;

  constructor(
    private userService: UserService,
    private connectionService: ConnectionService,
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

      const friendsSids = await this.userService.getOnlineFriendsSocketIds(
        user.id,
      );
      friendsSids.forEach((sids) => client.to(sids));
      client.emit(EmitEvent.UpdateFriendStatus, {
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
    const user = await this.userService.getBySocketId(client.id);

    if (user != null) {
      await this.userService.disconnect(user, client.id);

      const friendsSids = await this.userService.getOnlineFriendsSocketIds(
        user.id,
      );

      friendsSids.forEach((sids) => client.to(sids));
      client.emit(EmitEvent.UpdateFriendStatus, {
        id: user.id,
        status: null,
      });
    }
  }

  /**
   * Send private message and notify receiver.
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
    if (
      !(await this.userService.areFriends(client.userId, payload.receiverId))
    ) {
      throw new WsException('Only friends can send messages to each other!');
    }

    await this.messageService.createPrivateMessage(client.userId, payload);
    const sids = await this.userService.getSocketIds(payload.receiverId);

    this.server.to(sids as string[]).emit(EmitEvent.ReceivePrivateMessage, {
      senderId: client.userId,
      ...payload,
    });
  }

  /**
   * Create a new room.
   *
   * @param client socket client.
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.CreateRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.CreateRoom)
  async handleCreateRoom(@ConnectedSocket() client: Socket<null, EmitEvents>) {
    const room = await this.roomService.bookRoom(client.userId);
    client.join(room.id);

    client.emit(EmitEvent.ReceiveRoomChanges, {
      roomId: room.id,
      memberIds: room.memberIds,
      change: {
        type: RoomChange.Join,
        memeberId: client.userId,
      },
    });
  }

  /**
   * Join to a new room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.JoinToRoom),
    SocketUserIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.JoinToRoom)
  async handleJoinRoom(
    @ConnectedSocket() client: Socket<null, EmitEvents>,
    payload: JoinRoomDto,
  ) {
    const room = await this.roomService.joinRoom(payload.id, client.userId);
    client.join(room.id);

    this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      roomId: room.id,
      memberIds: room.memberIds,
      change: {
        type: RoomChange.Join,
        memeberId: client.userId,
      },
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
    payload: LeaveRoomDto,
  ) {
    const room = await this.roomService.leaveRoom(payload.id, client.userId);
    client.leave(room.id);

    client.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      roomId: room.id,
      memberIds: room.memberIds,
      change: {
        type: RoomChange.Leave,
        memeberId: client.userId,
      },
    });
  }

  /**
   * Kick member out or room.
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
    payload: KickOutOfRoomDto,
  ) {
    const room = await this.roomService.kickMember(
      payload.id,
      client.userId,
      payload.memberId,
    );

    this.server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
      roomId: room.id,
      memberIds: room.memberIds,
      change: {
        type: RoomChange.Kick,
        memeberId: payload.memberId,
      },
    });

    client.leave(room.id);
  }
}
