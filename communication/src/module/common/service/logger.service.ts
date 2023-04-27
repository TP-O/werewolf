import { ConsoleLogger, Injectable, Scope } from '@nestjs/common';
import { AppEnv } from 'src/common/enum';
import { AppConfig } from 'src/config';
import { createLogger, format, transports } from 'winston';

const writer = createLogger({
  format: format.combine(
    format.errors({ stack: true }),
    format.timestamp(),
    format.prettyPrint(),
  ),
  transports: [
    new transports.File({
      filename: 'log/error.log',
      level: 'error',
      maxsize: 10_000_000,
    }),
  ],
});

@Injectable({
  scope: Scope.TRANSIENT,
})
export class LoggerService extends ConsoleLogger {
  private readonly _logError!: boolean;

  private readonly _logWriter = writer;

  constructor(config: AppConfig) {
    super();
    this._logError = config.env === AppEnv.Production;
  }

  error(message: any, stack?: string, context?: string) {
    if (!this.isLevelEnabled('error')) {
      return;
    }

    if (context) {
      super.error(message, stack, context);
    } else {
      super.error(message, stack, this.context);
    }

    if (this._logError) {
      this._logWriter.error({
        context: this.context,
        message: message,
        stack,
      });
    }
  }
}
