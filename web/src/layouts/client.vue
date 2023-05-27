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
  <main h-screen text="center gray-700 dark:gray-200">
    <q-splitter :model-value="80" h-full>
      <template #before>
        <div h-full p-2>
          <RouterView />
        </div>
      </template>

      <template #after>
        <div h-full flex="~ col">
          <div flex="~ row items-center" gap-2 p-2>
            <q-avatar border="1 red-500" box-content p-1>
              <img src="https://cdn.quasar.dev/img/avatar2.jpg" />
            </q-avatar>
            <div min-w-0>
              <p overflow-hidden text-ellipsis font-bold>
                {{ player?.username }}
              </p>
            </div>
          </div>

          <q-separator />

          <div grow-1 overflow-y-scroll>
            <div v-for="i in 10" :key="i" flex="~ row items-center" gap-2 p-2>
              <q-avatar border="1 red-500" box-content p-1 size="md">
                <img src="https://cdn.quasar.dev/img/avatar2.jpg" />
              </q-avatar>
              <div min-w-0>
                <p text-md overflow-hidden text-ellipsis>
                  {{ player?.username }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </q-splitter>
  </main>
</template>
