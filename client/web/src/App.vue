<script setup lang="ts">
import log from 'loglevel'
import { useQuasar } from 'quasar'

useHead({
  title: 'Vitesse',
  meta: [
    { name: 'description', content: 'Opinionated Vite Starter Template' },
    {
      name: 'theme-color',
      content: () => (isDark.value ? '#00aba9' : '#ffffff'),
    },
  ],
  link: [
    {
      rel: 'icon',
      type: 'image/svg+xml',
      href: () => (preferredDark.value ? '/favicon-dark.svg' : '/favicon.svg'),
    },
  ],
})

const $q = useQuasar()

onErrorCaptured((err) => {
  log.error(err)
  $q.notify({
    color: 'red',
    message: err?.message ? err.message : 'Unknown error',
  })
  return false
})
</script>

<template>
  <RouterView />
</template>
