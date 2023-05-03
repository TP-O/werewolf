import { Injectable, OnModuleDestroy } from '@nestjs/common';
import { Redis } from 'ioredis';
import { RedisConfig } from 'src/config';
import { LoggerService } from './logger.service';

@Injectable()
export class RedisService implements OnModuleDestroy {
  private readonly _client: Redis;

  constructor(config: RedisConfig, logger: LoggerService) {
    this._client = new Redis({
      host: config.host,
      port: config.port,
      password: config.password,
      enableAutoPipelining: true,
      maxRetriesPerRequest: 20,
    });
    this._client.on('error', (err) => {
      logger.error(err.message, err.stack, RedisService.name);
    });
  }

  async onModuleDestroy() {
    await this._client.quit();
  }

  public get client() {
    return this._client;
  }
}
