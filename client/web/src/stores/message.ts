import { acceptHMRUpdate, defineStore } from 'pinia'
import type { PlayerId, RoomId } from '~/types'

interface RoomMessage {
  senderId: PlayerId
  content: string
}

type ReceiveRoomMessage = RoomMessage & {
  roomId: RoomId
}

type RoomMessages = Record<RoomId, RoomMessage[]>

export const useMessageStore = defineStore('message', () => {
  const privateMessages = reactive({})
  const roomMessages = reactive<RoomMessages>({})

  const addRoomMessage = (receive: ReceiveRoomMessage) => {
    const message = {
      senderId: receive.senderId,
      content: receive.content,
    }
    if (!roomMessages[receive.roomId])
      roomMessages[receive.roomId] = [message]
    else
      roomMessages[receive.roomId].push(message)
  }

  return {
    privateMessages,
    roomMessages,
    addRoomMessage,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useMessageStore as any, import.meta.hot))
}
