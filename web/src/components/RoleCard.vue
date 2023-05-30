<script setup lang="ts">
import type { Role } from '~/types'

defineProps<{ role: Role }>()
</script>

<template>
  <q-card w="700px" max-w="80vw" h="500px" max-h="80vh" flex="~ col" gap-4 p-4>
    <q-img src="https://cdn.quasar.dev/img/parallax2.jpg" fit="contain" />

    <q-card-section p-0>{{ role.description }} </q-card-section>
  </q-card>
</template>
