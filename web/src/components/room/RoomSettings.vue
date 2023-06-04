<script setup lang="ts">
import { Field, useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as zod from 'zod'
import { storeToRefs } from 'pinia'

const { player } = storeToRefs(usePlayerStore())
const { room, settings } = storeToRefs(useWaitingRoomStore())
const isOwner = computed(() => room.value?.ownerId === player.value?.id)
const info = [
  {
    label: 'Room ID',
    value: room.value?.id,
  },
  {
    label: 'Capacity',
    value: `${settings.value.capacity} players`,
  },
  {
    label: 'Password',
    value: settings.value.password,
  },
  {
    label: 'Turn duration',
    value: `${settings.value.gameSettings.turnDuration} seconds`,
  },
  {
    label: 'Discussion duration',
    value: `${settings.value.gameSettings.discussionDuration} seconds`,
  },
]

const schema = zod.object({
  password: zod
    .string()
    .optional()
    .refine((v) => !v || (v.length >= 4 && v.length <= 15), {
      message: 'Password must be between 4 and 15 characters',
    }),
  capacity: zod.number().min(5).max(20),
  gameSettings: zod.object({
    turnDuration: zod.number().min(30).max(60),
    discussionDuration: zod.number().min(60).max(300),
  }),
})
const { handleSubmit, useFieldModel, setValues } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: settings,
})

const isSettingsShowed = ref(false)
function showSettings() {
  setValues(settings.value)
  isSettingsShowed.value = true
}

const onSubmit = handleSubmit((d) => {
  console.log(d)
  isSettingsShowed.value = false
})
</script>

<template>
  <div flex="~ col" relative border rounded>
    <HeaderCard label="settings">
      <q-btn unelevated text-base @click="showSettings">
        <div i="carbon-settings"></div>

        <q-tooltip text-base> Update settings </q-tooltip>
      </q-btn>
    </HeaderCard>

    <div overflow-y-scroll p-2>
      <table>
        <tbody>
          <tr v-for="(row, i) of info" :key="i">
            <td pb-2 font-bold>{{ row.label }}:</td>
            <td pb-2 pl-4>{{ row.value }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <q-dialog
      v-model="isSettingsShowed"
      transition-show="rotate"
      transition-hide="rotate"
    >
      <q-card w="700px" max-w="80vw">
        <q-card-section>
          <div class="text-h6">Settings</div>
        </q-card-section>

        <q-card-section>
          <form grid gap-4 @submit="onSubmit">
            <Field v-slot="{ errorMessage }" name="password">
              <q-input
                v-model="useFieldModel('password').value"
                outlined
                :disable="!isOwner"
                label="Password"
                :error="!!errorMessage"
                :error-message="errorMessage"
              />
            </Field>

            <Field v-slot="{ errorMessage }" name="capacity">
              <q-input
                v-model.number="useFieldModel('capacity').value"
                outlined
                type="number"
                min="5"
                max="20"
                :disable="!isOwner"
                label="Capacity"
                :error="!!errorMessage"
                :error-message="errorMessage"
              />
            </Field>

            <Field v-slot="{ errorMessage }" name="gameSettings.turnDuration">
              <q-input
                v-model.number="
                  useFieldModel('gameSettings.turnDuration').value
                "
                outlined
                type="number"
                min="20"
                max="60"
                :disable="!isOwner"
                label="Turn duration"
                :error="!!errorMessage"
                :error-message="errorMessage"
              />
            </Field>

            <Field
              v-slot="{ errorMessage }"
              name="gameSettings.discussionDuration"
            >
              <q-input
                v-model.number="
                  useFieldModel('gameSettings.discussionDuration').value
                "
                outlined
                type="number"
                min="40"
                max="360"
                :disable="!isOwner"
                label="Discussion duration"
                :error="!!errorMessage"
                :error-message="errorMessage"
              />
            </Field>

            <q-btn
              type="submit"
              color="secondary"
              label="Save"
              w-full
              outline
            />
          </form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>
