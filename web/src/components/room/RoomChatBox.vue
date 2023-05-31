<script setup lang="ts">
import log from 'loglevel'
import { storeToRefs } from 'pinia'
import { CommEmitEvent } from '~/enums'
import type { RoomId } from '~/types'

const props = defineProps<{ id: RoomId }>()

const { player } = storeToRefs(usePlayerStore())
const { messages } = storeToRefs(useWaitingRoomStore())

const messageInput = ref('')
function sendMessage() {
  const data = {
    roomId: props.id,
    content: messageInput.value,
  }
  useCommSocket((socket) => socket.emit(CommEmitEvent.RoomMessage, data))

  messageInput.value = ''
  log.info(`Send event ${CommEmitEvent.RoomMessage}:`, data)
}

const boxChat = ref<HTMLDivElement>()
onMounted(() => {
  if (boxChat.value) {
    const observer = new MutationObserver(() => {
      boxChat.value!.scrollTop = boxChat.value!.scrollHeight
    })
    observer.observe(boxChat.value, { childList: true })
  }
})
</script>

<template>
  <div relative border rounded p-2>
    <div
      ref="boxChat"
      h="10/12"
      overflow-x-hidden
      overflow-y-scroll
      px-1
      text-left
    >
      <div v-for="(message, i) in messages" :key="i">
        <p mb-2 break-words>
          <b
            v-if="message.senderId"
            :class="message.senderId === player?.username ? 'text-blue' : ''"
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
        outlined
        dense
        autogrow
        placeholder="Say something..."
        h="max-[200px]"
        class="chatbox"
        @keydown.enter.prevent="sendMessage"
      />
    </div>
  </div>
</template>
