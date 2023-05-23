import log from 'loglevel'
import type { UserModule } from '~/types'

export const install: UserModule = ({ isClient }) => {
  if (!isClient)
    return

  import.meta.env.DEV ? log.setLevel('trace') : log.setLevel('silent')
}
