import { type ViteSSGContext } from 'vite-ssg'

export type UserModule = (ctx: ViteSSGContext) => void

export interface ErrorAlert {
  error: boolean
  message?: string
}
