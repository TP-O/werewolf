import { FactionId } from '~/enums'

export const factions: Record<FactionId, { id: FactionId; name: string }> = {
  [FactionId.Villager]: {
    id: FactionId.Villager,
    name: 'Villager',
  },
  [FactionId.Werewolf]: {
    id: FactionId.Werewolf,
    name: 'Werewolf',
  },
}
