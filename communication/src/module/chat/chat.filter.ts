import {
  Catch,
  ArgumentsHost,
  BadRequestException,
  HttpException,
} from '@nestjs/common';
import { BaseWsExceptionFilter, WsException } from '@nestjs/websockets';
import { Socket } from 'socket.io';
import { EmitEventFunc } from './chat.type';
import { ErrorMessage } from 'src/common/type';
import { EmitEvent } from './chat.enum';

@Catch(WsException, HttpException)
export class ChatExceptionFilter extends BaseWsExceptionFilter {
  catch(exception: Error, host: ArgumentsHost) {
    const client = host.switchToWs().getClient() as Socket<null, EmitEventFunc>;
    let message: ErrorMessage;

    if (exception instanceof BadRequestException) {
      message = (exception.getResponse() as Error).message;
    } else {
      message = exception.message;
    }

    client.emit(EmitEvent.Error, {
      event: client.eventName,
      message: message,
    });
  }
}
