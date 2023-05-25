import type { QLoadingShowOptions, QNotifyCreateOptions } from 'quasar'
import { Dialog, Loading, Notify, QSpinnerBall, Quasar } from 'quasar'
import type { UserModule } from '~/types'

import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/src/css/index.sass'

export const install: UserModule = ({ app }) => {
  app.use(Quasar, {
    plugins: {
      Dialog,
      Notify,
      Loading,
    },
    config: {
      notify: {
        position: 'top-right',
        timeout: 2500,
        textColor: 'white',
        progress: true,
        actions: [{ icon: 'close', color: 'white' }],
      } as QNotifyCreateOptions,
      loading: {
        spinner: QSpinnerBall,
        backgroundColor: 'white',
        spinnerColor: 'primary',
        messageColor: 'black',
        customClass: 'text-body1',
      } as QLoadingShowOptions,
    },
  })
}
