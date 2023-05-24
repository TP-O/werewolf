<script setup lang="ts">
const { player } = usePlayerStore()
const router = useRouter()

onMounted(async () => {
  await auth.waitForAuthState()

  if (player) {
    router.push('/')
  }
})
</script>

<template>
  <main h-screen px-4 py-10 text="center gray-700 dark:gray-200">
    <RouterView />
  </main>
</template>
