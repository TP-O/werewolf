import { acceptHMRUpdate, defineStore } from 'pinia'

interface Player {
  username: string
}

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null>(null)

  auth.raw.onAuthStateChanged((user) => {
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
    player,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(usePlayerStore as any, import.meta.hot))
}
