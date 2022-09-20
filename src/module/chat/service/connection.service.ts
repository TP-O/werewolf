import { Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import { User } from '@prisma/client';
import { Server, Socket } from 'socket.io';
import { AppConfig } from 'src/config';
import { UserId } from 'src/enum';
import { AuthService } from 'src/service/auth.service';
import { UserService } from 'src/service/user.service';

@Injectable()
export class ConnectionService {
  constructor(
    private authService: AuthService,
    private userService: UserService,
  ) {}

  /**
   * Verify token.
   *
   * @param headerAuthorization
   * @returns user record.
   */
  private async validateAuthorization(headerAuthorization: string) {
    const token = String(headerAuthorization).replace('Bearer ', '');
    const user = await this.authService.getUser(token);

    return user;
  }

  /**
   * Disconnect an user.
   *
   * @param server websocket server.
   * @param user
   */
  private async forceDisconnect(server: Server, user: User) {
    await this.userService.disconnect(user);
    user.sids.forEach((id) => server.sockets.sockets.get(id).disconnect());

    throw new WsException('Your account is in use. Please connect again!');
  }

  /**
   * Check if the connection satisfies some sepecific conditions
   * before allowing the connection.
   *
   * @param server websocket server.
   * @param client socket client.
   * @returns user record.
   */
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

    if (user.sids.length !== 0 && !AppConfig.allowDuplicateSignIn) {
      await this.forceDisconnect(server, user);
    }

    return user;
  }
}
