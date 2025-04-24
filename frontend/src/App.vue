<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { computed } from 'vue'
import logo from '@/assets/logo.svg'

const navItems = computed(() => {
  return [
    { name: 'Dashboard', path: '/' },
    { name: 'Stocks', path: '/all-stocks' },
    { name: 'Recommendations', path: '/stock-recommendations' },
  ]
})
</script>
<template>
  <header>
    <div class="w-full text-gray-700 bg-white">
      <div
        class="flex flex-col max-w-screen-xl px-4 mx-auto md:items-center md:justify-between md:flex-row md:px-6 lg:px-8"
      >
        <div class="flex flex-row items-center justify-between p-4">
          <img :src="logo" alt="Logo" class="h-8 w-auto mr-2" />
          <RouterLink class="text-lg font-semibold text-gray-900 uppercase" to="/">
            Stock API</RouterLink
          >
        </div>
        <nav class="flex-col flex-grow hidden pb-4 md:pb-0 md:flex md:justify-end md:flex-row">
          <RouterLink
            v-for="nav in navItems"
            :key="nav.name"
            :to="nav.path"
            class="px-2 py-2 mt-2 text-sm font-semibold bg-transparent rounded-lg md:mt-0 md:ml-4 hover:text-gray-900 focus:text-gray-900 hover:bg-gray-200 focus:bg-gray-200 focus:outline-none focus:shadow-outline"
            >{{ nav.name }}</RouterLink
          >
        </nav>
      </div>
    </div>
  </header>

  <RouterView />
</template>
