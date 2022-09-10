import { UseFilters, UsePipes, ValidationPipe } from '@nestjs/common';
import {
  GatewayMetadata,
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

  constructor(private connectionService: ConnectionService) {}

  async handleConnection(client: Socket) {
    try {
      const userId = await this.connectionService.validateConnection(client);

      await this.connectionService.connect(client, userId);
    } catch (error) {
      client.emit(EmitedEvent.Error, {
        event: ListenedEvent.Connect,
        error: (error as WsException).getError(),
      });

      client.disconnect();
    }
  }

  async handleDisconnect(client: Socket) {
    await this.connectionService.disconnect(client);
  }

  @SubscribeMessage(ListenedEvent.PrivateMessage)
  async handleMessage(client: any, payload: any) {
    console.log(payload);
  }
}
