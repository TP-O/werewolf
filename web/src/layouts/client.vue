<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()
const { player } = storeToRefs(usePlayerStore())
const { checkAuth } = usePlayerStore()
const { room } = storeToRefs(useWaitingRoomStore())
const { isCommServerConnected } = storeToRefs(useClientStore())

watch(room, () => {
  if (!room.value) {
    router.push('/')
    $q.notify('You left the room')
  } else {
    router.push(`/room/${room.value.id}`)
  }
})

watchEffect(() => {
  if (isCommServerConnected.value) {
    $q.loading.hide()
  } else {
    $q.loading.show({
      message: 'Connecting to communication server...',
    })
  }
})

watch(player, async () => {
  if (!(await checkAuth())) {
    setTimeout(() => router.push('/auth/sign-in'), 1000)
  }
})

onMounted(async () => {
  if (await checkAuth()) {
    await useCommSocket()
  } else {
    setTimeout(() => router.push('/auth/sign-in'), 1000)
  }
})

onUnmounted(() => {
  useCommSocket((socket) => socket.disconnect())
})
</script>

<template>
  <main h-screen p-2 text="center gray-700 dark:gray-200">
    <RouterView />
  </main>
</template>
