<script setup lang="ts">
defineProps<{ label: string }>()
</script>

<template>
  <q-card flat>
    <q-card-section horizontal justify-between>
      <q-item w="4/5">
        <q-item-section>
          <q-item-label text-base capitalize> {{ label }} </q-item-label>
        </q-item-section>
      </q-item>

      <q-card-actions>
        <slot />
      </q-card-actions>
    </q-card-section>
  </q-card>
</template>
