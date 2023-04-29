import {
  BadRequestException,
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException,
} from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import { PlayerService } from 'src/module/player';

@Injectable()
export class RequireActiveGuard implements CanActivate {
  constructor(private readonly playerService: PlayerService) {}

  async canActivate(ctx: ExecutionContext): Promise<boolean> {
    const request = ctx.switchToHttp().getRequest<FastifyRequest>();
    if (!request.player) {
      throw new UnauthorizedException('Unable to detect player!');
    }

    const sid = await this.playerService.getSocketId(request.player.id);
    if (!sid) {
      throw new BadRequestException('You are offline!');
    }

    request.player.sid = sid;
    return true;
  }
}
