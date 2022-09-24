import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { Socket } from 'socket.io';
import { EmitEvent } from 'src/enum';
import { UserService } from 'src/service/user.service';
import { EmitEvents } from 'src/type';

@Injectable()
export class SocketUserIdBindingInterceptor implements NestInterceptor {
  constructor(private userService: UserService) {}

  async intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Promise<Observable<any>> {
    const client = context.switchToWs().getClient() as Socket<null, EmitEvents>;
    const userId = await this.userService.getId(client.id);

    if (userId < 1) {
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
