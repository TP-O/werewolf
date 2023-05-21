<script setup lang="ts">
const playerStore = usePlayerStore()
const router = useRouter()

onMounted(async () => {
  await playerStore.waitForLoadingPlayer()

  if (!playerStore.player)
    return router.push('/auth/sign-in')

  await connectCommunicationServer()
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
