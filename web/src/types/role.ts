import type { FactionId, RoleId } from '~/enums'

export interface Role {
  id: RoleId
  name: string
  factionId: FactionId
  sets: number
  description: string
}
