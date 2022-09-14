import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { Socket } from 'socket.io';
import { EmitEvent } from 'src/enum/event.enum';
import { UserService } from 'src/module/common/user.service';
import { EmitEvents } from 'src/type/event.type';

@Injectable()
export class SocketUserIdBindingInterceptor implements NestInterceptor {
  constructor(private userService: UserService) {}

  async intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Promise<Observable<any>> {
    const client = context.switchToWs().getClient() as Socket<null, EmitEvents>;
    const userId = await this.userService.getId(client.id);

    if (!(userId > 0)) {
      client.emit(EmitEvent.Error, {
        event: client.eventName,
        message: 'Something went wrong. Please try to login again!',
      });

      client.disconnect();
    }

    client.userId = userId;

    return next.handle();
  }
}
