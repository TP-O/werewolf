<script setup lang="ts">
import { io } from 'socket.io-client'

const store = usePlayerStore()
const router = useRouter()

onMounted(async () => {
  await store.waitForLoadingPlayer()

  if (!store.player)
    return router.push('/auth/sign-in')

  const token = await firebaseAuth.currentUser?.getIdToken()
  const socket = io('http://127.0.0.1:8079/', {
    extraHeaders: {
      authorization: `Bearer ${token}`,
    },
    withCredentials: true,
  })

  socket.on('error', (e) => {
    console.log(e)
  })

  socket.on('error', (e) => {
    console.log(e)
  })

  socket.on('success', (e) => {
    console.log(e)
  })

  socket.on('friend_status', (e) => {
    console.log(e)
  })

  socket.on('private_message', (e) => {
    console.log(e)
  })

  socket.on('room_message', (e) => {
    console.log(e)
  })

  socket.on('room_change', (e) => {
    console.log(e)
  })

  axiosClient.post('/rooms/book', { password: '12345' })
})
</script>

<template>
  <main
    h-screen px-4 py-10
    text="center gray-700 dark:gray-200"
  >
    <RouterView />
  </main>
</template>
