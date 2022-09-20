export type Room = {
  id: string;
  ownerId: number;
  memberIds: number[];
  waitingIds: number[];
  refusedIds: number[];
};
