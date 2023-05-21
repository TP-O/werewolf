import { type ViteSSGContext } from 'vite-ssg'

export type UserModule = (ctx: ViteSSGContext) => void

export type PlayerId = string

export type RoomId = string

export type PlayerStatus = number

export interface WaitingRoom {
  id: RoomId
  isMuted: boolean
  password?: string
  ownerId: PlayerId
  memberIds: PlayerId[]
}
