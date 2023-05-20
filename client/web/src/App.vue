<script setup lang="ts">
import type { ErrorAlert } from '~/types'

// https://github.com/vueuse/head
// you can use this to manipulate the document head in any components,
// they will be rendered correctly in the html results with vite-ssg
useHead({
  title: 'Vitesse',
  meta: [
    { name: 'description', content: 'Opinionated Vite Starter Template' },
    {
      name: 'theme-color',
      content: () => isDark.value ? '#00aba9' : '#ffffff',
    },
  ],
  link: [
    {
      rel: 'icon',
      type: 'image/svg+xml',
      href: () => preferredDark.value ? '/favicon-dark.svg' : '/favicon.svg',
    },
  ],
})

const errorAlert = ref<ErrorAlert>({
  error: false,
})

onErrorCaptured((err) => {
  errorAlert.value = {
    error: true,
    message: err.message,
  }
  return false
})
</script>

<template>
  <RouterView />

  <q-dialog v-model="errorAlert.error">
    <q-card w-screen md="min-w-md">
      <q-card-section>
        <div text-xl font-bold color="error">
          Error!!!
        </div>
      </q-card-section>

      <q-card-section class="q-pt-none">
        {{ errorAlert.message }}
      </q-card-section>

      <q-card-actions align="right">
        <q-btn v-close-popup flat label="OK" color="primary" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>
