import { acceptHMRUpdate, defineStore } from 'pinia'

export const useClientStore = defineStore('client', () => {
  const connectedToCommunication = ref(false)

  function setConnectedToCommunication(isConneced: boolean) {
    connectedToCommunication.value = isConneced
  }

  return {
    connectedToCommunication,
    setConnectedToCommunication,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useClientStore as any, import.meta.hot))
}
