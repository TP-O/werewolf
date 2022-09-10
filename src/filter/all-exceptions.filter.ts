import { Catch, ArgumentsHost, BadRequestException } from '@nestjs/common';
import { BaseWsExceptionFilter } from '@nestjs/websockets';
import { EmitedEvent } from 'src/enum/event.enum';

@Catch()
export class AllExceptionsFilter extends BaseWsExceptionFilter {
  catch(exception: Error, host: ArgumentsHost) {
    const client = host.switchToWs().getClient();
    let errorResponse: string | object;

    if (exception instanceof BadRequestException) {
      const res = exception.getResponse();
      errorResponse =
        typeof res === 'string' ? res : (res as any)?.message?.[0];
    } else {
      errorResponse = exception.message;
    }

    client.emit(EmitedEvent.Error, {
      event: null,
      error: errorResponse,
    });
  }
}
