<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { getStocksMetrics, getStocksList, getStocksAnalysis } from '@/composables/stocks'
import InfoCard from '@/components/InfoCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import ViewHeader from '@/components/ViewHeader.vue'
import type { Stock, StockWithScore } from '@/types/types'

const stocks = ref<Stock[]>([])
const stocksScore = ref<StockWithScore[]>([])
const totalStocks = ref(0)
const upgrade = ref(0)
const downgrade = ref(0)
const remain = ref(0)

const dashboardTexts = computed(() => {
  return {
    title: 'Stock Ratings',
    description: 'Welcome to StockWise! Where you can find stock data at your fingertips',
  }
})

const sectionsTitle = computed(() => {
  return {
    insights: 'Best stocks to buy now',
    stocks: 'Latest analyst ratings',
  }
})

const metricsCards = computed(() => {
  return [
    {
      title: 'Total ratings',
      value: totalStocks.value,
      color: 'text-black',
      description: '',
    },
    {
      title: 'Upgrades',
      value: upgrade.value,
      color: 'text-green-600',
      description: metricsPercentage.value.upgrade,
    },
    {
      title: 'Downgrades',
      value: downgrade.value,
      color: 'text-red-600',
      description: metricsPercentage.value.downgrade,
    },
    {
      title: 'Remains',
      value: remain.value,
      color: 'text-slate-500',
      description: metricsPercentage.value.remain,
    },
  ]
})

const dashboardSections = computed(() => {
  return [
    {
      icon: 'pi pi-star',
      title: sectionsTitle.value.insights,
      link: '/stock-insights',
      stocks: stocksScore.value,
    },
    {
      icon: 'pi pi-eye',
      title: sectionsTitle.value.stocks,
      link: '/all-stocks',
      stocks: stocks.value,
    },
  ]
})

const loadStocks = async (page: number, limit: number) => {
  try {
    const result = await getStocksList(page, limit)
    stocks.value = result.stocks
    totalStocks.value = result.length
  } catch (error) {
    console.error('Error fetching stocks:', error)
  } finally {
  }
}

const loadMetrics = async () => {
  try {
    const result = await getStocksMetrics()
    totalStocks.value = result.total
    upgrade.value = result.upgrade
    downgrade.value = result.downgrade
    remain.value = result.remain
  } catch (error) {
    console.error('Error fetching stocks:', error)
  }
}

const loadAnalysis = async () => {
  try {
    const result = await getStocksAnalysis()
    stocksScore.value = result.stocks
  } catch (error) {
    console.error('Error fetching stocks:', error)
  }
}

const metricsPercentage = computed(() => {
  const calculatePercentage = (value: number) => {
    return totalStocks.value > 0
      ? `${((value / totalStocks.value) * 100).toFixed(2)}% of total stocks`
      : '0% of total stocks'
  }

  return {
    upgrade: calculatePercentage(upgrade.value),
    downgrade: calculatePercentage(downgrade.value),
    remain: calculatePercentage(remain.value),
  }
})

onMounted(async () => {
  const isLoading = ref(true)
  try {
    await Promise.all([loadStocks(1, 10), loadMetrics(), loadAnalysis()])
  } finally {
    isLoading.value = false
  }
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
      />
    </section>
    <DashboardSection
      v-for="(section, index) in dashboardSections"
      :key="index"
      :icon="section.icon"
      :title="section.title"
      :link="section.link"
      :stocks="section.stocks"
    />
  </div>
</template>
