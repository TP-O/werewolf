import { Type } from 'class-transformer';
import { ValidateNested } from 'class-validator';
import { AppConfig } from './app.config';
import { DbConfig } from './db.config';
import { RedisConfig } from './redis.config';
import { FirebaseConfig } from './firebase.config';

export class RootConfig {
  @Type(() => AppConfig)
  @ValidateNested()
  public readonly app!: AppConfig;

  @Type(() => DbConfig)
  @ValidateNested()
  public readonly db!: DbConfig;

  @Type(() => RedisConfig)
  @ValidateNested()
  public readonly redis!: RedisConfig;

  @Type(() => FirebaseConfig)
  @ValidateNested()
  public readonly firebase!: FirebaseConfig;
}
