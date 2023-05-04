import { PrismaClient } from '@prisma/client';
import { PlayerStatus } from 'src/module/player';

export async function seedPlayer(client: PrismaClient) {
  return client.player.createMany({
    data: [
      {
        id: 'RUckX55wLWWMtI7uuB0IxYTXipE2',
        statusId: PlayerStatus.Offline,
        username: 'player 01',
      },
      {
        id: 'tqR9BYe4RjQEuPH8QJVoBsknoDJ2',
        statusId: PlayerStatus.Offline,
        username: 'player 02',
      },
      {
        id: 'PmuDVqanntY08j0saW1qPKZo1Yl1',
        statusId: PlayerStatus.Offline,
        username: 'player 03',
      },
      {
        id: '1ZEa2Qma8FfFrxQHn3z6T0eFhPr1',
        statusId: PlayerStatus.Offline,
        username: 'player 04',
      },
      {
        id: 'Aw8sHDLaG4bS0EG9zrmgYzcueox1',
        statusId: PlayerStatus.Offline,
        username: 'player 05',
      },
    ],
  });
}
