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
  <div>
    <div mb-4>Or join with</div>
    <div flex="~ justify-around">
      <q-btn capitalize @click="signInWithGoogle">
        <div i-devicon-google mr-2 />
        Google
      </q-btn>
      <q-btn capitalize @click="signInWithGoogle">
        <div i-devicon-facebook mr-2 />
        Facebook
      </q-btn>
    </div>
  </div>
</template>
