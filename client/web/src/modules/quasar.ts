import { Quasar } from 'quasar'
import type { UserModule } from '~/types'

import '@quasar/extras/material-icons/material-icons.css'
import 'quasar/src/css/index.sass'

export const install: UserModule = ({ app }) => {
  app.use(Quasar, {
    plugins: {},
  })
}
