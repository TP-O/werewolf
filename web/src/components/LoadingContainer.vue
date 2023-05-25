<script setup lang="ts">
defineProps({
  loading: {
    type: Boolean,
    required: true,
  },
  message: {
    type: String,
    default: 'Connecting...',
  },
})
</script>

<template>
  <div relative h-full w-full>
    <slot />

    <q-inner-loading :showing="loading">
      <q-spinner-ball color="primary" size="6em" mb-6 />
      <div text-xl>{{ message }}</div>
    </q-inner-loading>
  </div>
</template>
