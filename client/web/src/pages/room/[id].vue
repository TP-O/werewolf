<script setup lang="ts">
const route = useRoute()
const roomStore = useRoomStore()
const room = computed(() => roomStore.waitingRoom)
const roomId = String(route.params.id)
const messageStore = useMessageStore()
const messages = computed(() => messageStore.roomMessages[roomId])
const { player } = usePlayerStore()
const content = ref('')

function send() {
  commSocket.sendRoomMessage({
    roomId,
    content: content.value,
  })
}

onBeforeMount(async () => {
  if (!room.value || room.value.id !== roomId) {
    await roomStore.joinRoom(roomId)

    if (!room.value)
      console.log('Cant join')
  }
})
</script>

<template>
  <div h-full flex>
    <div w="1/3" p-2>
      <q-splitter
        :model-value="50"
        horizontal
        h-full
      >
        <template #before>
          <div v-for="memberId, i of room?.memberIds" :key="i" p-2>
            Player {{ memberId }}
          </div>
        </template>

        <template #after>
          <div flex="~ col" h-full>
            <div grow-1>
              <q-chat-message
                v-for="message, i in messages"
                :key="i"
                :name="message.senderId === player?.username
                  ? 'me' : message.senderId"
                :text="[`${message.content}`]"
                :sent="message.senderId === player?.username"
              />
            </div>

            <q-input v-model="content" placeholder="Say something" />
            <q-btn label="Send" @click="send" />
          </div>
        </template>
      </q-splitter>
    </div>
    <div w="1/3" p-2>
      Selected role
    </div>
    <div w="1/3" p-2>
      Settings
    </div>
  </div>
</template>
