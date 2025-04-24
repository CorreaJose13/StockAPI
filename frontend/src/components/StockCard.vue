<script setup lang="ts">
const props = defineProps({
  ticker: {
    type: String,
    required: true,
  },
  action: {
    type: String,
    required: true,
  },
  company: {
    type: String,
    required: true,
  },
  targetFrom: {
    type: Number,
    required: true,
  },
  targetTo: {
    type: Number,
    required: true,
  },
  ratingFrom: {
    type: String,
    required: true,
  },
  ratingTo: {
    type: String,
    required: true,
  },
  brokerage: {
    type: String,
    required: true,
  },
  time: {
    type: String,
    required: true,
  },
})
</script>
<template>
  <div class="bg-white rounded-lg shadow-md p-4 min-w-3xs w-full max-w-sm">
    <div class="flex gap-2 items-center mb-1">
      <span class="font-bold text-lg"> {{ props.ticker }} </span>
      <span class="bg-blue-500 text-white text-xs px-2 py-0.5 rounded-full">
        {{ props.action }}
      </span>
    </div>
    <div class="text-gray-600 text-sm mb-3">{{ props.company }}</div>

    <div class="mb-3">
      <div class="flex items-center">
        <span class="text-black font-semibold">Target Price:</span>
        <span class="text-gray-500 line-through ml-1">${{ props.targetFrom }}</span>
        <span class="ml-1">→</span>
        <span class="ml-1 font-bold">${{ props.targetTo }}</span>
      </div>
    </div>

    <div class="mb-3">
      <div class="flex flex-col">
        <span class="text-black font-bold">Rating:</span>
        <div>
          <span class="text-gray-500 line-through"> {{ props.ratingFrom }} </span>
          <span class="ml-1">→</span>
          <span class="ml-1 font-bold"> {{ props.ratingTo }} </span>
        </div>
      </div>
    </div>

    <div class="text-xs text-gray-500">{{ props.brokerage }} • {{ props.time }}</div>
  </div>
</template>
