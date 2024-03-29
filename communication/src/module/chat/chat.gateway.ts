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
import { AllExceptionFilter, WsExceptionFilter } from 'src/common/filter';
import { EventBindingInterceptor } from './interceptor';
import { SendRoomMessageDto, SendPrivateMessageDto } from './dto';
import { ChatService } from './chat.service';
import { ListenEvent } from './chat.enum';
import { EmitEventMap } from './chat.type';
import { loadConfig } from 'src/utils/load-config';

const appConfig = loadConfig().app;

@Injectable()
@UseFilters(AllExceptionFilter, WsExceptionFilter)
@UsePipes(
  new ValidationPipe({
    whitelist: true,
  }),
)
@WebSocketGateway<GatewayMetadata>({
  namespace: '/',
  cors: {
    origin: appConfig.cors.origins.map((v) => new RegExp(v)),
    methods: ['GET', 'POST'],
  },
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer()
  readonly server!: Server<EmitEventMap>;

  constructor(private chatService: ChatService) {}

  /**
   * Store player state before connection.
   *
   * @param client
   */
  async handleConnection(client: Socket): Promise<void> {
    await this.chatService.connect(this.server, client);
  }

  /**
   * Remove player state after disconnection.
   *
   * @param client
   */
  async handleDisconnect(client: Socket): Promise<void> {
    await this.chatService.disconnect(this.server, client);
  }

  /**
   * Send private message.
   *
   * @param client
   * @param payload
   */
  @UseInterceptors(new EventBindingInterceptor(ListenEvent.PrivateMessage))
  @SubscribeMessage(ListenEvent.PrivateMessage)
  async sendPrivateMessage(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: SendPrivateMessageDto,
  ): Promise<void> {
    await this.chatService.sendPrivateMessage(this.server, client, payload);
  }

  /**
   * Send room message.
   *
   * @param client
   * @param payload
   */
  @UseInterceptors(new EventBindingInterceptor(ListenEvent.RoomMessage))
  @SubscribeMessage(ListenEvent.RoomMessage)
  async sendRoomMesage(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: SendRoomMessageDto,
  ): Promise<void> {
    await this.chatService.sendRoomMessage(this.server, client, payload);
  }
}
