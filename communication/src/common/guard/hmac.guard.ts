import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { Hmac, createHmac } from 'crypto';
import { FastifyRequest } from 'fastify';
import { AppConfig } from 'src/config';

const HMAC_ALGORITHM = 'sha256';

@Injectable()
export class HmacGuard implements CanActivate {
  /**
   * HMAC generator.
   */
  private readonly _hmac!: Hmac;

  constructor(config: AppConfig) {
    this._hmac = createHmac(HMAC_ALGORITHM, config.secret);
  }

  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest<FastifyRequest>();
    const hmac = String(request.headers['X-HMAC-Signature']).replace(
      'HMAC ',
      '',
    );

    if (!hmac) {
      throw new UnauthorizedException('HMAC is required!');
    }

    const expectedHmac = this._hmac
      .update(JSON.stringify(request.body))
      .digest('hex');
    if (hmac !== expectedHmac) {
      throw new UnauthorizedException('Request is invalid!');
    }

    return true;
  }
}
