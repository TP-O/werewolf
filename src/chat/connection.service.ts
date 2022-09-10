import { CACHE_MANAGER, Inject, Injectable } from '@nestjs/common';
import { WsException } from '@nestjs/websockets';
import { Cache } from 'cache-manager';
import { Socket } from 'socket.io';
import { AppConfig } from 'src/config/app.config';
import { CacheNamespace } from 'src/enum/cache.enum';
import { AuthService } from './auth.service';

@Injectable()
export class ConnectionService {
  constructor(
    private authService: AuthService,
    @Inject(CACHE_MANAGER) private cacheManager: Cache,
  ) {}

  private async parseUserId(headerAuthorization: string) {
    const token = String(headerAuthorization).replace('Bearer ', '');
    const userId = await this.authService.getUserId(token);

    return userId;
  }

  private async isDuplicateConnection(userId: string) {
    const socketId = await this.cacheManager.get<string>(
      `${CacheNamespace.UId2SId}${userId}`,
    );

    return socketId != null;
  }

  async validateConnection(client: Socket) {
    const userId = await this.parseUserId(
      client.handshake.headers.authorization,
    );

    if (userId === '') {
      throw new WsException('Invalid access token!');
    }

    if (
      !AppConfig.allowDuplicateSignIn &&
      (await this.isDuplicateConnection(userId))
    ) {
      throw new WsException('Your account is being used elsewhere!');
    }

    return userId;
  }

  async connect(client: Socket, userId: string) {
    await this.cacheManager.set(
      `${CacheNamespace.UId2SId}${userId}`,
      client.id,
    );
    await this.cacheManager.set(
      `${CacheNamespace.SId2UId}${client.id}`,
      userId,
    );

    console.log('connected!');
  }

  async disconnect(client: Socket) {
    const userId = this.cacheManager.get<string>(
      `${CacheNamespace.SId2UId}${client.id}`,
    );

    if (userId !== null) {
      await this.cacheManager.del(`${CacheNamespace.SId2UId}${client.id}`);
      await this.cacheManager.del(`${CacheNamespace.UId2SId}${userId}`);
    }

    console.log('disconnected');
  }
}
