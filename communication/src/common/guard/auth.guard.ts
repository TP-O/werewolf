import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import { AuthService } from 'src/module/auth';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private authService: AuthService) {}

  async canActivate(context: ExecutionContext) {
    const request = context.switchToHttp().getRequest<FastifyRequest>();
    const token = String(request.headers.authorization).replace('Bearer ', '');
    const player = await this.authService.getPlayer(token);
    request.player = player;
    return true;
  }
}
