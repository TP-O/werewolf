import type { Socket } from 'socket.io-client'
import { io } from 'socket.io-client'
import merge from 'just-merge'
import type { PlayerId, PlayerStatus, WaitingRoom } from '~/types'

enum ListenEvent {
  Error = 'error',
  Success = 'success',
  FriendStatus = 'friend_status',
  PrivateMessage = 'private_message',
  RoomMessage = 'room_message',
  RoomChange = 'room_change',
}

export enum RoomChangeType {
  Join,
  Leave,
  Owner,
  Setting,
}

interface SuccessResponse {
  message: string
}

interface ErrorResponse {
  event: ListenEvent
  message: string | string[]
}

interface FriendStatusData {
  id: PlayerId
  status: PlayerStatus
}

interface PrivateMessageData {
  senderId: PlayerId
  content: string
}

type RoomMessageData = PrivateMessageData & {
  roomId: string
}

interface RoomData {
  changeType: RoomChangeType
  changerId?: PlayerId
  room: Pick<WaitingRoom, 'id'> & Partial<WaitingRoom>
}

const roomStore = useRoomStore()
const messageStore = useMessageStore()
// eslint-disable-next-line import/no-mutable-exports
let communicationSocket: Socket

export async function connectCommunicationServer() {
  const token = await firebaseAuth.currentUser?.getIdToken()
  communicationSocket = io('http://127.0.0.1:8079/', {
    extraHeaders: {
      authorization: `Bearer ${token}`,
    },
  })

  communicationSocket.on(ListenEvent.Error, onError)
  communicationSocket.on(ListenEvent.Success, onSuccess)
  communicationSocket.on(ListenEvent.FriendStatus, onFriendStatus)
  communicationSocket.on(ListenEvent.RoomChange, onRoomChange)
  communicationSocket.on(ListenEvent.PrivateMessage, onPrivateMessage)
  communicationSocket.on(ListenEvent.RoomMessage, onRoomMessage)
}

function onError(res: ErrorResponse) {
  throw new Error(String(res.message))
}

function onSuccess(res: SuccessResponse) {
  console.log(res)
}

function onFriendStatus(data: FriendStatusData) {
  console.log(data)
}

function onRoomChange(data: RoomData) {
  console.log(data)

  if (roomStore.waitingRoom)
    merge(roomStore.waitingRoom, data.room)
}

function onPrivateMessage(data: PrivateMessageData) {
  console.log(data)
}

function onRoomMessage(data: RoomMessageData) {
  console.debug(data)

  messageStore.receiveRoomMessage(data)
}

export { communicationSocket }
