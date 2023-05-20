import { acceptHMRUpdate, defineStore } from 'pinia'

interface Player {
  username: string
}

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null>(null)
  const p = computed(() => player)

  firebaseAuth.onAuthStateChanged((user) => {
    if (user) {
      player.value = {
        username: user.uid,
      }
    }
    else {
      player.value = null
    }
  })

  return {
    player: p,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(usePlayerStore as any, import.meta.hot))
