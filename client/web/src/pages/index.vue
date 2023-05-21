<script setup lang="ts">
defineOptions({
  name: 'HomePage',
})
const playerStore = usePlayerStore()
const roomStore = useRoomStore()
const router = useRouter()

async function bookRoom() {
  await roomStore.bookRoom()
  router.push(`/room/${roomStore.waitingRoom?.id}`)
}

async function joinRoom() {
  const id = prompt('Enter room ID:')

  if (id) {
    await roomStore.joinRoom(id)
    router.push(`/room/${roomStore.waitingRoom?.id}`)
  }
}
</script>

<template>
  <div>
    Hello {{ playerStore.player?.username }}
  </div>

  <q-btn color="blue" label="Sign out!" px-8 py-2 capitalize @click="signOut" />
  <q-btn color="blue" label="Create room" px-8 py-2 capitalize @click="bookRoom" />
  <q-btn color="blue" label="Join room" px-8 py-2 capitalize @click="joinRoom" />
</template>

<route lang="yaml">
meta:
  layout: client
</route>
