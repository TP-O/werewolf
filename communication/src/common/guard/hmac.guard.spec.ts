import { createMock } from '@golevelup/ts-jest';
import { ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { HmacGuard } from './hmac.guard';
import { AppConfig } from 'src/config';
import { createHmac } from 'crypto';

describe('HmacGuard', () => {
  const secret = 'secret_xxx';

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
      name: 'Empty HMAC',
      error: new UnauthorizedException('HMAC is required!'),
      setup: () => {
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            'X-HMAC-Signature': '',
          },
        });
      },
    },
    {
      name: 'Invalid request',
      error: new UnauthorizedException('Request is invalid!'),
      setup: () => {
        const body = {
          foo: 'baz',
          bar: 0,
        };
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            'X-HMAC-Signature': 'HMAC invalid hmac',
          },
          body,
        });
      },
    },
    {
      name: 'Ok',
      setup: () => {
        const body = {
          foo: 'baz',
          bar: 0,
        };
        const hmac = createHmac('sha256', secret)
          .update(JSON.stringify(body))
          .digest('hex');
        mockExecutionContext.switchToHttp().getRequest.mockReturnValueOnce({
          headers: {
            'X-HMAC-Signature': `HMAC ${hmac}`,
          },
          body,
        });
      },
    },
  ];

  tests.forEach((t) => {
    test(t.name, async () => {
      t.setup();

      const tokenGuard = new HmacGuard({ secret } as AppConfig);

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
