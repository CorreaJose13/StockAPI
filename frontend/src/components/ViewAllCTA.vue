<script setup lang="ts">
import { RouterLink } from 'vue-router'
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  link: {
    type: String,
    required: true,
  },
})
</script>
<template>
  <div class="flex items-center justify-between">
    <h2 class="text-start text-xl font-bold text-white">{{ props.title }}</h2>
    <RouterLink
      :to="props.link"
      class="rounded-lg p-2 text-sm font-semibold text-white hover:bg-white hover:text-black"
    >
      <span class="flex flex-row items-center justify-center gap-2">
        <span>View all</span>
        <i class="pi pi-arrow-right"></i>
      </span>
    </RouterLink>
  </div>
</template>
