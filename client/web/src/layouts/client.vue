<script setup lang="ts">
import { storeToRefs } from 'pinia'

const router = useRouter()
const playerStore = usePlayerStore()
const { isCommServerConnected } = storeToRefs(useClientStore())

onMounted(async () => {
  await auth.waitForAuthState()

  if (!playerStore.player)
    return router.push('/auth/sign-in')

  await useCommSocket()
})
</script>

<template>
  <main
    h-screen p-2
    text="center gray-700 dark:gray-200"
  >
    <div
      v-if="!isCommServerConnected"
      h-full
      flex="~ col justify-center items-center"
    >
      <q-spinner-ball
        color="primary"
        size="6em"
        mb-6
      />
      <div text-xl>
        Connecting...
      </div>
    </div>
    <RouterView v-else />
  </main>
</template>
