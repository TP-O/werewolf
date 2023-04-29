import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { Socket } from 'socket.io';
import { ListenEvent } from '../chat.enum';

@Injectable()
export class EventBindingInterceptor implements NestInterceptor {
  constructor(private readonly event: ListenEvent) {}

  async intercept(
    context: ExecutionContext,
    next: CallHandler,
  ): Promise<Observable<any>> {
    const client = context.switchToWs().getClient<Socket>();
    client.event = this.event;
    return next.handle();
  }
}
