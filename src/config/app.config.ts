import { env } from 'process';
import { AppEnv } from 'src/enum';

export const AppConfig = Object.freeze(
  (() => {
    const parsedPort = parseInt(env.APP_PORT, 10);

    return {
      env: Object.values(AppEnv).includes(env.APP_ENV as AppEnv)
        ? env.APP_ENV
        : AppEnv.Development,
      debug: env.APP_DEBUG === 'true',
      port: parsedPort >= 0 && parsedPort < 65536 ? parsedPort : 3000,
      allowDuplicateSignIn: false,
      allowJoinMultipleRooms: true,
    };
  })(),
);
