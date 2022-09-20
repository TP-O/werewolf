import { ActiveStatus, EmitEvent, RoomEvent } from 'src/enum';
import { Room } from 'src/module/chat/type';
import { WsErrorResponse } from './error.type';

type SuccessResponse = {
  message: string;
};

type UpdateFriendStatusData = {
  id: number;
  status: ActiveStatus | null;
};

type ReceivePrivateMessageData = {
  senderId: number;
  content: string;
};

type ReceiveGroupMessageData = ReceivePrivateMessageData & {
  roomId: string;
};

type ReceiveRoomInvitationData = {
  inviterId: number;
  roomId: string;
};

type KickedOutOfRoomData = {
  roomId: string;
  kickerId: number;
};

type ReceiveRoomChangesData = {
  event: RoomEvent;
  actorId: number;
  room: Partial<Room> & Pick<Room, 'id'>;
};

export type EmitEvents = {
  [EmitEvent.Error]: (response: WsErrorResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.UpdateFriendStatus]: (data: UpdateFriendStatusData) => void;
  [EmitEvent.ReceivePrivateMessage]: (data: ReceivePrivateMessageData) => void;
  [EmitEvent.ReceiveGroupMessage]: (data: ReceiveGroupMessageData) => void;
  [EmitEvent.ReceiveRoomInvitation]: (data: ReceiveRoomInvitationData) => void;
  [EmitEvent.KickedOutOfRoom]: (data: KickedOutOfRoomData) => void;
  [EmitEvent.ReceiveRoomChanges]: (data: ReceiveRoomChangesData) => void;
};
