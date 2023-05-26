<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()
const { player } = storeToRefs(usePlayerStore())
const { checkAuth } = usePlayerStore()
const { room } = storeToRefs(useWaitingRoomStore())
const { isCommServerConnected } = storeToRefs(useClientStore())

watch(
  isCommServerConnected,
  () => {
    if (isCommServerConnected.value) {
      $q.loading.hide()
    } else {
      $q.loading.show({
        message: 'Connecting to communication server...',
      })
    }
  },
  { immediate: true }
)

watch(
  player,
  async () => {
    if (!(await checkAuth())) {
      setTimeout(() => router.push('/auth/sign-in'), 1000)
    }
  },
  {
    immediate: true,
  }
)

watch(room, () => {
  if (!room.value) {
    router.push('/')
    $q.notify('You left the room')
  } else {
    router.push(`/room/${room.value.id}`)
  }
})

onBeforeMount(async () => {
  if (await checkAuth()) {
    await useCommSocket()
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
