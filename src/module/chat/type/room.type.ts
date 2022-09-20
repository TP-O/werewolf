export type Room = {
  id: string;
  isPublic: boolean;
  ownerId: number;
  memberIds: number[];
  waitingIds: number[];
  refusedIds: number[];
};
