import { RedisOptions } from 'ioredis';
import { Time } from 'src/enum/time.enum';
import { env } from 'process';

export const RedisConfig: RedisOptions & { url: string } = Object.freeze({
  host: env.REDIS_HOST || 'redis',
  port: parseInt(env.REDIS_PORT, 10) || 6379,
  url: `redis://${env.REDIS_HOST || 'redis'}:${
    parseInt(env.REDIS_PORT, 10) || 6379
  }`,
  password: env.REDIS_PASSWORD || '',
  ttl: 1 * Time.Hour,
});
