import { Injectable } from '@nestjs/common';
import { Redis } from 'ioredis';
import { RedisConfig } from 'src/config/redis';

@Injectable()
export class RedisService {
  private _client: Redis;

  constructor(config: RedisConfig) {
    this._client = new Redis({
      host: config.host,
      port: config.port,
      password: config.password,
    });
  }

  public get client() {
    return this._client;
  }
}
