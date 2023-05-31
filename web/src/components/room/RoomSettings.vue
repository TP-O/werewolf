<script setup lang="ts">
import { storeToRefs } from 'pinia'

const { player } = storeToRefs(usePlayerStore())
const { room, settings } = storeToRefs(useWaitingRoomStore())
const isOwner = computed(() => room.value?.ownerId === player.value?.id)
</script>

<template>
  <div flex="~ col" relative border rounded p-1>
    <div mb-2>Settings</div>

    <div grid="~ cols-2" mb-4 gap-4>
      <q-input
        v-model="settings.password"
        outlined
        :disable="!isOwner"
        label="Password"
      />
      <q-input
        v-model="settings.capacity"
        outlined
        type="number"
        min="5"
        max="20"
        :disable="!isOwner"
        label="Capacity"
      />
      <q-input
        v-model="settings.gameSettings.turnDuration"
        outlined
        type="number"
        min="20"
        max="60"
        :disable="!isOwner"
        label="Turn duration"
      />
      <q-input
        v-model="settings.gameSettings.discussionDuration"
        outlined
        type="number"
        min="40"
        max="360"
        :disable="!isOwner"
        label="Discussion duration"
      />
    </div>

    <q-carousel
      v-model="settings.gameSettings.mapId"
      animated
      arrows
      navigation
      infinite
    >
      <q-carousel-slide
        :name="1"
        img-src="https://cdn.quasar.dev/img/mountains.jpg"
      />
      <q-carousel-slide
        :name="2"
        img-src="https://cdn.quasar.dev/img/parallax1.jpg"
      />
      <q-carousel-slide
        :name="3"
        img-src="https://cdn.quasar.dev/img/parallax2.jpg"
      />
      <q-carousel-slide
        :name="4"
        img-src="https://cdn.quasar.dev/img/quasar.jpg"
      />
    </q-carousel>

    <div absolute bottom-0 left-0 w-full p-2>
      <q-btn color="secondary" label="Save" w-full outline />
    </div>
  </div>
</template>
