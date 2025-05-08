<script setup lang="ts">
import { RouterView } from 'vue-router'
import MainHeader from '@/components/MainHeader.vue'
</script>
<template>
  <div class="flex min-h-screen flex-col font-mono">
    <header>
      <MainHeader />
    </header>
    <main class="h-full flex-grow bg-gradient-to-b from-sky-700 from-10% to-black to-35% px-6 pb-4">
      <RouterView />
    </main>
  </div>
</template>
