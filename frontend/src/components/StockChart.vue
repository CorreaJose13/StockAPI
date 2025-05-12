<script setup lang="ts">
import type { ChartData } from '@/types/chart'
import { onMounted, ref, computed, nextTick } from 'vue'
import { fetchChartData } from '@/services/chart'
import Chart from 'primevue/chart'

const props = defineProps({
  ticker: {
    type: String,
    required: true,
  },
})

const data = ref<ChartData[]>([])
const isLoading = ref(false)
const chartReady = ref(false)

const chartData = ref()
const chartOptions = ref()

const hasValidData = computed(() => {
  return data.value && data.value.length > 0
})

const getRawData = async () => {
  try {
    isLoading.value = true
    const rawData = await fetchChartData(props.ticker)
    data.value = rawData || []
  } catch (error) {
    data.value = []
  } finally {
    isLoading.value = false
  }
}

const extractDates = (): string[] => {
  if (!hasValidData.value) return []
  return data.value.map((item) => item.date)
}

const extractClose = (): number[] => {
  if (!hasValidData.value) return []
  return data.value.map((item) => item.close)
}

const extractVolume = (): number[] => {
  if (!hasValidData.value) return []
  return data.value.map((item) => item.volume)
}

const prepareChartData = () => {
  if (!hasValidData.value) return

  const documentStyle = getComputedStyle(document.documentElement)

  chartData.value = {
    labels: extractDates(),
    datasets: [
      {
        type: 'line',
        label: 'Close Price',
        fill: false,
        backgroundColor: 'black',
        borderColor: 'black',
        yAxisID: 'close',
        tension: 0.4,
        data: extractClose(),
      },
      {
        type: 'bar',
        label: 'Volume',
        fill: false,
        yAxisID: 'volume',
        backgroundColor: documentStyle.getPropertyValue('--p-sky-900') || '#0c4a6e',
        tension: 0.4,
        data: extractVolume(),
      },
    ],
  }

  chartOptions.value = setChartOptions()
}

const setChartOptions = () => {
  const documentStyle = getComputedStyle(document.documentElement)
  const textColor = documentStyle.getPropertyValue('--p-text-color')
  const textColorSecondary = documentStyle.getPropertyValue('--p-text-muted-color')
  const surfaceBorder = documentStyle.getPropertyValue('--p-content-border-color')

  return {
    stacked: true,
    maintainAspectRatio: false,
    aspectRatio: 0.6,
    plugins: {
      legend: {
        labels: {
          color: textColor,
        },
      },
    },
    scales: {
      x: {
        ticks: {
          color: textColorSecondary,
        },
        grid: {
          color: surfaceBorder,
        },
      },
      close: {
        type: 'linear',
        display: true,
        position: 'right',
        ticks: {
          color: textColorSecondary,
        },
        grid: {
          color: surfaceBorder,
        },
      },
      volume: {
        type: 'linear',
        display: true,
        position: 'left',
        ticks: {
          color: textColorSecondary,
        },
        grid: {
          color: surfaceBorder,
        },
      },
    },
  }
}

onMounted(async () => {
  await getRawData()

  await nextTick()

  if (hasValidData.value) {
    prepareChartData()

    await nextTick()

    chartReady.value = true
  }
})
</script>
<template>
  <div v-if="isLoading" class="flex items-center justify-center">
    <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
  </div>
  <div v-else-if="!hasValidData" class="flex items-center justify-center p-8">
    <span class="text-gray-500">No chart data available for {{ props.ticker }}</span>
  </div>
  <div v-else-if="chartReady && chartData && chartOptions" class="card">
    <Chart type="line" :data="chartData" :options="chartOptions" class="h-[20rem]" />
  </div>
</template>
