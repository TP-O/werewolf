const a = 1

<script setup lang="ts">
import { factions } from '~/constants'
import { RoleId } from '~/enums'
import type { Role } from '~/types'

withDefaults(
  defineProps<{
    role: Role
    picked: boolean
    locked: boolean
    disabled: boolean
    pick: (id: Role['id']) => void
    markAsRequired: (id: Role['id']) => void
  }>(),
  { picked: false, required: false, disabled: false }
)
</script>

<template>
  <q-card cursor-pointer>
    <q-card-section horizontal justify-between>
      <q-item>
        <q-item-section avatar>
          <q-avatar>
            <img :src="`https://picsum.photos/seed/${Math.random()}/200/200`" />
          </q-avatar>
        </q-item-section>

        <q-item-section>
          <q-item-label>{{ role.name }}</q-item-label>
          <q-item-label caption>{{
            factions[role.factionId].name
          }}</q-item-label>
        </q-item-section>
      </q-item>

      <q-card-actions
        v-if="![RoleId.Villager, RoleId.Werewolf].includes(role.id)"
      >
        <q-btn flat :disabled="disabled" @click.stop="markAsRequired(role.id)">
          <div v-if="locked" i-mdi-lock-outline></div>
          <div v-else i-mdi-lock-open-variant-outline></div>

          <q-tooltip text-xs>
            {{
              locked
                ? 'This role is chosen by default'
                : 'This role is chosen at random'
            }}
          </q-tooltip>
        </q-btn>

        <q-btn flat :disabled="disabled" @click.stop="pick(role.id)">
          <div v-if="picked" i-mdi-close />
          <div v-else i-mdi-plus></div>

          <q-tooltip text-xs>
            {{ picked ? 'Delete role' : 'Add role' }}
          </q-tooltip>
        </q-btn>
      </q-card-actions>
    </q-card-section>
  </q-card>
</template>
