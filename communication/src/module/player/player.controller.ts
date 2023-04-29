import { Controller, Get, Res, UseGuards } from '@nestjs/common';
import { FastifyReply } from 'fastify';
import { PlayerService } from './player.service';
import { TokenGuard } from 'src/common/guard';
import { HttpPlayer } from 'src/common/decorator';
import { Player } from '@prisma/client';

@Controller('players')
@UseGuards(TokenGuard)
export class PlayerController {
  constructor(private playerService: PlayerService) {}

  /**
   * Get friends of the logged-in player.
   *
   * @param request
   * @param response
   */
  @Get('friends')
  async getFriends(
    @HttpPlayer() player: Player,
    @Res() response: FastifyReply,
  ): Promise<void> {
    const friends = await this.playerService.getFriends(player.id);
    response.code(200).send({
      data: friends,
    });
  }
}
