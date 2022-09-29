export type Room = {
  id: string;
  isPublic: boolean;
  isPersistent: boolean;
  ownerId: number;
  memberIds: number[];
  waitingIds: number[];
  refusedIds: number[];
};
