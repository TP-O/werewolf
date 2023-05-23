<script setup lang="ts">
import { info } from 'loglevel'
import { storeToRefs } from 'pinia'
import { EmitEvent } from '~/composables/socketio/communication/types'

const route = useRoute()
const router = useRouter()
const { player } = storeToRefs(usePlayerStore())
const { room, messages } = storeToRefs(useWaitingRoomStore())
const { join, leave, kick } = useWaitingRoomStore()
const roomId = String(route.params.id)
const messageInput = ref('')
const boxChat = ref<HTMLDivElement>()

function sendMessage() {
  const data = {
    roomId,
    content: messageInput.value,
  }

  messageInput.value = ''
  info(`Send event ${EmitEvent.RoomMessage}:`, data)
  useCommSocket(socket => socket.emit(EmitEvent.RoomMessage, data))
}

function leaveRoom() {
  leave(roomId)
  router.push('/')
}

watch(messages.value, () => {
  if (boxChat.value)
    boxChat.value.scrollTop = boxChat.value.scrollHeight
})

onBeforeMount(async () => {
  if (!room.value || room.value?.id !== roomId)
    await join(roomId)
})
</script>

<template>
  <div h-full>
    <div h="1/12" grid="~ cols-[0.5fr_11fr_0.5fr]" py-4>
      <q-btn color="negative" label="Exit" @click="leaveRoom" />

      <p text-xl>
        Room: <b>{{ roomId }}</b>
      </p>
    </div>

    <div h="11/12" flex gap-1>
      <div flex="~ col" gap-1 w="1/3">
        <div
          h="1/2" overflow-x-hidden overflow-y-scroll border rounded p-2
        >
          <div mb-2>
            Joined players
          </div>

          <div
            v-for="memberId, i of room?.memberIds" :key="i"
            flex="~  justify-between"
            my-2 gap-2 border p-2
          >
            <p
              flex="~ items-center"
              overflow-hidden text-ellipsis whitespace-nowrap
            >
              {{ memberId }}
            </p>
            <q-btn>
              <div i="carbon-overflow-menu-horizontal" />
              <q-menu min-w="100px" dense>
                <q-list>
                  <q-item
                    v-if="room?.ownerId === player?.username"
                    v-close-popup clickable @click="kick(memberId)"
                  >
                    <q-item-section>Kick</q-item-section>
                  </q-item>
                  <q-item
                    v-close-popup clickable
                  >
                    <q-item-section>Profile</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
        </div>

        <div h="1/2" relative border rounded p-2>
          <div
            ref="boxChat"
            h="10/12" overflow-x-hidden overflow-y-scroll px-1 text-left
          >
            <div v-for="message, i in messages" :key="i">
              <p mb-2 break-words>
                <b
                  v-if="message.senderId"
                  :class="message.senderId === player?.username
                    ? 'text-blue' : ''"
                >
                  {{ `${message.senderId}:` }}
                </b>
                {{ message.content }}
              </p>
            </div>
          </div>

          <div absolute bottom-0 left-0 w-full bg-white p-2>
            <q-input
              v-model="messageInput"
              dense autogrow outlined
              h="max-[200px]"
              @keydown.enter.prevent="sendMessage"
            />
          </div>
        </div>
      </div>

      <div border rounded p-1 w="1/3">
        Selected role
      </div>

      <div border rounded p-1 w="1/3">
        Settings
      </div>
    </div>
  </div>
</template>

<style>
textarea {
  @apply max-h-[100px];
}
</style>

<route lang="yaml">
meta:
  layout: client
</route>
