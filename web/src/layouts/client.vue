<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()
const { player } = usePlayerStore()
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

onMounted(async () => {
  await auth.waitForAuthState()

  if (!player) {
    return router.push('/auth/sign-in')
  } else {
    await useCommSocket()
  }
})
</script>

<template>
  <main h-screen p-2 text="center gray-700 dark:gray-200">
    <div
      v-if="!isCommServerConnected"
      h-full
      flex="~ col justify-center items-center"
    >
      <q-spinner-ball color="primary" size="6em" mb-6 />
      <div text-xl>Connecting...</div>
    </div>
    <RouterView v-else />
  </main>
</template>
