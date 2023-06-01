<script setup lang="ts">
import { useQuasar } from 'quasar'
import { Field, useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as zod from 'zod'

defineOptions({
  name: 'SignInPage',
})

const router = useRouter()
const $q = useQuasar()

const schema = zod.object({
  email: zod.string().email().default(''),
  password: zod
    .string()
    .nonempty('Password is required')
    .min(8, { message: 'Password must have at least 8 characters' })
    .default(''),
})
const { handleSubmit, useFieldModel } = useForm({
  validationSchema: toTypedSchema(schema),
})

const onSubmit = handleSubmit(async (values: zod.infer<typeof schema>) => {
  $q.loading.show({
    message: 'Accessing...',
  })
  await auth.signIn(values.email, values.password)
  $q.loading.hide()

  router.push('/')
})
</script>

<template>
  <div flex="~ justify-center items-center" h-full>
    <div md="w-2xl">
      <div mb-8>
        <div text="2xl" font-bold>Join to player with your friends</div>
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

        <div flex="~ justify-between items-center">
          <router-link to="sign-up" color="blue" font-bold>
            Create account
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
