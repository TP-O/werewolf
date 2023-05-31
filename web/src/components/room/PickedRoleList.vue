<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { roles } from '~/constants'
import type { RoleId } from '~/enums'

const { player } = storeToRefs(usePlayerStore())
const { room, settings } = storeToRefs(useWaitingRoomStore())
const isOwner = computed(() => room.value?.ownerId === player.value?.id)

const isRoleSelectionShowed = ref(false)
function showRoleSelection() {
  isRoleSelectionShowed.value = true
}

const isRoleDetailsShowed = ref(false)
const exploredRoleId = ref<RoleId>(1)
function showRoleDetails(id: number) {
  exploredRoleId.value = id
  isRoleDetailsShowed.value = true
}

function pickRole(id: number) {
  const i = settings.value.gameSettings.roleIds.indexOf(id)
  if (i === -1) {
    settings.value.gameSettings.roleIds.unshift(id)
  } else {
    settings.value.gameSettings.roleIds.splice(i, 1)

    const j = settings.value.gameSettings.requiredRoleIds.indexOf(id)
    if (j !== -1) {
      settings.value.gameSettings.requiredRoleIds.splice(j, 1)
    }
  }
}

function pickRequiredRole(id: number) {
  const i = settings.value.gameSettings.requiredRoleIds.indexOf(id)
  if (i === -1) {
    settings.value.gameSettings.requiredRoleIds.unshift(id)

    if (!settings.value.gameSettings.roleIds.includes(id)) {
      settings.value.gameSettings.roleIds.unshift(id)
    }
  } else {
    settings.value.gameSettings.requiredRoleIds.splice(i, 1)
  }
}
</script>

<template>
  <div flex="~ col" border rounded p-2>
    <div mb-2>Picked roles</div>

    <q-btn
      color="secondary"
      :label="isOwner ? 'Add role' : 'View roles'"
      mb-4
      w-full
      outline
      @click="showRoleSelection"
    />

    <div grid="~ cols-1" gap-4 overflow-y-scroll>
      <RoleSummaryCard
        v-for="(role, i) of settings.gameSettings.roleIds.map(
          (id) => roles[id]
        )"
        :key="i"
        :role="role"
        :picked="settings.gameSettings.roleIds.includes(role.id)"
        :required="settings.gameSettings.requiredRoleIds.includes(role.id)"
        :pick="pickRole"
        :mark-as-required="pickRequiredRole"
        @click="showRoleDetails(role.id)"
      />
    </div>

    <q-dialog
      v-model="isRoleSelectionShowed"
      transition-show="rotate"
      transition-hide="rotate"
      full-height
    >
      <q-card w="700px" max-w="80vw">
        <q-card-section>
          <div class="text-h6">Roles</div>
        </q-card-section>

        <q-card-section flex="~ col" gap-4 pt-0>
          <RoleSummaryCard
            v-for="(role, i) of roles"
            :key="i"
            :role="role"
            :picked="settings.gameSettings.roleIds.includes(role.id)"
            :required="settings.gameSettings.requiredRoleIds.includes(role.id)"
            :pick="pickRole"
            :mark-as-required="pickRequiredRole"
            @click="showRoleDetails(role.id)"
          />
        </q-card-section>
      </q-card>
    </q-dialog>

    <q-dialog
      v-model="isRoleDetailsShowed"
      transition-show="rotate"
      transition-hide="rotate"
    >
      <RoleCard :role="roles[exploredRoleId]" />
    </q-dialog>
  </div>
</template>
