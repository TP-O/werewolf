<script setup lang="ts">
const playerStore = usePlayerStore()
const router = useRouter()

const clientStore = useClientStore()
const isCommServerConnected = computed(() => clientStore.isCommServerConnected)

onMounted(async () => {
  await auth.waitForAuthState()

  if (!playerStore.player)
    return router.push('/auth/sign-in')

  await commSocket.connect()
})
</script>

<template>
  <main
    h-screen p-4
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
