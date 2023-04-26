import { PrismaClient } from '@prisma/client';
import { PlayerStatus } from 'src/enum';

export async function seedPlayerStatus(client: PrismaClient) {
  return client.playerStatus.createMany({
    data: [
      { id: PlayerStatus.Offline, name: 'offline' },
      { id: PlayerStatus.Online, name: 'online' },
    ],
  });
}
