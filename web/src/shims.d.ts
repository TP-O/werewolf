import { AxiosInstance } from 'axios'

declare interface Window {
  // extend the window
}

declare module '*.vue' {
  import { type DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface ImportMetaEnv {
  [key: string]: any
  BASE_URL: string
  MODE: string
  DEV: boolean
  PROD: boolean
  SSR: boolean

  VITE_COMMUNICATION_SERVER: string
  VITE_GAME_SERVER: string
  VITE_FIREBASE_CONFIG: string
}
