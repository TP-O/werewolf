import { Injectable, OnModuleDestroy } from '@nestjs/common';
import { Redis } from 'ioredis';
import { RedisConfig } from 'src/config';

@Injectable()
export class RedisService implements OnModuleDestroy {
  private _client: Redis;

  constructor(config: RedisConfig) {
    this._client = new Redis({
      host: config.host,
      port: config.port,
      password: config.password,
      enableAutoPipelining: true,
    });
  }

  async onModuleDestroy() {
    await this._client.quit();
  }

  public get client() {
    return this._client;
  }
}
