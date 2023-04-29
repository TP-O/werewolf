import {
  createParamDecorator,
  ExecutionContext,
  UnauthorizedException,
} from '@nestjs/common';
import { Player } from '@prisma/client';
import { FastifyRequest } from 'fastify';

export const HttpPlayer = createParamDecorator<ExecutionContext, Player>(
  (ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest<FastifyRequest>();
    if (!request.player) {
      throw new UnauthorizedException('Unable to detect player!');
    }

    return request.player;
  },
);
