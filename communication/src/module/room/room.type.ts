import { PlayerId } from '../player';

export type Room = {
  id: string;
  isPublic: boolean;
  isPersistent: boolean;
  isMuted: boolean;
  gameId: number;
  ownerId: PlayerId;
  memberIds: PlayerId[];
  waitingIds: PlayerId[];
  refusedIds: PlayerId[];
};
