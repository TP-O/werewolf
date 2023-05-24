<script setup lang="ts">
import { useVuelidate } from '@vuelidate/core'
import { email, helpers, minLength, required } from '@vuelidate/validators'

defineOptions({
  name: 'SignUpPage',
})

const data = reactive({
  email: '',
  password: '',
  confirmPassword: '',
})
const schema = {
  email: {
    required: helpers.withMessage('Email is required', required),
    email: helpers.withMessage('Invalid email', email),
  },
  password: {
    required: helpers.withMessage('Password is required', required),
    minLength: helpers.withMessage(
      ({ $params }) => `Password must be at least ${$params.min} characters`,
      minLength(8)
    ),
  },
  confirmPassword: {
    sameAsPassword: helpers.withMessage(
      'Password does not match',
      (v) => v === data.password
    ),
  },
}
const form = useVuelidate(schema, data)
const router = useRouter()

async function onSubmit() {
  if (!(await form.value.$validate())) {
    return
  }

  await auth.signUp(data.email, data.password)
  router.push('/')
}
</script>

<template>
  <div flex="~ justify-center items-center" h-full>
    <div md="w-2xl">
      <div mb-8>
        <div text="2xl" font-bold>Create new account</div>
      </div>

      <form flex="~ col justify-between" gap-4 px-4 @submit.prevent="onSubmit">
        <q-input
          v-model="form.email.$model"
          :debounce="200"
          outlined
          type="email"
          label="Email"
          :error="form.email.$error"
          :error-message="form.email.$errors[0]?.$message.toString()"
        />
        <q-input
          v-model="form.password.$model"
          :debounce="200"
          outlined
          type="password"
          label="Password"
          :error="form.password.$error"
          :error-message="form.password.$errors[0]?.$message.toString()"
        />
        <q-input
          v-model="form.confirmPassword.$model"
          :debounce="200"
          outlined
          type="password"
          label="Confirm password"
          :error="form.confirmPassword.$error"
          :error-message="form.confirmPassword.$errors[0]?.$message.toString()"
        />
        <div flex="~ justify-between items-center">
          <router-link to="sign-in" color="blue" font-bold>
            Sign in instead
          </router-link>
          <q-btn color="blue" label="Go!" type="submit" px-8 py-2 capitalize />
        </div>
      </form>

      <q-separator my-6 />

      <div>
        <div mb-4>Or join with</div>
        <div flex="~ justify-around">
          <q-btn capitalize @click="auth.signInWithGoogle">
            <div i-devicon-google mr-2 />
            Google
          </q-btn>
          <q-btn capitalize @click="auth.signInWithGoogle">
            <div i-devicon-facebook mr-2 />
            Facebook
          </q-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: guest
</route>
