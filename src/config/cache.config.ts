import { RedisClientOptions } from 'redis';
import * as redisStore from 'cache-manager-redis-store';
import { CacheModuleOptions } from '@nestjs/common';
import { Time } from 'src/enum/time.enum';
import { env } from 'process';

export const CacheConfig: CacheModuleOptions & RedisClientOptions =
  Object.freeze({
    store: redisStore,
    url: `redis://${env.REDIS_HOST || 'redis'}:${
      parseInt(env.REDIS_PORT, 10) || 6379
    }`,
    password: env.REDIS_PASSWORD || '',
    ttl: 1 * Time.Hour,
  });
