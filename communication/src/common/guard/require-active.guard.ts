import {
  BadRequestException,
  CanActivate,
  ExecutionContext,
  Inject,
  Injectable,
  UnauthorizedException,
  forwardRef,
} from '@nestjs/common';
import { FastifyRequest } from 'fastify';
import { OnlinePlayer, PlayerService } from 'src/module/player';

@Injectable()
export class RequireActiveGuard implements CanActivate {
  constructor(
    @Inject(forwardRef(() => PlayerService))
    private readonly playerService: PlayerService,
  ) {}

  async canActivate(ctx: ExecutionContext): Promise<boolean> {
    const request = ctx.switchToHttp().getRequest<FastifyRequest>();
    if (!request.player) {
      throw new UnauthorizedException('Unable to detect player!');
    }

    const sid = await this.playerService.getSocketId(request.player.id);
    if (!sid) {
      throw new BadRequestException('You are offline!');
    }

    (request.player as OnlinePlayer).sid = sid;
    return true;
  }
}
