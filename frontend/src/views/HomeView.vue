<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import LatestSection from '@/components/LatestSection.vue'
import StockCarousel from '@/components/StockCarousel.vue'
import { useAnalysisStore } from '@/stores/analysis'
import { useLatestStore } from '@/stores/latest'
import HomeHeader from '@/components/HomeHeader.vue'

const isLoading = ref(true)
const analysisStore = useAnalysisStore()
const latestStore = useLatestStore()

const dashboardTexts = computed(() => {
  return {
    title: 'Welcome to StockWise!',
    description: ' Where you can find stock data at your fingertips.',
  }
})

const loadAllData = async () => {
  await analysisStore.fetchIfNeeded()
  await latestStore.fetchInitialDataIfNeeded()
}

onMounted(async () => {
  isLoading.value = true
  await loadAllData()
  isLoading.value = false
})
</script>
<template>
  <div>
    <HomeHeader
      :title="dashboardTexts.title"
      :description="dashboardTexts.description"
    ></HomeHeader>
    <div class="mx-auto flex max-w-screen-2xl flex-col">
      <StockCarousel
        v-if="!isLoading"
        title="⭐ Top rated stocks"
        link="/stock-insights"
        :stocks="analysisStore.stocks"
        class="animate-fade-in-x"
      />
      <LatestSection
        v-if="!isLoading"
        title="⏱️ Latest analyst ratings"
        link="/all-stocks"
        :stocks="latestStore.stocks"
        class="animate-fade-in-x-reverse"
      />
    </div>
  </div>
</template>
