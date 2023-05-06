import { PlayerId } from '../player/player.type';

export type RoomId = string;

export type Room = {
  id: RoomId;
  isMuted: boolean;
  password?: string;
  ownerId?: PlayerId;
  memberIds: PlayerId[];
};
