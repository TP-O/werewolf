import { UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
import {
  GatewayMetadata,
  OnGatewayConnection,
  OnGatewayDisconnect,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import Redis from 'ioredis';
import { Server, Socket } from 'socket.io';
import { ValidationConfig } from 'src/config/validation.config';
import { RedisClient } from 'src/decorator/redis.decorator';
import { EmitedEvent, ListenedEvent } from 'src/enum/event.enum';
import { AllExceptionsFilter } from 'src/filter/all-exceptions.filter';
import { ConnectionService } from '../connection.service';
import { PrismaService } from '../prisma.service';
import { UserService } from '../user.service';

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

  @RedisClient()
  private readonly redis: Redis;

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
  async handleMessage(client: any, payload: any) {
    console.log(payload);
  }
}
