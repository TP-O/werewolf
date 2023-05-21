import { acceptHMRUpdate, defineStore } from 'pinia'

interface ErrorDialog {
  error: boolean
  message?: string
}

export const useDialogStore = defineStore('dialog', () => {
  const errorDialog = ref<ErrorDialog>({
    error: false,
  })

  const openErrorDialog = (message: string) => {
    errorDialog.value = {
      error: true,
      message,
    }
  }

  return {
    errorDialog,
    openErrorDialog,
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useDialogStore as any, import.meta.hot))
}
