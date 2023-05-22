<script setup lang="ts">
const route = useRoute()
const roomStore = useRoomStore()
const room = computed(() => roomStore.waitingRoom)
const roomId = String(route.params.id)
const messageStore = useMessageStore()
const messages = computed(() => messageStore.roomMessages[roomId])
const { player } = usePlayerStore()
const messageInput = ref('')

function sendMessage() {
  commSocket.sendRoomMessage({
    roomId,
    content: messageInput.value,
  })
}

onBeforeMount(async () => {
  if (!room.value || room.value.id !== roomId) {
    await roomStore.joinRoom(roomId)

    if (!room.value)
      throw new Error('Unable to join this room')
  }
})
</script>

<template>
  <div h-full flex="~ col">
    <div mb-4 text-xl>
      Room: <b>{{ roomId }}</b>
    </div>

    <div grid="~ cols-3" h-full grow-1 gap-1>
      <div flex="~ col" gap-1>
        <div h="1/2" overflow-scroll overflow-x-hidden border rounded p-2>
          <div>Joined players</div>
          <div v-for="memberId, i of room?.memberIds" :key="i" p-2>
            Player {{ memberId }}
          </div>
        </div>

        <div h="1/2" flex="~ col" border rounded p-2>
          <div grow-1 text-left>
            <q-chat-message
              v-for="message, i in messages"
              :key="i"
              :name="message.senderId === player?.username
                ? '' : message.senderId"
              :text="[`${message.content}`]"
              :sent="message.senderId === player?.username"
            />
          </div>

          <div h="max-[300px]">
            <q-input
              v-model="messageInput"
              autogrow dense outlined rounded
            >
              <template #after>
                <q-btn round dense flat icon="send" @click="sendMessage" />
              </template>
            </q-input>
          </div>
        </div>
      </div>
      <div border rounded p-1>
        Selected role
      </div>
      <div border rounded p-1>
        Settings
      </div>
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: client
</route>
