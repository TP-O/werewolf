import type { Socket } from 'socket.io-client'
import { io } from 'socket.io-client'
import merge from 'just-merge'
import { error, info } from 'loglevel'
import type {
  EmitEventMap,
  ErrorResponse,
  FriendStatusData,
  ListenEventMap,
  PrivateMessageData,
  RoomData,
  RoomMessageData,
  SuccessResponse,
} from './types'
import {
  ListenEvent,
  RoomChangeType,
} from './types'
import { useClientStore } from '~/stores/client'

const roomStore = useRoomStore()
const messageStore = useMessageStore()
const dialogStore = useDialogStore()
const clientStore = useClientStore()

let socket: Socket<ListenEventMap, EmitEventMap>

/**
 * Use this to make sure socket client is always ready.
 *
 * @param cb The callback function.
 */
export async function useCommSocket(
  cb?: (socket: Socket<ListenEventMap, EmitEventMap>) => void,
) {
  if (!socket || !socket.connected)
    await connect()

  if (cb)
    cb(socket)
}

/**
 * Connect to the server.
 *
 * @param reconnect Force reconnect to the server.
 */
async function connect(reconnect = false) {
  if (socket?.connected && !reconnect)
    return

  const token = await auth.getIdToken()
  socket = io(import.meta.env.VITE_COMMUNICATION_SERVER, {
    extraHeaders: {
      authorization: `Bearer ${token}`,
    },
  })

  socket.on('connect', onConnect)
  socket.on('connect_error', onConnectError)
  socket.on('disconnect', onDisconnect)
  socket.on(ListenEvent.Error, onError)
  socket.on(ListenEvent.Success, onSuccess)
  socket.on(ListenEvent.FriendStatus, onFriendStatus)
  socket.on(ListenEvent.RoomChange, onRoomChange)
  socket.on(ListenEvent.PrivateMessage, onPrivateMessage)
  socket.on(ListenEvent.RoomMessage, onRoomMessage)
}

function onConnect() {
  info('Connected to communication server')
  clientStore.setIsCommServerConnected(true)
}

function onDisconnect() {
  info('Disconnected from communication server')
  clientStore.setIsCommServerConnected(false)
}

function onConnectError() {
  dialogStore.openErrorDialog('Unable to connect to server')
}

function onError(res: ErrorResponse) {
  error(`Data from event ${ListenEvent.Error}:`, res)
  dialogStore.openErrorDialog(String(res.message))
}

function onSuccess(res: SuccessResponse) {
  info(`Data from event ${ListenEvent.Success}:`, res)
}

function onFriendStatus(data: FriendStatusData) {
  info(`Data from event ${ListenEvent.FriendStatus}:`, data)
}

function onRoomChange(data: RoomData) {
  info(`Data from event ${ListenEvent.RoomChange}:`, data)

  if (roomStore.waitingRoom?.id === data.room.id) {
    merge(roomStore.waitingRoom, data.room)

    if (data.changeType === RoomChangeType.Join) {
      data.room.memberIds?.forEach((mId) => {
        messageStore.addRoomMessage({
          roomId: data.room.id,
          content: `${mId} has joined`,
          senderId: '',
        })
      })
    }
    else if (data.changeType === RoomChangeType.Leave) {
      data.room.memberIds?.forEach((mId) => {
        if (mId === '') {
          // Kicked =))
        }

        messageStore.addRoomMessage({
          roomId: data.room.id,
          content: `${mId} has leaved`,
          senderId: '',
        })
      })
    }
  }
}

function onPrivateMessage(data: PrivateMessageData) {
  info(`Data from event ${ListenEvent.PrivateMessage}:`, data)
}

function onRoomMessage(data: RoomMessageData) {
  info(`Data from event ${ListenEvent.RoomMessage}:`, data)
  messageStore.addRoomMessage(data)
}
