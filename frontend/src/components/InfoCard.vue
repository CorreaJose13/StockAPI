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
  <article class="flex flex-col w-40 p-4 bg-white rounded-lg shadow-md">
    <span class="text-lg font-semibold text-gray-700">{{ props.title }}</span>
    <div v-if="props.single">
      <span :class="['text-5xl font-bold', props.color]">{{ props.value }}</span>
    </div>
    <div v-else>
      <span :class="['text-3xl font-bold', props.color]">{{ props.value }}</span>
      <p class="text-xs mt-2 text-gray-700">{{ props.description }}</p>
    </div>
  </article>
</template>
