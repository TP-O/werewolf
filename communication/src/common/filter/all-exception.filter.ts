import { ExceptionFilter, Catch, ArgumentsHost } from '@nestjs/common';
import { FastifyReply, FastifyRequest } from 'fastify';
import { Socket } from 'socket.io';
import { LoggedError } from '../type';
import { EmitEvent, EmitEventFunc } from 'src/module/chat';
import { LoggerService } from 'src/module/common';

/**
 * Filter all unexpected exceptions.
 */
@Catch()
export class AllExceptionFilter implements ExceptionFilter {
  constructor(private readonly logger: LoggerService) {
    this.logger.setContext(AllExceptionFilter.name);
  }

  catch(exception: Error, host: ArgumentsHost): void {
    let loggedErr: LoggedError;
    switch (host.getType()) {
      case 'ws':
        loggedErr = this._handleWsException(exception, host);
        break;

      case 'http':
        loggedErr = this._handleHttpException(exception, host);
        break;

      default:
        break;
    }

    if (loggedErr) {
      this.logger.error(loggedErr, exception.stack);
    }
  }

  private _handleWsException(
    exception: Error,
    host: ArgumentsHost,
  ): LoggedError {
    const client = host.switchToWs().getClient() as Socket<null, EmitEventFunc>;
    const loggedError: LoggedError = {
      name: exception.name,
      message: exception.message,
      hostType: 'ws',
      event: client.eventName,
      payload: host.switchToWs().getData(),
    };

    client.emit(EmitEvent.Error, {
      event: client.eventName,
      message: 'Unknown error!',
    });

    return loggedError;
  }

  private _handleHttpException(
    exception: Error,
    host: ArgumentsHost,
  ): LoggedError {
    const response = host.switchToHttp().getResponse<FastifyReply>();
    const request = host.switchToHttp().getRequest<FastifyRequest>();
    const loggedError: LoggedError = {
      name: exception.name,
      message: exception.message,
      hostType: 'http',
      url: request.url,
      payload: request.body,
    };

    response.code(500).send({
      statusCode: 500,
      message: 'Unknown error!',
    });

    return loggedError;
  }
}
