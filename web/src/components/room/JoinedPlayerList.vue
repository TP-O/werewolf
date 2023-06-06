<script setup lang="ts">
import { storeToRefs } from 'pinia'

const { player } = storeToRefs(usePlayerStore())
const { room } = storeToRefs(useWaitingRoomStore())
const { kick } = useWaitingRoomStore()
</script>

<template>
  <div relative overflow-x-hidden overflow-y-scroll border rounded>
    <header-card label="players" p-1.15 />

    <q-card v-for="(memberId, i) of room?.memberIds" :key="i" m-2>
      <q-card-section horizontal justify-between>
        <q-item w="4/5">
          <q-item-section avatar>
            <q-avatar>
              <img
                :src="`https://picsum.photos/seed/${Math.random()}/200/200`"
              />
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label>
              <p w-full overflow-hidden text-ellipsis whitespace-nowrap>
                {{ memberId }}
              </p>
            </q-item-label>
          </q-item-section>
        </q-item>

        <q-card-actions>
          <q-btn-dropdown unelevated>
            <q-list>
              <q-item v-close-popup clickable>
                <q-item-section avatar>
                  <q-avatar color="black" text-color="white">
                    <q-icon>
                      <div i-mdi-account-circle-outline />
                    </q-icon>
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>Profile</q-item-label>
                </q-item-section>
              </q-item>

              <q-item v-close-popup clickable>
                <q-item-section avatar>
                  <q-avatar color="black" text-color="white">
                    <q-icon>
                      <div i-mdi-account-plus-outline />
                    </q-icon>
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>Add friend</q-item-label>
                </q-item-section>
              </q-item>

              <q-item
                v-if="room?.ownerId === player?.username"
                v-close-popup
                clickable
                @click="kick(memberId)"
              >
                <q-item-section avatar>
                  <q-avatar color="red" text-color="white">
                    <q-icon>
                      <div i-mdi-exit-run />
                    </q-icon>
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>Kick</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </q-card-actions>
      </q-card-section>
    </q-card>
  </div>
</template>
