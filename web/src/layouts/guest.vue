<script setup lang="ts">
import { storeToRefs } from 'pinia'

const router = useRouter()
const { player } = storeToRefs(usePlayerStore())
const { checkAuth } = usePlayerStore()
const isMounted = ref(false)

/**
 * Add delay time for player observation.
 */
function deplayNavigate() {
  setTimeout(() => router.push('/'), 1000)
}

watch(player, async () => {
  if (await checkAuth()) {
    deplayNavigate()
  }
})

onMounted(async () => {
  if (await checkAuth()) {
    deplayNavigate()
  } else {
    isMounted.value = true
  }
})
</script>

<template>
  <main h-screen px-4 py-10 text="center gray-700 dark:gray-200">
    <LoadingContainer
      :loading="!isMounted"
      message="Checking authentication..."
    >
      <RouterView />
    </LoadingContainer>
  </main>
</template>
