<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { computed } from 'vue'
import logo from '@/assets/logo.svg'

const navItems = computed(() => {
  return [
    { name: 'Dashboard', path: '/' },
    { name: 'Stocks', path: '/all-stocks' },
    { name: 'Insights', path: '/stock-insights' },
  ]
})
</script>
<template>
  <div class="flex flex-col min-h-screen">
    <header>
      <div class="w-full text-black bg-white">
        <div
          class="flex flex-col max-w-screen-2xl mx-auto md:items-center md:justify-between md:flex-row"
        >
          <div class="flex flex-row items-center justify-between p-4">
            <img :src="logo" alt="Logo" class="h-8 w-auto mr-2" />
            <RouterLink class="text-xl font-semibold flex" to="/">
              <p class="text-blue-900">Stock</p>
              <p class="text-blue-500">Wise</p>
            </RouterLink>
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

    <main class="flex-grow bg-gradient-to-r from-blue-900 via-blue-700 to-blue-800 pb-20">
      <RouterView />
    </main>
  </div>
</template>
