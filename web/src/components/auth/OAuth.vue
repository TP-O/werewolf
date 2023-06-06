<script setup lang="ts">
import { useQuasar } from 'quasar'

const $q = useQuasar()

async function signInWithGoogle() {
  $q.loading.show({
    message: 'Accessing...',
  })
  await auth.signInWithGoogle()
  $q.loading.hide()
}
</script>

<template>
  <div flex="~ col justify-around" gap-4>
    <q-btn rounded @click="signInWithGoogle">
      <div i-devicon-google mr-2 />
      <div>Join with Google</div>
    </q-btn>

    <q-btn rounded @click="signInWithGoogle">
      <div i-devicon-facebook mr-2 />
      <div>Join with Facebook</div>
    </q-btn>
  </div>
</template>
