import { Controller, Get, Req, Res, UseGuards } from '@nestjs/common';
import { FastifyReply, FastifyRequest } from 'fastify';
import { PlayerService } from './player.service';
import { TokenGuard } from 'src/common/guard';

@Controller('players')
@UseGuards(TokenGuard)
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
