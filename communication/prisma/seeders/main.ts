import { PrismaClient } from '@prisma/client';
import { seedFriendRelationship } from './friend-relationship';
import { seedPlayerStatus } from './player-status';
import { seedPlayer } from './player';
const client = new PrismaClient();

async function main() {
  await seedPlayerStatus(client);
  await seedPlayer(client);
  await seedFriendRelationship(client);
}

main()
  .then(async () => {
    await client.$disconnect();
  })
  .catch(async (e) => {
    await client.$disconnect();
    console.log(e);
    process.exit(1);
  });
