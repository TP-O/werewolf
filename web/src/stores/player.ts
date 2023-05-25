import { acceptHMRUpdate, defineStore } from 'pinia'
import type { PlayerId } from '~/types'

interface Player {
  id: PlayerId
  username: string
}

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null | undefined>()
  const isAuthChecked = new Promise<boolean>((resolve) => {
    const stop = watch(player, () => {
      stop()
      resolve(true)
    })
  })

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

  /**
   * Wait until the frist time auth state is changed.
   *
   * @returns True if the player is authenticated.
   */
  async function checkAuth() {
    if (await isAuthChecked) {
      return player.value !== null
    }
  }

  return {
    player,
    checkAuth,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(usePlayerStore as any, import.meta.hot)
  )
}
