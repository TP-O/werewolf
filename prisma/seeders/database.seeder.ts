import { PrismaClient } from '@prisma/client';
import { seedFriendRelationship } from './friend-relationship.seeder';
import { seedUser } from './user.seeder';
const client = new PrismaClient();

async function main() {
  await seedUser(client);
  await seedFriendRelationship(client);
}

main()
  .catch((e) => {
    console.error(e);
  })
  .finally(async () => {
    await client.$disconnect();
  });
