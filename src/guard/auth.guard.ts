import {
  CanActivate,
  ExecutionContext,
  Injectable,
  InternalServerErrorException,
  UnauthorizedException,
} from '@nestjs/common';
import { UserId } from 'src/enum/user.enum';
import { AuthService } from 'src/module/common/auth.service';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private authService: AuthService) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest();
    const token = String(request.headers.authorization).replace('Bearer ', '');
    const user = await this.authService.getUser(token);

    if (user.id === UserId.NonExist) {
      throw new UnauthorizedException('Invalid access token!');
    }

    if (user.id === UserId.Asynchronous) {
      throw new InternalServerErrorException(
        'Please connect again after a while!',
      );
    }

    request.user = user;

    return true;
  }
}
