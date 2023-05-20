<script setup lang="ts">
defineOptions({
  name: 'HomePage',
})
const { player } = usePlayerStore()
</script>

<template>
  <div>
    hello {{ player?.username }}
  </div>

  <q-btn color="blue" label="Sign ou!" px-8 py-2 capitalize @click="signOut" />
</template>

<route lang="yaml">
meta:
  layout: client
</route>
