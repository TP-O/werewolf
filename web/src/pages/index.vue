<script setup lang="ts">
defineOptions({
  name: 'HomePage',
})

const { player } = usePlayerStore()
const { book } = useWaitingRoomStore()
const router = useRouter()

function bookRoom() {
  book()
    .then(({ id }) => router.push(`/room/${id}`))
    .catch(() => {
      throw new Error('Unable to create room')
    })
}

async function joinRoom() {
  const id = prompt('Enter room ID:')
  if (id) {
    router.push(`/room/${id}`)
  }
}
</script>

<template>
  <div>Hello {{ player?.username }}</div>

  <q-btn
    color="blue"
    label="Sign out!"
    px-8
    py-2
    capitalize
    @click="auth.signOut"
  />
  <q-btn
    color="blue"
    label="Create room"
    px-8
    py-2
    capitalize
    @click="bookRoom"
  />
  <q-btn
    color="blue"
    label="Join room"
    px-8
    py-2
    capitalize
    @click="joinRoom"
  />
</template>

<route lang="yaml">
meta:
  layout: client
</route>
