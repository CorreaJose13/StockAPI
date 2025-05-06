<script setup lang="ts">
import { onMounted, computed } from 'vue'
import InfoCard from '@/components/InfoCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import StockCarousel from '@/components/StockCarousel.vue'
import ViewHeader from '@/components/ViewHeader.vue'
import { useMetricsStore } from '@/stores/metrics'
import { useAnalysisStore } from '@/stores/analysis'
import { useStocksStore } from '@/stores/stocks'

const metricsStore = useMetricsStore()
const analysisStore = useAnalysisStore()
const stocksStore = useStocksStore()

const dashboardTexts = computed(() => {
  return {
    title: 'Stock Ratings',
    description: 'Welcome to StockWise! Where you can find stock data at your fingertips',
  }
})

const metricsCards = computed(() => {
  return [
    {
      title: 'Total ratings',
      value: metricsStore.total,
      color: 'text-black',
      description: '',
      single: true,
    },
    {
      title: 'Upgrades',
      value: metricsStore.upgrade,
      color: 'text-green-600',
      description: metricsPercentage.value.upgrade,
    },
    {
      title: 'Downgrades',
      value: metricsStore.downgrade,
      color: 'text-red-600',
      description: metricsPercentage.value.downgrade,
    },
    {
      title: 'Remains',
      value: metricsStore.remain,
      color: 'text-slate-500',
      description: metricsPercentage.value.remain,
    },
  ]
})

const metricsPercentage = computed(() => {
  const calculatePercentage = (value: number) => {
    return metricsStore.total > 0
      ? `${((value / metricsStore.total) * 100).toFixed(2)}% of total stocks`
      : '0% of total stocks'
  }

  return {
    upgrade: calculatePercentage(metricsStore.upgrade),
    downgrade: calculatePercentage(metricsStore.downgrade),
    remain: calculatePercentage(metricsStore.remain),
  }
})

const loadAllData = async () => {
  await Promise.all([
    stocksStore.fetchInitialDataIfNeeded(),
    metricsStore.fetchIfNeeded(),
    analysisStore.fetchIfNeeded(),
  ])
}

onMounted(async () => {
  loadAllData()
})
</script>
<template>
  <div class="flex flex-col max-w-screen-2xl max-h-screen mx-auto p-4">
    <ViewHeader :title="dashboardTexts.title" :description="dashboardTexts.description" />
    <section class="flex flex-row max-w-screen-xl mx-auto gap-4 justify-around">
      <InfoCard
        v-for="(item, key) in metricsCards"
        :key="key"
        :title="item.title"
        :value="item.value"
        :color="item.color"
        :description="item.description"
        :single="item.single"
      />
    </section>
    <StockCarousel
      icon="pi pi-star"
      title="Top rated stocks"
      link="/stock-insights"
      :stocks="analysisStore.stocks"
    />
    <DashboardSection
      icon="pi pi-eye"
      title="Latest analyst ratings"
      link="/all-stocks"
      :stocks="stocksStore.stocks"
    />
  </div>
</template>
