import { EmitEvent } from 'src/enum/event.enum';
import { RoomChange } from 'src/enum/room.enum';
import { ActiveStatus } from 'src/enum/user.enum';
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

type ReceiveGroupMessageData = {
  roomId: string;
  senderId: number;
  content: string;
};

type ReceiveRoomChangesData = {
  roomId: string;
  memberIds: number[];
  change: {
    type: RoomChange;
    memeberId: number;
  };
};

type ReceiveRoomInvitationData = {
  inviterId: number;
  roomId: string;
};

type KickedOutOfRoomData = {
  roomId: string;
  kickerId: number;
};

export type EmitEvents = {
  [EmitEvent.Error]: (response: WsErrorResponse) => void;
  [EmitEvent.Success]: (response: SuccessResponse) => void;
  [EmitEvent.UpdateFriendStatus]: (data: UpdateFriendStatusData) => void;
  [EmitEvent.ReceivePrivateMessage]: (data: ReceivePrivateMessageData) => void;
  [EmitEvent.ReceiveGroupMessage]: (data: ReceiveGroupMessageData) => void;
  [EmitEvent.ReceiveRoomChanges]: (data: ReceiveRoomChangesData) => void;
  [EmitEvent.ReceiveRoomInvitation]: (data: ReceiveRoomInvitationData) => void;
  [EmitEvent.KickedOutOfRoom]: (data: KickedOutOfRoomData) => void;
};
