import { Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import { User } from '@prisma/client';
import { Server, Socket } from 'socket.io';
import { AppConfig } from 'src/config';
import { EmitEvent, RoomEvent, UserId } from 'src/enum';
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
   * Solve conflict if multiple people connect to the
   * same account.
   *
   * @param server websocket server.
   * @param user
   */
  private async handleConflict(server: Server, user: User) {
    if (!AppConfig.disconnectIfConflict) {
      throw new WsException('Your account is in use!');
    }

    const { disconnectedId, leftRooms } = await this.userService.disconnect(
      user,
    );

    server.sockets.sockets.get(disconnectedId).disconnect();

    leftRooms.forEach((room) => {
      if (room.memberIds.length > 0) {
        server.to(room.id).emit(EmitEvent.ReceiveRoomChanges, {
          event: RoomEvent.Leave,
          actorId: user.id,
          room,
        });
      }
    });
  }

  /**
   * Check if the connection satisfies some sepecific conditions
   * before allowing the connection.
   *
   * @param server websocket server.
   * @param client socket client.
   * @returns updated user.
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

    if (user.statusId != null) {
      await this.handleConflict(server, user);
    }

    return user;
  }
}
