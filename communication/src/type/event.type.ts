import { ActiveStatus, EmitEvent, RoomEvent } from 'src/enum';
import { Room } from 'src/module/room/room.type';
import { WsErrorResponse } from './error.type';
import { PlayerId } from 'src/module/user/player.type';

type SuccessResponse = {
  message: string;
};

type UpdateFriendStatusData = {
  id: PlayerId;
  status: ActiveStatus | null;
};

type ReceivePrivateMessageData = {
  senderId: PlayerId;
  content: string;
};

type ReceiveRoomMessageData = ReceivePrivateMessageData & {
  roomId: string;
};

type ReceiveRoomInvitationData = {
  inviterId: PlayerId;
  roomId: string;
};

type ReceiveRoomChangesData = {
  event: RoomEvent;
  actorIds: PlayerId[];
  room: Partial<Room> & Pick<Room, 'id'>;
};

export type EmitEvents = {
  [EmitEvent.Error]: (response: WsErrorResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.UpdateFriendStatus]: (data: UpdateFriendStatusData) => void;
  [EmitEvent.ReceivePrivateMessage]: (data: ReceivePrivateMessageData) => void;
  [EmitEvent.ReceiveRoomMessage]: (data: ReceiveRoomMessageData) => void;
  [EmitEvent.ReceiveRoomInvitation]: (data: ReceiveRoomInvitationData) => void;
  [EmitEvent.ReceiveRoomChanges]: (data: ReceiveRoomChangesData) => void;
};
