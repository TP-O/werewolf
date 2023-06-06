<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'

const router = useRouter()
const $q = useQuasar()
const { player } = storeToRefs(usePlayerStore())
const { checkAuth } = usePlayerStore()

watch(
  player,
  async () => {
    if (await checkAuth()) {
      setTimeout(() => router.push('/'), 1000)
    } else {
      $q.loading.hide()
    }
  },
  {
    immediate: true,
  }
)

onBeforeMount(() => {
  $q.loading.show({
    message: 'Checking authentication...',
  })
})

onUnmounted(() => {
  if ($q.loading.isActive) {
    $q.loading.hide()
  }
})
</script>

<template>
  <main text="center gray-700 dark:gray-200">
    <RouterView />
  </main>
</template>
