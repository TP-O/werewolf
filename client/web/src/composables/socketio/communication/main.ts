import type { Socket } from 'socket.io-client'
import { io } from 'socket.io-client'
import { error, info } from 'loglevel'
import type {
  EmitEventMap,
  ErrorResponse,
  FriendStatusData,
  ListenEventMap,
  PrivateMessageData,
  SuccessResponse,
} from './types'
import {
  ListenEvent,
} from './types'

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

  const roomStore = useWaitingRoomStore()
  const clientStore = useClientStore()

  socket.on('connect', () => {
    info('Connected to communication server')
    clientStore.onConnectComm()
  })

  socket.on('connect_error', () => {
    error('Unable to connect to communication server')
    clientStore.onErroConnectComm()
  })

  socket.on('disconnect', () => {
    error('Disconnected from communication server')
    clientStore.onDisconnectComm()
  })

  socket.on(ListenEvent.Error, (event: ErrorResponse) => {
    error(`Data from event ${ListenEvent.Error}:`, event)
    clientStore.onCommError(event)
  })

  socket.on(ListenEvent.Success, (event: SuccessResponse) => {
    info(`Data from event ${ListenEvent.Success}:`, event)
    clientStore.onCommSuccess(event)
  })

  socket.on(ListenEvent.FriendStatus, (event: FriendStatusData) => {
    info(`Data from event ${ListenEvent.FriendStatus}:`, event)
  })

  socket.on(ListenEvent.RoomChange, (event) => {
    info(`Data from event ${ListenEvent.RoomChange}:`, event)
    roomStore.onRoomChange(event)
  })

  socket.on(ListenEvent.PrivateMessage, (event: PrivateMessageData) => {
    info(`Data from event ${ListenEvent.PrivateMessage}:`, event)
  })

  socket.on(ListenEvent.RoomMessage, (event) => {
    info(`Data from event ${ListenEvent.RoomMessage}:`, event)
    roomStore.onRoomMessage(event)
  })
}
