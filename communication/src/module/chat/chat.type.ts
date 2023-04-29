import { WsErrorResponse } from 'src/common/type';
import { PlayerId, PlayerStatus } from '../player';
import { EmitEvent, RoomChangeType } from './chat.enum';
import { Room } from '../room';

type SuccessResponse = {
  message: string;
};

type FriendStatusData = {
  id: PlayerId;
  status: PlayerStatus;
};

type PrivateMessageData = {
  senderId: PlayerId;
  content: string;
};

type RoomMessageData = PrivateMessageData & {
  roomId: string;
};

type RoomData = {
  changeType: RoomChangeType;
  changerId?: PlayerId;
  room: Pick<Room, 'id'> & Partial<Room>;
};

export type EmitEventMap = {
  [EmitEvent.Error]: (response: WsErrorResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.FriendStatus]: (data: FriendStatusData) => void;
  [EmitEvent.PrivateMessage]: (data: PrivateMessageData) => void;
  [EmitEvent.RoomMessage]: (data: RoomMessageData) => void;
  [EmitEvent.RoomChange]: (data: RoomData) => void;
};
