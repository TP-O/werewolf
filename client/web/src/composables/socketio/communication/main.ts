import type { Socket } from 'socket.io-client'
import { io } from 'socket.io-client'
import merge from 'just-merge'
import log from 'loglevel'
import type {
  EmitEventMap,
  ErrorResponse,
  FriendStatusData,
  ListenEventMap,
  PrivateMessageData,
  RoomData,
  RoomMessageData,
  SendRoomMessage,
  SuccessResponse,
} from './types'
import { EmitEvent, ListenEvent } from './types'
import { useClientStore } from '~/stores/client'

const roomStore = useRoomStore()
const messageStore = useMessageStore()
const dialogStore = useDialogStore()
const clientStore = useClientStore()

let commSocket: Socket<ListenEventMap, EmitEventMap>

/**
 * Connect to the server.
 *
 * @param reconnect Force reconnect to the server.
 */
export async function connect(reconnect = false) {
  if (commSocket?.connected && !reconnect)
    return

  const token = await auth.getIdToken()
  commSocket = io(import.meta.env.VITE_COMMUNICATION_SERVER, {
    extraHeaders: {
      authorization: `Bearer ${token}`,
    },
  })

  commSocket.on('connect', onConnect)
  commSocket.on('connect_error', onConnectError)
  commSocket.on('disconnect', onDisconnect)
  commSocket.on(ListenEvent.Error, onError)
  commSocket.on(ListenEvent.Success, onSuccess)
  commSocket.on(ListenEvent.FriendStatus, onFriendStatus)
  commSocket.on(ListenEvent.RoomChange, onRoomChange)
  commSocket.on(ListenEvent.PrivateMessage, onPrivateMessage)
  commSocket.on(ListenEvent.RoomMessage, onRoomMessage)
}

function onConnect() {
  log.debug('Connected to communication server')
  clientStore.setIsCommServerConnected(true)
}

function onDisconnect() {
  log.info('Disconnected from communication server')
  clientStore.setIsCommServerConnected(false)
}

function onConnectError() {
  dialogStore.openErrorDialog('Unable to connect to server')
}

function onError(res: ErrorResponse) {
  log.info(`Data from event ${ListenEvent.Error}:`, res)
  dialogStore.openErrorDialog(String(res.message))
}

function onSuccess(res: SuccessResponse) {
  log.info(`Data from event ${ListenEvent.Success}:`, res)
}

function onFriendStatus(data: FriendStatusData) {
  log.info(`Data from event ${ListenEvent.FriendStatus}:`, data)
}

function onRoomChange(data: RoomData) {
  log.info(`Data from event ${ListenEvent.RoomChange}:`, data)

  if (roomStore.waitingRoom)
    merge(roomStore.waitingRoom, data.room)
}

function onPrivateMessage(data: PrivateMessageData) {
  log.info(`Data from event ${ListenEvent.PrivateMessage}:`, data)
}

function onRoomMessage(data: RoomMessageData) {
  log.info(`Data from event ${ListenEvent.RoomMessage}:`, data)
  messageStore.addRoomMessage(data)
}

export function sendRoomMessage(data: SendRoomMessage) {
  log.info(`Send event ${EmitEvent.RoomMessage}:`, data)
  commSocket.emit(EmitEvent.RoomMessage, data)
}
