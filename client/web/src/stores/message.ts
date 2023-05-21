import { acceptHMRUpdate, defineStore } from 'pinia'
import type { PlayerId, RoomId } from '~/types'

interface RoomMessage {
  senderId: PlayerId
  content: string
}

type ReceiveRoomMessage = RoomMessage & {
  roomId: RoomId
}

interface SendRoomMessage {
  roomId: RoomId
  content: string
}

type RoomMessages = Record<RoomId, RoomMessage[]>

export const useMessageStore = defineStore('message', () => {
  const privateMessages = reactive({})
  const roomMessages = reactive<RoomMessages>({})

  const receiveRoomMessage = (receive: ReceiveRoomMessage) => {
    const message = {
      senderId: receive.senderId,
      content: receive.content,
    }
    if (!roomMessages[receive.roomId])
      roomMessages[receive.roomId] = [message]
    else
      roomMessages[receive.roomId].push(message)
  }

  const sendRoomMessage = (send: SendRoomMessage) => {
    communicationSocket.emit(
      'room_message',
      send,
    )
  }

  return {
    privateMessages,
    roomMessages,
    receiveRoomMessage,
    sendRoomMessage,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useMessageStore as any, import.meta.hot))
