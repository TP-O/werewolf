import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import { AuthService } from 'src/module/auth/auth.service';

@Injectable()
export class TokenGuard implements CanActivate {
  constructor(private readonly authService: AuthService) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest<FastifyRequest>();
    const token = String(request.headers.authorization).replace('Bearer ', '');
    if (!token) {
      throw new UnauthorizedException('Token is required!');
    }

    const player = await this.authService.getPlayer(token);
    request.player = player;
    return true;
  }
}
