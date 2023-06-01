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
  <div flex="~ justify-center items-center" h-full>
    <div md="w-2xl">
      <div mb-8>
        <div text="2xl" font-bold>Create new account</div>
      </div>

      <form flex="~ col justify-between" gap-4 px-4 @submit="onSubmit">
        <Field v-slot="{ errorMessage }" name="email">
          <q-input
            v-model="useFieldModel('email').value"
            outlined
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
            type="password"
            label="Confirm password"
            :error="!!errorMessage"
            :error-message="errorMessage"
          />
        </Field>

        <div flex="~ justify-between items-center">
          <router-link to="sign-in" color="blue" font-bold>
            Sign in instead
          </router-link>
          <q-btn color="blue" label="Go!" type="submit" px-8 py-2 capitalize />
        </div>
      </form>

      <q-separator my-6 />

      <OAuth />
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: guest
</route>
