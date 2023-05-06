import { Player } from '@prisma/client';

export type PlayerId = string;

export type SocketId = string;

export type OnlinePlayer = Player & {
  sid: SocketId;
};
