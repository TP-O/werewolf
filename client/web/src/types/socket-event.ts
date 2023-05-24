import type { PlayerId } from './player'
import type { RoomId, WaitingRoom } from './room'
import type {
  CommEmitEvent,
  CommListenEvent,
  PlayerStatus,
  RoomChangeType,
} from '~/enums'

export interface SuccessResponse {
  message: string
}

export interface ErrorResponse {
  event: CommListenEvent
  message: string | string[]
}

export interface FriendStatusEvent {
  id: PlayerId
  status: PlayerStatus
}

export interface PrivateMessageEvent {
  senderId: PlayerId
  content: string
}

export type RoomMessageEvent = PrivateMessageEvent & {
  roomId: string
}

export interface RoomEvent {
  changeType: RoomChangeType
  room: Pick<WaitingRoom, 'id'> & Partial<WaitingRoom>
}

export interface SendRoomMessageEvent {
  roomId: RoomId
  content: string
}

export interface CommListenEventMap {
  [CommListenEvent.Error]: (res: ErrorResponse) => void
  [CommListenEvent.Success]: (res: SuccessResponse) => void
  [CommListenEvent.FriendStatus]: (event: FriendStatusEvent) => void
  [CommListenEvent.PrivateMessage]: (event: PrivateMessageEvent) => void
  [CommListenEvent.RoomMessage]: (event: RoomMessageEvent) => void
  [CommListenEvent.RoomChange]: (event: RoomEvent) => void
}

export interface CommEmitEventMap {
  [CommEmitEvent.RoomMessage]: (event: SendRoomMessageEvent) => void
}
