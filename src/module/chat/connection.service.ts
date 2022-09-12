import { Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import { User } from '@prisma/client';
import { Server, Socket } from 'socket.io';
import { AppConfig } from 'src/config/app.config';
import { UserId } from 'src/enum/user.enum';
import { AuthService } from '../common/auth.service';
import { UserService } from '../common/user.service';

@Injectable()
export class ConnectionService {
  constructor(
    private authService: AuthService,
    private userService: UserService,
  ) {}

  private async validateAuthorization(headerAuthorization: string) {
    const token = String(headerAuthorization).replace('Bearer ', '');
    const user = await this.authService.getUser(token);

    return user;
  }

  private async forceDisconnect(server: Server, user: User) {
    (user.sids as string[]).forEach((sid) =>
      server.sockets.sockets.get(sid).disconnect(),
    );

    await this.userService.disconnect(user);

    throw new WsException('Your account is in use. Please connect again!');
  }

  async validateConnection(server: Server, client: Socket) {
    const user = await this.validateAuthorization(
      client.handshake.headers.authorization,
    );

    if (user.id === UserId.NonExist) {
      throw new WsException('Invalid access token!');
    }

    if (user.id === UserId.Asynchronous) {
      throw new WsException('Please connect again after a while!');
    }

    if (
      (user.sids as string[]).length !== 0 &&
      !AppConfig.allowDuplicateSignIn
    ) {
      await this.forceDisconnect(server, user);
    }

    return user;

    // return {
    //   id: 1,
    //   fid: '',
    //   sids: [],
    // };
  }
}
