import { acceptHMRUpdate, defineStore } from 'pinia'
import type { WaitingRoom } from '~/types'

export const useRoomStore = defineStore('room', () => {
  const waitingRoom = ref<WaitingRoom | null>(null)

  async function bookRoom() {
    const res = await communicationClient.post('/rooms/book', {
      password: '12345',
    })
    waitingRoom.value = res.data.data
  }

  async function joinRoom(id: string) {
    const res = await communicationClient.post('/rooms/join', {
      id,
      password: '12345',
    })
    waitingRoom.value = res.data.data
  }

  async function addMember(memberId: string) {
    if (waitingRoom.value)
      waitingRoom.value.memberIds.push(memberId)
  }

  return {
    waitingRoom,
    bookRoom,
    addMember,
    joinRoom,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useRoomStore as any, import.meta.hot))
