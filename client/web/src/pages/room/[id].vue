<script setup lang="ts">
import { info } from 'loglevel'
import { EmitEvent } from '~/composables/socketio/communication/types'

const route = useRoute()
const router = useRouter()
const roomStore = useRoomStore()
const room = computed(() => roomStore.waitingRoom)
const roomId = String(route.params.id)
const messageStore = useMessageStore()
const messages = computed(() => messageStore.roomMessages[roomId])
const { player } = usePlayerStore()
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
  roomStore.leaveRoom(roomId)
  router.push('/')
}

watch(messages.value, () => {
  if (boxChat.value)
    boxChat.value.scrollTop = boxChat.value.scrollHeight
})

onBeforeMount(async () => {
  if (!room.value || room.value.id !== roomId)
    await roomStore.joinRoom(roomId)
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
        <div h="1/2" overflow-x-hidden overflow-y-scroll border rounded p-2>
          <div>Joined players</div>

          <div v-for="memberId, i of room?.memberIds" :key="i" p-2>
            Player {{ memberId }}
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
              autogrow dense outlined
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
