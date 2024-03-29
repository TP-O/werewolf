<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import type { ResponseError } from '~/types'

defineOptions({
  name: 'WaitingRoomPage',
})

const route = useRoute()
const router = useRouter()
const $q = useQuasar()
const { room } = storeToRefs(useWaitingRoomStore())
const { join, leave } = useWaitingRoomStore()

const roomId = String(route.params.id)

async function leaveRoom() {
  $q.loading.show({
    message: 'Leaving...',
  })
  await leave(roomId)
}

function joinRoom(password?: string) {
  return join(roomId, password).catch(
    ({ data: { statusCode } }: ResponseError) => {
      if (statusCode === 400) {
        $q.dialog({
          title: 'Enter room password',
          prompt: {
            model: '',
            type: 'text',
          },
          cancel: true,
          persistent: true,
        }).onOk((password: string) => joinRoom(password))
      } else {
        router.push('/')

        if (statusCode === 404) {
          throw new Error('Room does not exist')
        } else {
          throw new Error('Unable to join this room')
        }
      }
    }
  )
}

onBeforeMount(async () => {
  $q.loading.show({
    message: 'Loading room...',
  })

  if (!room.value || room.value?.id !== roomId) {
    await joinRoom()
  }

  $q.loading.hide()
})

onUnmounted(() => {
  if ($q.loading.isActive) {
    $q.loading.hide()
  }
})
</script>

<template>
  <div h-full flex gap-1>
    <JoinedPlayerList w="1/3" />

    <PickedRoleList w="1/3" />

    <div flex="~ col" gap-1 w="1/3">
      <RoomSettings h="1/2" />
      <RoomChatBox :id="room?.id || ''" h="1/2" />

      <div p-2>
        <q-btn color="negative" label="Exit" w-full @click="leaveRoom" />
      </div>
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: client
</route>
