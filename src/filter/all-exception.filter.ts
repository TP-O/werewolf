import { ExceptionFilter, Catch, ArgumentsHost } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import { Socket } from 'socket.io';
import { AppConfig } from 'src/config';
import { EmitEvent } from 'src/enum';
import { EmitEvents } from 'src/type';

@Catch()
export class AllExceptionFilter implements ExceptionFilter {
  catch(exception: Error, host: ArgumentsHost) {
    if (AppConfig.debug) {
      throw exception;
    }

    switch (host.getType()) {
      case 'ws':
        this.handleWsException(exception, host);
        break;

      case 'http':
        this.handleHttpException(exception, host);
        break;

      case 'rpc':
        break;
    }
  }

  private handleWsException(exception: Error, host: ArgumentsHost) {
    const client = host.switchToWs().getClient() as Socket<null, EmitEvents>;

    client.emit(EmitEvent.Error, {
      event: client.eventName,
      message: 'Unknown error!',
    });
  }

  private handleHttpException(exception: Error, host: ArgumentsHost) {
    const response = host.switchToHttp().getResponse<FastifyReply>();

    response.code(500).send({
      statusCode: 500,
      message: 'Unknown error!',
    });
  }
}
