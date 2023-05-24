import { acceptHMRUpdate, defineStore } from 'pinia'
import { useQuasar } from 'quasar'
import type { ErrorResponse, SuccessResponse } from '~/types'

export const useClientStore = defineStore('client', () => {
  const $q = useQuasar()
  const isCommServerConnected = ref(false)

  function onConnectComm() {
    isCommServerConnected.value = true
  }

  function onErroConnectComm() {
    $q.notify({
      color: 'red',
      message: 'Unable to connect to communication server',
    })
  }

  function onDisconnectComm() {
    isCommServerConnected.value = false
    $q.notify({
      color: 'red',
      message: 'Disconnected from communication server',
    })
  }

  function onCommError(res: ErrorResponse) {
    if (typeof res.message === 'string') {
      $q.notify({
        color: 'red',
        message: res.message,
      })
    } else {
      res.message.forEach((m) =>
        $q.notify({
          color: 'red',
          message: m,
        })
      )
    }
  }

  function onCommSuccess(res: SuccessResponse) {
    $q.notify({
      color: 'green',
      message: res.message,
    })
  }

  return {
    isCommServerConnected,
    onConnectComm,
    onErroConnectComm,
    onDisconnectComm,
    onCommError,
    onCommSuccess,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useClientStore as any, import.meta.hot)
  )
}
