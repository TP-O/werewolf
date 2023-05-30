import type { Role } from '~/types'
import { FactionId, RoleId } from '~/enums'

export const roles: Record<Role['id'], Role> = {
  [RoleId.Villager]: {
    id: RoleId.Villager,
    name: 'Villager',
    factionId: FactionId.Villager,
    description:
      'Every participant who is not the moderator, seer, medic, or werewolf is a villager. Think of these characters as potential food for the werewolf. Their only goal is to figure out who the werewolf is and stop the feeding frenzy.',
    sets: -1,
  },
  [RoleId.Werewolf]: {
    id: RoleId.Werewolf,
    name: 'Werewolf',
    factionId: FactionId.Werewolf,
    description:
      'The werewolf’s only job is to stalk the villagers at night and kill them without getting caught. This role probably has the most power in the game. Each night, the moderator wakes the werewolf up and asks them who they want to eat.',
    sets: -1,
  },
  [RoleId.Hunter]: {
    id: RoleId.Hunter,
    name: 'Hunter',
    factionId: FactionId.Villager,
    description:
      'The Hunter is a Villager that chooses a player to kill whenever the Hunter is killed. In a no-reveal game the Hunter tells the Moderator who to kill on the following night before anyone has woken up.',
    sets: 1,
  },
  [RoleId.Seer]: {
    id: RoleId.Seer,
    name: 'Seer',
    factionId: FactionId.Villager,
    description:
      'The seer is the anti-werewolf. This character is attempting to catch the werewolf and save the villagers. The seer has the chance to discover the werewolf, and then try to persuade other players.',
    sets: 1,
  },
  [RoleId.TwoSisters]: {
    id: RoleId.TwoSisters,
    name: 'Two Sisters',
    factionId: FactionId.Villager,
    description:
      'The Hunter is a Villager that chooses a player to kill whenever the Hunter is killed. In a no-reveal game the Hunter tells the Moderator who to kill on the following night before anyone has woken up.',
    sets: 2,
  },
}
