import { acceptHMRUpdate, defineStore } from 'pinia'

interface Player {
  username: string
}

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null>(null)

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

  let isLoaded = false
  function waitForLoadingPlayer() {
    return new Promise<void>((resolve) => {
      if (isLoaded)
        return resolve()

      isLoaded = true
      firebaseAuth.onAuthStateChanged(() => resolve())
    })
  }

  return {
    player,
    waitForLoadingPlayer,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(usePlayerStore as any, import.meta.hot))
