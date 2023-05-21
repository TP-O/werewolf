import type { PlayerId, PlayerStatus, RoomId, WaitingRoom } from '~/types'

export enum ListenEvent {
  Error = 'error',
  Success = 'success',
  FriendStatus = 'friend_status',
  PrivateMessage = 'private_message',
  RoomMessage = 'room_message',
  RoomChange = 'room_change',
}

export enum EmitEvent {
  RoomMessage = 'room_message',
}

export enum RoomChangeType {
  Join,
  Leave,
  Owner,
  Setting,
}

export interface SuccessResponse {
  message: string
}

export interface ErrorResponse {
  event: ListenEvent
  message: string | string[]
}

export interface FriendStatusData {
  id: PlayerId
  status: PlayerStatus
}

export interface PrivateMessageData {
  senderId: PlayerId
  content: string
}

export type RoomMessageData = PrivateMessageData & {
  roomId: string
}

export interface RoomData {
  changeType: RoomChangeType
  changerId?: PlayerId
  room: Pick<WaitingRoom, 'id'> & Partial<WaitingRoom>
}

export interface SendRoomMessage {
  roomId: RoomId
  content: string
}

export interface ListenEventMap {
  [ListenEvent.Error]: (response: ErrorResponse) => void
  [ListenEvent.Success]: (response: SuccessResponse) => void
  [ListenEvent.FriendStatus]: (data: FriendStatusData) => void
  [ListenEvent.PrivateMessage]: (data: PrivateMessageData) => void
  [ListenEvent.RoomMessage]: (data: RoomMessageData) => void
  [ListenEvent.RoomChange]: (data: RoomData) => void
}

export interface EmitEventMap {
  [EmitEvent.RoomMessage]: (data: SendRoomMessage) => void
}
