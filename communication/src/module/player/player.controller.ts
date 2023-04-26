import { Controller, Get, Req, Res, UseGuards } from '@nestjs/common';
import { FastifyReply, FastifyRequest } from 'fastify';
import { AuthGuard } from 'src/common/guard';
import { PlayerService } from './player.service';

@Controller('players')
@UseGuards(AuthGuard)
export class PlayerController {
  constructor(private playerService: PlayerService) {}

  /**
   * Get friend list of logged in player.
   *
   * @param request
   * @param response
   */
  @Get('friends')
  async getFriendList(
    @Req() request: FastifyRequest,
    @Res() response: FastifyReply,
  ) {
    const friendList = await this.playerService.getFriendList(
      request.player.id,
    );

    response.code(200).send({
      data: friendList,
    });
  }
}
