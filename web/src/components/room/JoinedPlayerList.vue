<script setup lang="ts">
import { storeToRefs } from 'pinia'

const { player } = storeToRefs(usePlayerStore())
const { room } = storeToRefs(useWaitingRoomStore())
const { kick } = useWaitingRoomStore()
</script>

<template>
  <div overflow-x-hidden overflow-y-scroll border rounded p-2>
    <div mb-2>Joined players</div>

    <div
      v-for="(memberId, i) of room?.memberIds"
      :key="i"
      flex="~  justify-between"
      my-2
      gap-2
      border
      p-2
    >
      <p flex="~ items-center" overflow-hidden text-ellipsis whitespace-nowrap>
        {{ memberId }}
      </p>
      <q-btn>
        <div i="carbon-overflow-menu-horizontal" />
        <q-menu min-w="100px" dense>
          <q-list>
            <q-item
              v-if="room?.ownerId === player?.username"
              v-close-popup
              clickable
              @click="kick(memberId)"
            >
              <q-item-section>Kick</q-item-section>
            </q-item>
            <q-item v-close-popup clickable>
              <q-item-section>Profile</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>
  </div>
</template>
