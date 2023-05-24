import { acceptHMRUpdate, defineStore } from 'pinia'
import type { PlayerId } from '~/types'

interface Player {
  id: PlayerId
  username: string
}

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null>(null)

  auth.raw.onAuthStateChanged((user) => {
    if (user) {
      player.value = {
        id: user.uid,
        username: user.uid,
      }
    } else {
      player.value = null
    }
  })

  return {
    player,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(usePlayerStore as any, import.meta.hot)
  )
}
