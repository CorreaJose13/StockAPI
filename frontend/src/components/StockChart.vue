<script setup lang="ts">
import type { ChartData } from '@/types/chart'
import { onMounted, ref, nextTick, computed } from 'vue'
import { fetchChartData } from '@/services/chart'
import Chart from 'primevue/chart'

const props = defineProps({
  ticker: {
    type: String,
    required: true,
  },
})

const data = ref<ChartData[]>()
const isLoading = ref(false)

const getRawData = async () => {
  const rawData = await fetchChartData(props.ticker)
  data.value = rawData
}

const hasChartData = computed(() => {
  return data.value && data.value.length > 0
})

const extractDates = (): string[] => {
  if (!data.value) return []
  return data.value.map((item) => item.date)
}

const extractClose = (): number[] => {
  if (!data.value) return []
  return data.value.map((item) => item.close)
}

const extractVolume = (): number[] => {
  if (!data.value) return []
  return data.value.map((item) => item.volume)
}

const setChartData = () => {
  const documentStyle = getComputedStyle(document.documentElement)

  return {
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
        backgroundColor: documentStyle.getPropertyValue('--p-sky-900'),
        tension: 0.4,
        data: extractVolume(),
      },
    ],
  }
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
          drawOnChartArea: false,
          color: surfaceBorder,
        },
      },
    },
  }
}

const chartData = ref()
const chartOptions = ref()

onMounted(async () => {
  isLoading.value = true
  await getRawData()
  await nextTick()
  isLoading.value = false
  chartData.value = setChartData()
  chartOptions.value = setChartOptions()
})
</script>
<template>
  <div v-if="isLoading" class="flex items-center justify-center">
    <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
  </div>
  <div v-else class="card">
    <Chart type="line" :data="chartData" :options="chartOptions" class="h-[20rem]" />
  </div>
</template>
