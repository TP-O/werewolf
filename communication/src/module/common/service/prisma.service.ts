import { INestApplication, Injectable, OnModuleInit } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { AppConfig, DbConfig } from 'src/config';
import { LoggerService } from './logger.service';

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit {
  private readonly _debug: boolean;

  private readonly logger: LoggerService;

  constructor(appConfig: AppConfig, dbConfig: DbConfig, logger: LoggerService) {
    const datasources = {
      db: {
        // eslint-disable-next-line prettier/prettier
        url: `postgresql://${dbConfig.username}:${ dbConfig.password}@${
          dbConfig.host}:${dbConfig.port}/${
          dbConfig.database}?schema=public&connection_limit=${dbConfig.pollSize}&pool_timeout=0`,
      },
    };
    if (appConfig.debug) {
      super({
        log: [
          {
            emit: 'event',
            level: 'query',
          },
        ],
        datasources,
      });
    } else {
      super({
        datasources,
      });
    }

    this._debug = appConfig.debug;
    this.logger = logger;
    this.logger.setContext(PrismaService.name);
  }

  async onModuleInit() {
    await this.$connect();

    if (this._debug) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      this.$on('query', async (e) => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        this.logger.log(`${e.query} ${e.params}`);
      });
    }
  }

  enableShutdownHooks(app: INestApplication) {
    this.$on('beforeExit', async () => {
      await app.close();
    });
  }
}
