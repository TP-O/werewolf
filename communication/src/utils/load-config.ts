import { fileLoader } from 'nest-typed-config';
import { RootConfig } from 'src/config';

let config: RootConfig;

export const loadConfig = () => {
  if (config) {
    return config;
  }

  config = fileLoader({
    absolutePath: process.env.CONFIG_FILE ?? 'config.yaml',
  })() as RootConfig;
  return config;
};
