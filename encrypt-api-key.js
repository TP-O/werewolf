const dotenv = require('dotenv');
const aes256 = require('aes256');
const { env } = require('process');

dotenv.config();

console.log(aes256.encrypt(env.APP_SECRET_KEY, env.APP_DECRYPTED_API_KEY));
