<script setup lang="ts">
import { RouterView } from 'vue-router'
import MainHeader from '@/components/MainHeader.vue'
</script>
<template>
  <div class="flex min-h-screen flex-col font-mono">
    <header>
      <MainHeader />
    </header>
    <main class="h-full flex-grow bg-black px-6 pb-4">
      <RouterView />
    </main>
  </div>
</template>
