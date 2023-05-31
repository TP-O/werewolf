<script setup lang="ts">
import useVuelidate from '@vuelidate/core'
import {
  helpers,
  minLength,
  minValue,
  required,
  requiredIf,
} from '@vuelidate/validators'
import { storeToRefs } from 'pinia'

const { player } = storeToRefs(usePlayerStore())
const { room, settings } = storeToRefs(useWaitingRoomStore())
const isOwner = computed(() => room.value?.ownerId === player.value?.id)

const schema = {
  capacity: {
    required: helpers.withMessage('Capacity is required', required),
    min: helpers.withMessage('Invalid capacity', minValue(5)),
  },
  password: {
    requiredIf: requiredIf(settings.value.password !== ''),
    minLength: helpers.withMessage(
      ({ $params }) => `Password must be at least ${$params.min} characters`,
      minLength(4)
    ),
  },
  gameSettings: {
    turnDuration: {
      required: helpers.withMessage('Capacity is required', required),
      min: helpers.withMessage('Invalid capacity', minValue(30)),
    },
    discussionDuration: {
      required: helpers.withMessage('Capacity is required', required),
      min: helpers.withMessage('Invalid capacity', minValue(60)),
    },
    mapId: {
      required: helpers.withMessage('Capacity is required', required),
      min: helpers.withMessage('Invalid capacity', minValue(1)),
    },
  },
}
const form = useVuelidate(schema, settings.value)
</script>

<template>
  <div flex="~ col" relative border rounded p-1>
    <div mb-2>Settings</div>

    <div grid="~ cols-2" mb-4 gap-4>
      <q-input
        v-model="form.password.$model"
        outlined
        :disable="!isOwner"
        label="Password"
        :error="form.password.$error"
        :error-message="form.password.$errors[0]?.$message.toString()"
      />
      <q-input
        v-model="form.capacity.$model"
        outlined
        type="number"
        min="5"
        max="20"
        :disable="!isOwner"
        label="Capacity"
        :error="form.capacity.$error"
        :error-message="form.capacity.$errors[0]?.$message.toString()"
      />
      <q-input
        v-model="form.gameSettings.turnDuration.$model"
        outlined
        type="number"
        min="20"
        max="60"
        :disable="!isOwner"
        label="Turn duration"
        :error="form.gameSettings.turnDuration.$error"
        :error-message="
          form.gameSettings.turnDuration.$errors[0]?.$message.toString()
        "
      />
      <q-input
        v-model="form.gameSettings.discussionDuration.$model"
        outlined
        type="number"
        min="40"
        max="360"
        :disable="!isOwner"
        label="Discussion duration"
        :error="form.gameSettings.discussionDuration.$error"
        :error-message="
          form.gameSettings.discussionDuration.$errors[0]?.$message.toString()
        "
      />
    </div>

    <q-carousel
      v-model="form.gameSettings.mapId.$model"
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
