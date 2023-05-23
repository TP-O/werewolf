<script setup lang="ts">
import { error } from 'loglevel'
import { useQuasar } from 'quasar'

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

const $q = useQuasar()

onErrorCaptured((err) => {
  error(err)
  $q.notify({
    color: 'red',
    message: err.message,
  })
  return false
})
</script>

<template>
  <RouterView />
</template>
