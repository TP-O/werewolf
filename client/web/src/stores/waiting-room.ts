import { info } from 'loglevel'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { useQuasar } from 'quasar'
import type {
  RoomData,
  RoomMessageData,
} from '~/composables/socketio/communication/types'
import {
  EmitEvent,
  RoomChangeType,
} from '~/composables/socketio/communication/types'
import type { PlayerId, RoomId } from '~/types'

export interface WaitingRoom {
  id: RoomId
  isMuted: boolean
  password?: string
  ownerId: PlayerId
  memberIds: PlayerId[]
}

interface WaitingRoomMessage {
  senderId?: PlayerId
  content: string
}

export const useWaitingRoomStore = defineStore('waiting_room', () => {
  const room = ref<WaitingRoom | null>(null)
  const messages = ref<WaitingRoomMessage[]>([])
  const player = usePlayerStore()
  const router = useRouter()
  const $q = useQuasar()

  async function book() {
    const { data: { data } } = await commApi.post('/rooms', {
      password: '12345',
    })
    room.value = data

    return data
  }

  async function join(id: string, password?: string) {
    if (room.value)
      return

    const res = await commApi.post('/rooms/join', {
      id,
      password,
    })

    if (res.statusCode === 404) {
      router.push('/')
      $q.notify({
        color: 'error',
        message: 'Room does not exist',
      })
      return null
    }
    else if (res.statusCode === 400) {
      router.push('/')
      $q.dialog({
        title: 'Enter room password',
        prompt: {
          model: '',
          type: 'text',
        },
        cancel: true,
        persistent: true,
      }).onOk((password) => {
        join(id, password).then(() =>
          router.push(`/room/${id}`))
      })
      return null
    }

    room.value = res.data.data
    messages.value = [] // For sure :D

    return null
  }

  async function leave(id: string) {
    if (!room.value)
      return

    await commApi.post('/rooms/leave', {
      id,
    })
    room.value = null
    messages.value = []
  }

  async function kick(memberId: PlayerId) {
    if (!room.value)
      return

    await commApi.post('/rooms/kick', {
      id: room.value.id,
      memberId,
    })
  }

  function sendMessage(content: string) {
    if (!room.value)
      return

    const data = {
      roomId: room.value.id,
      content,
    }
    info(`Send event ${EmitEvent.RoomMessage}:`, data)
    useCommSocket(socket => socket.emit(EmitEvent.RoomMessage, data))
  }

  function onRoomChange(event: RoomData) {
    if (room.value && event.changeType === RoomChangeType.Leave) {
      if (event.room.memberIds?.includes(player.player?.username || '')) {
        room.value = null
        messages.value = []
        router.push('/')
        $q.notify('You left the room')
      }
      else {
        removeMember(...(event.room.memberIds || []))
      }
    }
    else if (event.changeType === RoomChangeType.Join) {
      if (room.value) {
        addMember(...(event.room.memberIds || []))
      }
      else {
        room.value = event.room as WaitingRoom
        messages.value = []
        router.push(`/room/${room.value.id}`)
      }
    }
  }

  function removeMember(...memberId: PlayerId[]) {
    if (!room.value)
      return

    memberId.forEach((mId) => {
      const i = room.value!.memberIds.indexOf(mId)
      if (i === -1)
        return

      room.value!.memberIds.splice(i, 1)
      messages.value.push({
        content: `${mId} has leaved`,
      })
    })
  }

  function addMember(...memberId: PlayerId[]) {
    if (!room.value)
      return

    memberId.forEach((mId) => {
      if (!room.value!.memberIds.includes(mId)) {
        room.value!.memberIds.push(mId)
        messages.value.push({
          content: `${mId} has joined`,
        })
      }
    })
  }

  function onRoomMessage(data: RoomMessageData) {
    messages.value.push(data)
  }

  return {
    room,
    messages,
    book,
    join,
    leave,
    kick,
    sendMessage,
    onRoomChange,
    onRoomMessage,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useWaitingRoomStore as any, import.meta.hot))
}
