import {
  CallHandler,
  ExecutionContext,
  Inject,
  Injectable,
  NestInterceptor,
  forwardRef,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { Socket } from 'socket.io';
import { PlayerService } from 'src/module/player/player.service';
import { EmitEventFunc } from '../chat.type';
import { EmitEvent } from '../chat.enum';

@Injectable()
export class SocketPlayerIdBindingInterceptor implements NestInterceptor {
  constructor(
    @Inject(forwardRef(() => PlayerService))
    private playerService: PlayerService,
  ) {}

  async intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Promise<Observable<any>> {
    const client = context.switchToWs().getClient() as Socket<
      null,
      EmitEventFunc
    >;
    const playerId = await this.playerService.getIdBySocketId(client.id);

    if (playerId) {
      client.playerId = playerId;

      return next.handle();
    }

    client.emit(EmitEvent.Error, {
      event: client.eventName,
      message: 'Something went wrong. Please try to login again!',
    });

    client.disconnect();
  }
}
