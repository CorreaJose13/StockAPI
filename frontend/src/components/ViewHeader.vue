<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>
<template>
  <section class="flex flex-col gap-1 w-full py-4">
    <h1 class="text-3xl font-bold text-start text-white">{{ title }}</h1>
    <p class="text-start text-slate-200">{{ description }}</p>
  </section>
</template>
