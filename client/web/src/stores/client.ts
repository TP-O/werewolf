import { acceptHMRUpdate, defineStore } from 'pinia'
import { useQuasar } from 'quasar'
import type {
  ErrorResponse,
  SuccessResponse,
} from '~/composables/socketio/communication/types'

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

  function onCommError(event: ErrorResponse) {
    if (typeof event.message === 'string') {
      $q.notify({
        color: 'red',
        message: event.message,
      })
    }
    else {
      event.message.forEach(m => $q.notify({
        color: 'red',
        message: m,
      }))
    }
  }

  function onCommSuccess(event: SuccessResponse) {
    $q.notify({
      color: 'green',
      message: event.message,
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
    acceptHMRUpdate(useClientStore as any, import.meta.hot))
}
