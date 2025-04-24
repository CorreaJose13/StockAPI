<script setup lang="ts">
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  value: {
    type: Number,
    required: true,
  },
  color: {
    type: String,
    default: 'bg-blue-500',
  },
  description: {
    type: String,
    required: true,
  },
})
</script>
<template>
  <article class="flex flex-col min-w-3xs w-full max-w-sm p-4 bg-white rounded-lg shadow-md">
    <span class="text-lg font-semibold text-gray-700">{{ props.title }}</span>
    <span :class="['text-3xl font-bold', props.color]">{{ props.value }}</span>
    <p class="text-xs mt-2 text-gray-700">{{ props.description }}</p>
  </article>
</template>
