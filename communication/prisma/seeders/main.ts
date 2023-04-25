import { PrismaClient } from '@prisma/client';
import { seedFriendRelationship } from './friend-relationship.seeder';
import { seedGame } from './game.seeder';
import { seedStatus } from './status.seeder';
import { seedUser } from './user.seeder';
const client = new PrismaClient();

async function main() {
  await seedStatus(client);
  await seedUser(client);
  await seedFriendRelationship(client);
  await seedGame(client);
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
