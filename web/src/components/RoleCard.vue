<script setup lang="ts">
import type { Role } from '~/types'

defineProps<{ role: Role }>()
</script>

<template>
  <q-card w="400px" h="500px" flex="~ col" gap-4 p-4>
    <q-img
      :src="`https://picsum.photos/seed/${Math.random()}/200/200`"
      fit="contain"
      h="1/2"
    />

    <q-card-section p-0>{{ role.description }} </q-card-section>
  </q-card>
</template>
