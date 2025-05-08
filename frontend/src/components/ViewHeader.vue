<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>
<template>
  <section class="flex w-full flex-col gap-1 py-4">
    <h2 class="text-start text-3xl font-bold text-white">{{ title }}</h2>
    <p class="text-start text-slate-200">{{ description }}</p>
  </section>
</template>
