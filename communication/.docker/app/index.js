// eslint-disable-next-line @typescript-eslint/no-var-requires
const { existsSync } = require('fs');

(async () => {
  while (!existsSync(process.env.CONFIG_FILE ?? 'config.yaml')) {
    console.log("Config file doesn't exist!");
    await new Promise((resolve) => setTimeout(resolve, 5000));
  }

  require('../../dist/main');
})();
