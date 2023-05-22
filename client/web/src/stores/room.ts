import { acceptHMRUpdate, defineStore } from 'pinia'
import type { WaitingRoom } from '~/types'

export const useRoomStore = defineStore('room', () => {
  const waitingRoom = ref<WaitingRoom | null>(null)

  async function bookRoom() {
    const res = await commApi.post('/rooms', {
      password: 'default_password',
    })
    waitingRoom.value = res.data.data
  }

  async function joinRoom(id: string) {
    const res = await commApi.post('/rooms/join', {
      id,
      password: 'default_password',
    })
    waitingRoom.value = res.data.data
  }

  async function leaveRoom(id: string) {
    await commApi.post('/rooms/leave', {
      id,
    })
    waitingRoom.value = null
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
    leaveRoom,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useRoomStore as any, import.meta.hot))
