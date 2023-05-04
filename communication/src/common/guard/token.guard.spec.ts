import { AuthService } from 'src/module/auth/auth.service';
import { TokenGuard } from './token.guard';
import { ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { createMock } from '@golevelup/ts-jest';
import { Player } from '@prisma/client';
import { PlayerStatus } from 'src/module/player/player.enum';

describe('TokenGuard', () => {
  const player: Player = {
    id: 'player_01_id',
    statusId: PlayerStatus.Offline,
    username: 'player 01',
  };

  const mockAuthService = createMock<AuthService>({
    getPlayer: jest.fn(),
  });

  const mockSwitchToHttp = {
    getRequest: jest.fn(),
  };
  const mockExecutionContext = createMock<ExecutionContext>({
    switchToHttp: () => mockSwitchToHttp,
  });

  const tests: {
    name: string;
    error?: Error;
    setup: () => void;
  }[] = [
    {
      name: 'Empty token',
      error: new UnauthorizedException('Token is required!'),
      setup: () => {
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            authorization: '',
          },
        });
      },
    },
    {
      name: 'Invalid token',
      error: new UnauthorizedException('Invalid token!'),
      setup: () => {
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            authorization: 'token',
          },
        });
        mockAuthService.getPlayer.mockImplementationOnce(() => {
          throw new UnauthorizedException('Invalid token!');
        });
      },
    },
    {
      name: 'Ok',
      setup: () => {
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            authorization: 'token',
          },
        });
        mockAuthService.getPlayer.mockResolvedValueOnce(player);
      },
    },
  ];

  tests.forEach((t) => {
    test(t.name, async () => {
      t.setup();

      const tokenGuard = new TokenGuard(mockAuthService);

      if (t.error) {
        await expect(
          tokenGuard.canActivate(mockExecutionContext),
        ).rejects.toEqual(t.error);
      } else {
        await expect(
          tokenGuard.canActivate(mockExecutionContext),
        ).resolves.toEqual(true);
      }
    });
  });
});
