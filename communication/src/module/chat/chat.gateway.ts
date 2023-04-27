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
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { AllExceptionFilter } from 'src/common/filter';
import {
  EventNameBindingInterceptor,
  SocketPlayerIdBindingInterceptor,
} from './interceptor';
import {
  BookRoomDto,
  InviteToRoomDto,
  JoinRoomDto,
  KickOutOfRoomDto,
  LeaveRoomDto,
  RespondRoomInvitationDto,
  TransferOwnershipDto,
} from '../room/dto';
import { SendRoomMessageDto, SendPrivateMessageDto } from './dto';
import { ChatExceptionFilter } from './chat.filter';
import { EmitEventFunc } from './chat.type';
import { ChatService } from './chat.service';
import { ListenEvent } from './chat.enum';

@Injectable()
@UseFilters(AllExceptionFilter, ChatExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
    stopAtFirstError: false,
    transform: true,
  }),
)
@WebSocketGateway<GatewayMetadata>({
  namespace: '/',
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer()
  readonly server: Server<null, EmitEventFunc>;

  constructor(private chatService: ChatService) {}

  /**
   * Store player state before connection.
   *
   * @param client socket client.
   */
  async handleConnection(client: Socket) {
    await this.chatService.connect(this.server, client);
  }

  /**
   * Remove player state after disconnection.
   *
   * @param client socket client.
   */
  async handleDisconnect(client: Socket) {
    await this.chatService.disconnect(this.server, client);
  }

  /**
   * Send private message.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.SendPrivateMessage),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.SendPrivateMessage)
  async sendPrivateMessage(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: SendPrivateMessageDto,
  ) {
    await this.chatService.sendPrivateMessage(this.server, client, payload);
  }

  /**
   * Send room message.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.SendRoomMessage),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.SendRoomMessage)
  async handleSendRoomMesage(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: SendRoomMessageDto,
  ) {
    await this.chatService.sendRoomMessage(this.server, client, payload);
  }

  /**
   * Book a new room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.BookRoom),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.BookRoom)
  async handleBookRoom(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: BookRoomDto,
  ) {
    await this.chatService.createTemporaryRoom(client, payload);
  }

  /**
   * Join to a new room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.JoinRoom),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.JoinRoom)
  async handleJoinRoom(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: JoinRoomDto,
  ) {
    await this.chatService.joinRoom(this.server, client, payload);
  }

  /**
   * Leave the room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.LeaveRoom),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.LeaveRoom)
  async handleLeaveRoom(
    @ConnectedSocket() client: Socket<null, EmitEventFunc>,
    @MessageBody() payload: LeaveRoomDto,
  ) {
    await this.chatService.leaveRoom(this.server, client, payload);
  }

  /**
   * Kick member out of room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.KickOutOfRoom),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.KickOutOfRoom)
  async handleKickOutOfRoom(
    @ConnectedSocket() client: Socket<null, EmitEventFunc>,
    @MessageBody() payload: KickOutOfRoomDto,
  ) {
    await this.chatService.kickOutOfRoom(this.server, client, payload);
  }

  /**
   * Transfer ownership to another member in room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.TranserOwnership),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.TranserOwnership)
  async handleTransferOwnership(
    @ConnectedSocket() client: Socket<null, EmitEventFunc>,
    @MessageBody() payload: TransferOwnershipDto,
  ) {
    await this.chatService.transferRoomOwnership(this.server, client, payload);
  }

  /**
   * Invite a guest into room.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.InviteToRoom),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.InviteToRoom)
  async handleInviteToRoom(
    @ConnectedSocket() client: Socket<null, EmitEventFunc>,
    @MessageBody() payload: InviteToRoomDto,
  ) {
    await this.chatService.inviteToRoom(this.server, client, payload);
  }

  /**
   * Respond to room invitation.
   *
   * @param client socket client.
   * @param payload
   */
  @UseInterceptors(
    new EventNameBindingInterceptor(ListenEvent.RespondRoomInvitation),
    SocketPlayerIdBindingInterceptor,
  )
  @SubscribeMessage(ListenEvent.RespondRoomInvitation)
  async handleRespondInvitation(
    @ConnectedSocket() client: Socket<null, EmitEventFunc>,
    @MessageBody() payload: RespondRoomInvitationDto,
  ) {
    await this.chatService.respondRoomInvitation(this.server, client, payload);
  }
}
