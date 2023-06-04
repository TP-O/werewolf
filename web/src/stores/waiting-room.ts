import log from 'loglevel'
import { acceptHMRUpdate, defineStore } from 'pinia'
import { CommEmitEvent, RoleId, RoomChangeType } from '~/enums'
import type {
  PlayerId,
  ResponseData,
  RoomEvent,
  RoomMessageEvent,
  WaitingRoom,
} from '~/types'

interface WaitingRoomMessage {
  senderId?: PlayerId
  content: string
}

interface WaitingRoomSettings {
  capacity: number
  password?: string
  gameSettings: {
    roleIds: RoleId[]
    requiredRoleIds: RoleId[]
    turnDuration: number
    discussionDuration: number
  }
}

export const useWaitingRoomStore = defineStore('waiting_room', () => {
  const room = ref<WaitingRoom | null>(null)
  const settings = ref<WaitingRoomSettings>({
    capacity: 5,
    gameSettings: {
      roleIds: [RoleId.Villager, RoleId.Werewolf],
      requiredRoleIds: [RoleId.Villager, RoleId.Werewolf],
      turnDuration: 30,
      discussionDuration: 60,
    },
  })
  const messages = ref<WaitingRoomMessage[]>([])
  const player = usePlayerStore()

  async function book() {
    const {
      data: { data },
    } = await commApi.post<ResponseData<WaitingRoom>>('/rooms', {
      // password: '12345',
    })
    room.value = data
    return data
  }

  async function join(id: string, password?: string) {
    if (room.value) {
      return
    }

    const {
      data: { data },
    } = await commApi.post<ResponseData<WaitingRoom>>('/rooms/join', {
      id,
      password,
    })
    room.value = data
    messages.value = []
    return data
  }

  async function leave(id: string) {
    if (!room.value) {
      return
    }

    await commApi.post('/rooms/leave', {
      id,
    })
    room.value = null
    messages.value = []
  }

  async function kick(memberId: PlayerId) {
    if (!room.value) {
      return
    }

    await commApi.post('/rooms/kick', {
      id: room.value.id,
      memberId,
    })
  }

  function sendMessage(content: string) {
    if (!room.value) {
      return
    }

    const data = {
      roomId: room.value.id,
      content,
    }
    log.info(`Send event ${CommEmitEvent.RoomMessage}:`, data)
    useCommSocket((socket) => socket.emit(CommEmitEvent.RoomMessage, data))
  }

  function onRoomChange(event: RoomEvent) {
    if (room.value && event.changeType === RoomChangeType.Leave) {
      if (event.room.memberIds?.includes(player.player?.username || '')) {
        room.value = null
        messages.value = []
      } else {
        removeMember(...(event.room.memberIds || []))
      }
    } else if (event.changeType === RoomChangeType.Join) {
      if (room.value) {
        addMember(...(event.room.memberIds || []))
      } else {
        room.value = event.room as WaitingRoom
        messages.value = []
      }
    }
  }

  function removeMember(...memberId: PlayerId[]) {
    if (!room.value) {
      return
    }

    memberId.forEach((mId) => {
      const i = room.value!.memberIds.indexOf(mId)
      if (i === -1) {
        return
      }

      room.value!.memberIds.splice(i, 1)
      messages.value.push({
        content: `${mId} has leaved`,
      })
    })
  }

  function addMember(...memberId: PlayerId[]) {
    if (!room.value) {
      return
    }

    memberId.forEach((mId) => {
      if (!room.value!.memberIds.includes(mId)) {
        room.value!.memberIds.push(mId)
        messages.value.push({
          content: `${mId} has joined`,
        })
      }
    })
  }

  function onRoomMessage(event: RoomMessageEvent) {
    messages.value.push(event)
  }

  return {
    room,
    settings,
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
    acceptHMRUpdate(useWaitingRoomStore as any, import.meta.hot)
  )
}
