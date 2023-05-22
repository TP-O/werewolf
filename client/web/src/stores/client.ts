import { acceptHMRUpdate, defineStore } from 'pinia'

export const useClientStore = defineStore('client', () => {
  const isCommServerConnected = ref(false)

  function setIsCommServerConnected(isConneced: boolean) {
    isCommServerConnected.value = isConneced
  }

  return {
    isCommServerConnected,
    setIsCommServerConnected,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useClientStore as any, import.meta.hot))
}
