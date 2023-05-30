const a = 1

<script setup lang="ts">
import { factions } from '~/constants'
import { RoleId } from '~/enums'
import type { Role } from '~/types'

withDefaults(
  defineProps<{
    role: Role
    picked: boolean
    required: boolean
    pick: (id: Role['id']) => void
    markAsRequired: (id: Role['id']) => void
  }>(),
  { picked: false, required: false }
)
</script>

<template>
  <q-card flat bordered cursor-pointer :border="picked ? 'green-400' : 'gray'">
    <q-card-section horizontal justify-between>
      <q-item>
        <q-item-section avatar>
          <q-avatar>
            <img src="https://cdn.quasar.dev/img/avatar2.jpg" />
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
        <q-btn
          flat
          round
          :color="required ? 'green' : ''"
          @click.stop="markAsRequired(role.id)"
        >
          <div i="carbon-manage-protection"></div>
        </q-btn>

        <q-btn v-if="picked" flat round color="red" @click.stop="pick(role.id)">
          <div i="carbon-close-filled"></div>
        </q-btn>
        <q-btn v-else flat round color="green" @click.stop="pick(role.id)">
          <div i="carbon-checkmark-filled"></div>
        </q-btn>
      </q-card-actions>
    </q-card-section>
  </q-card>
</template>
