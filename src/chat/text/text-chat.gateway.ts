import { UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
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
import { EmitedEvent, ListenedEvent } from 'src/enum/event.enum';
import { AllExceptionsFilter } from 'src/filter/all-exceptions.filter';
import { ConnectionService } from '../connection.service';
import { PrismaService } from '../prisma.service';
import { UserService } from '../user.service';
import { SendPrivateMessageDto } from './dto/send-private-message.dto';

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
  private readonly server: Server;

  constructor(
    private userService: UserService,
    private prismaService: PrismaService,
    private connectionService: ConnectionService,
  ) {}

  async handleConnection(client: Socket) {
    try {
      const user = await this.connectionService.validateConnection(
        this.server,
        client,
      );
      await this.userService.connect(user, client.id);
      const sIds = await this.userService.getOnlineFriendSocketIds(user.id);

      this.server.to(sIds).emit(EmitedEvent.FriendStatus, {
        data: {
          id: user.id,
          online: true,
        },
      });
    } catch (error: any) {
      client.emit(EmitedEvent.Error, {
        event: ListenedEvent.Connect,
        error: error.message,
      });

      client.disconnect();
    }
  }

  async handleDisconnect(client: Socket) {
    const user = await this.userService.getBySocketId(client.id);

    if (user != null) {
      await this.userService.disconnect(user, client.id);

      const sIds = await this.userService.getOnlineFriendSocketIds(user.id);

      this.server.to(sIds).emit(EmitedEvent.FriendStatus, {
        data: {
          id: user.id,
          online: false,
        },
      });
    }
  }

  @SubscribeMessage(ListenedEvent.PrivateMessage)
  async handleMessage(
    @ConnectedSocket() client: Socket,
    @MessageBody() payload: SendPrivateMessageDto,
  ) {
    const user = await this.userService.getBySocketId(client.id);

    if (!(await this.userService.areFriends(user.id, payload.receivedId))) {
      throw new WsException('Only friends can send messages to each other!');
    }

    const sids = await this.userService.getSocketIds(payload.receivedId);

    this.server.to(sids as string[]).emit(EmitedEvent.PrivateMessage, {
      data: {
        senderId: user.id,
        ...payload,
      },
    });
  }
}
