<script setup lang="ts">
defineOptions({
  name: 'HomePage',
})
const playerStore = usePlayerStore()
const { book, join } = useWaitingRoomStore()
const router = useRouter()

async function bookRoom() {
  const room = await book()
  router.push(`/room/${room.id}`)
}

async function joinRoom() {
  const id = prompt('Enter room ID:')

  if (id) {
    const room = await join(id)
    if (room)
      router.push(`/room/${room.id}`)
  }
}
</script>

<template>
  <div>
    Hello {{ playerStore.player?.username }}
  </div>

  <q-btn
    color="blue" label="Sign out!" px-8 py-2 capitalize @click="auth.signOut"
  />
  <q-btn
    color="blue" label="Create room" px-8 py-2 capitalize @click="bookRoom"
  />
  <q-btn
    color="blue" label="Join room" px-8 py-2 capitalize @click="joinRoom"
  />
</template>

<route lang="yaml">
meta:
  layout: client
</route>
