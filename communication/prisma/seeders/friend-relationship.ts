import { PrismaClient } from '@prisma/client';

export async function seedFriendRelationship(client: PrismaClient) {
  return client.friendRelationship.createMany({
    data: [
      {
        senderId: 'RUckX55wLWWMtI7uuB0IxYTXipE2',
        acceptorId: 'tqR9BYe4RjQEuPH8QJVoBsknoDJ2',
      },
      {
        senderId: 'RUckX55wLWWMtI7uuB0IxYTXipE2',
        acceptorId: 'PmuDVqanntY08j0saW1qPKZo1Yl1',
      },
      {
        senderId: '1ZEa2Qma8FfFrxQHn3z6T0eFhPr1',
        acceptorId: 'RUckX55wLWWMtI7uuB0IxYTXipE2',
      },
      {
        senderId: '1ZEa2Qma8FfFrxQHn3z6T0eFhPr1',
        acceptorId: 'tqR9BYe4RjQEuPH8QJVoBsknoDJ2',
      },
      {
        senderId: '1ZEa2Qma8FfFrxQHn3z6T0eFhPr1',
        acceptorId: 'PmuDVqanntY08j0saW1qPKZo1Yl1',
      },
    ],
  });
}
