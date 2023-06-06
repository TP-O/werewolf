<script setup lang="ts">
import { Field, useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as zod from 'zod'
import { useQuasar } from 'quasar'

defineOptions({
  name: 'SignUpPage',
})

const router = useRouter()
const $q = useQuasar()

const schema = zod
  .object({
    email: zod.string().email().default(''),
    password: zod
      .string()
      .nonempty('Password is required')
      .min(8, { message: 'Password must have at least 8 characters' })
      .default(''),
    confirmPassword: zod
      .string()
      .nonempty('Confirm password is required')
      .default(''),
  })
  .refine(({ password, confirmPassword }) => password === confirmPassword, {
    message: 'Password does not match',
    path: ['confirmPassword'],
  })
const { handleSubmit, useFieldModel } = useForm({
  validationSchema: toTypedSchema(schema),
})

const onSubmit = handleSubmit(async (values: zod.infer<typeof schema>) => {
  $q.loading.show({
    message: 'Creating account...',
  })
  await auth.signUp(values.email, values.password)
  $q.loading.hide()

  router.push('/')
})
</script>

<template>
  <div h-screen w-screen flex overflow-hidden>
    <div flex="~ col justify-center" w-full gap-8 px-8 md="px-24" lg="w-xl">
      <div>
        <div mb-8>
          <div text="5xl creepy" font-bold>New account</div>
        </div>

        <form flex="~ col justify-between" gap-2 @submit="onSubmit">
          <Field v-slot="{ errorMessage }" name="email">
            <q-input
              v-model="useFieldModel('email').value"
              outlined
              dense
              rounded
              type="text"
              label="Email"
              :error="!!errorMessage"
              :error-message="errorMessage"
            />
          </Field>

          <Field v-slot="{ errorMessage }" name="password">
            <q-input
              v-model="useFieldModel('password').value"
              outlined
              dense
              rounded
              type="password"
              label="Password"
              :error="!!errorMessage"
              :error-message="errorMessage"
            />
          </Field>

          <Field v-slot="{ errorMessage }" name="confirmPassword">
            <q-input
              v-model="useFieldModel('confirmPassword').value"
              outlined
              dense
              rounded
              type="password"
              label="Confirm password"
              :error="!!errorMessage"
              :error-message="errorMessage"
            />
          </Field>

          <div flex="~ justify-end" mb-3>
            <router-link to="sign-in" color="blue" font-bold>
              Sign in instead
            </router-link>
          </div>

          <q-btn
            color="blue"
            label="Go!"
            type="submit"
            rounded
            px-8
            py-2
            capitalize
          />
        </form>
      </div>

      <q-separator />

      <OAuth />
    </div>

    <div lg="block" hidden grow-1 bg-sign-in />
  </div>
</template>

<route lang="yaml">
meta:
  layout: guest
</route>
