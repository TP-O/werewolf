<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()
const { player } = storeToRefs(usePlayerStore())
const { checkAuth } = usePlayerStore()
const { room } = storeToRefs(useWaitingRoomStore())
const { isCommServerConnected } = storeToRefs(useClientStore())

/**
 * Add delay time for player observation.
 */
function deplayNavigate() {
  setTimeout(() => router.push('/auth/sign-in'), 1000)
}

watch(room, () => {
  if (!room.value) {
    router.push('/')
    $q.notify('You left the room')
  } else {
    router.push(`/room/${room.value.id}`)
  }
})

watch(player, async () => {
  if (!(await checkAuth())) {
    // router.push('/auth/sign-in')
    deplayNavigate()
  }
})

onMounted(async () => {
  if (await checkAuth()) {
    await useCommSocket()
  } else {
    // router.push('/auth/sign-in')
    deplayNavigate()
  }
})
</script>

<template>
  <main h-screen p-2 text="center gray-700 dark:gray-200">
    <LoadingContainer
      :loading="!isCommServerConnected"
      message="Connecting to communication server..."
    >
      <RouterView />
    </LoadingContainer>
  </main>
</template>
