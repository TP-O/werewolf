import type { Socket } from 'socket.io-client'
import { io } from 'socket.io-client'
import log from 'loglevel'
import type { CommEmitEventMap, CommListenEventMap } from '~/types'
import { CommListenEvent } from '~/enums'

let socket: Socket<CommListenEventMap, CommEmitEventMap>

/**
 * Use this to make sure socket client is always ready.
 *
 * @param cb The callback function.
 */
export async function useCommSocket(
  cb?: (socket: Socket<CommListenEventMap, CommEmitEventMap>) => void
) {
  if (!socket || !socket.connected) {
    await connect()
  }

  if (cb) {
    cb(socket)
  }
}

/**
 * Connect to the server.
 *
 * @param reconnect Force reconnect to the server.
 */
async function connect(reconnect = false) {
  if (!reconnect && socket?.connected) {
    return
  }

  // Prevent connection conflict
  if (socket) {
    socket.close()
  }

  const token = await auth.getIdToken()
  socket = io(import.meta.env.VITE_COMMUNICATION_SERVER, {
    extraHeaders: {
      authorization: `Bearer ${token}`,
    },
  })

  const roomStore = useWaitingRoomStore()
  const clientStore = useClientStore()

  socket.on('connect', () => {
    log.info('Connected to communication server')
    clientStore.onConnectComm()
  })

  socket.on('connect_error', () => {
    log.error('Unable to connect to communication server')
    clientStore.onErroConnectComm()
  })

  socket.on('disconnect', () => {
    log.error('Disconnected from communication server')
    clientStore.onDisconnectComm()
  })

  socket.on(CommListenEvent.Error, (res) => {
    log.error(`Data from event ${CommListenEvent.Error}:`, res)
    clientStore.onCommError(res)
  })

  socket.on(CommListenEvent.Success, (res) => {
    log.info(`Data from event ${CommListenEvent.Success}:`, res)
    clientStore.onCommSuccess(res)
  })

  socket.on(CommListenEvent.FriendStatus, (event) => {
    log.info(`Data from event ${CommListenEvent.FriendStatus}:`, event)
  })

  socket.on(CommListenEvent.RoomChange, (event) => {
    log.info(`Data from event ${CommListenEvent.RoomChange}:`, event)
    roomStore.onRoomChange(event)
  })

  socket.on(CommListenEvent.PrivateMessage, (event) => {
    log.info(`Data from event ${CommListenEvent.PrivateMessage}:`, event)
  })

  socket.on(CommListenEvent.RoomMessage, (event) => {
    log.info(`Data from event ${CommListenEvent.RoomMessage}:`, event)
    roomStore.onRoomMessage(event)
  })
}
