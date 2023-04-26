import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
} from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import * as aes256 from 'aes256';
import { AppConfig } from 'src/config';

@Injectable()
export class ApiKeyGuard implements CanActivate {
  _cipher: any;

  constructor(private appConfig: AppConfig) {
    this._cipher = aes256.createCipher(appConfig.secret);
  }

  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest<FastifyRequest>();
    const apiKey = request.query['apiKey'];

    if (apiKey == undefined) {
      throw new ForbiddenException('Api key is required!');
    }

    try {
      const decryptedApiKey = this._cipher.decrypt(apiKey);

      if (decryptedApiKey !== this.appConfig.secret) {
        throw new ForbiddenException('Api key is invalid!');
      }
    } catch (_) {
      throw new ForbiddenException('Api key is invalid!');
    }

    return true;
  }
}
