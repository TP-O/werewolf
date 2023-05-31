import type { PlayerId } from './player'

export type RoomId = string

export interface WaitingRoom {
  id: RoomId
  isMuted: boolean
  ownerId: PlayerId
  memberIds: PlayerId[]
}
