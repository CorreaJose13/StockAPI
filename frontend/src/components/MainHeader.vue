<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { computed } from 'vue'
import logo from '@/assets/logo.svg'

const navItems = computed(() => {
  return [
    { name: 'Stocks', path: '/all-stocks' },
    { name: 'Insights', path: '/stock-insights' },
  ]
})
</script>
<template>
  <div class="w-full bg-stone-900 px-6 py-4">
    <div class="mx-auto flex max-w-screen-2xl flex-row items-center justify-between">
      <div class="flex flex-row items-center gap-4">
        <img :src="logo" alt="Logo" class="h-12 w-auto" />
        <RouterLink class="text-2xl font-bold" to="/">
          <p class="text-white">StockWise</p>
        </RouterLink>
      </div>
      <nav class="flex flex-row gap-4">
        <RouterLink
          v-for="nav in navItems"
          :key="nav.name"
          :to="nav.path"
          class="p-2 text-slate-200 hover:bg-gray-200 hover:text-black focus:bg-gray-200 focus:text-black"
          >{{ nav.name }}</RouterLink
        >
      </nav>
    </div>
  </div>
</template>
