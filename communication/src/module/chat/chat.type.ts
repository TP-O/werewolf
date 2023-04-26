import { WsErrorResponse } from 'src/common/type';
import { PlayerId, PlayerStatus } from '../player';
import { EmitEvent, RoomEvent } from './chat.enum';
import { Room } from '../room';

type SuccessResponse = {
  message: string;
};

type UpdateFriendStatusData = {
  id: PlayerId;
  status: PlayerStatus;
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

export type EmitEventFunc = {
  [EmitEvent.Error]: (response: WsErrorResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.UpdateFriendStatus]: (data: UpdateFriendStatusData) => void;
  [EmitEvent.ReceivePrivateMessage]: (data: ReceivePrivateMessageData) => void;
  [EmitEvent.ReceiveRoomMessage]: (data: ReceiveRoomMessageData) => void;
  [EmitEvent.ReceiveRoomInvitation]: (data: ReceiveRoomInvitationData) => void;
  [EmitEvent.ReceiveRoomChanges]: (data: ReceiveRoomChangesData) => void;
};
