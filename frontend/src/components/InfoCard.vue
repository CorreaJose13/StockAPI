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
  single: {
    type: Boolean,
    default: false,
  },
})
</script>
<template>
  <article class="flex w-40 flex-col rounded-lg bg-white p-4">
    <span class="text-lg font-semibold text-slate-800">{{ props.title }}</span>
    <div v-if="props.single">
      <span :class="['text-5xl font-bold', props.color]">{{ props.value }}</span>
    </div>
    <div v-else>
      <span :class="['text-3xl font-bold', props.color]">{{ props.value }}</span>
      <p class="mt-2 text-xs text-gray-700">{{ props.description }}</p>
    </div>
  </article>
</template>
