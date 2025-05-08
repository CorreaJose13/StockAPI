<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import type { ChartData, ChartOptions } from '@/types/chart'
import { priceAbsDiff, pricePercDiff, formatDate } from '@/utils/stock'

const props = defineProps({
  ticker: { type: String, required: true },
  action: { type: String, required: true },
  company: { type: String, required: true },
  targetFrom: { type: Number, required: true },
  targetTo: { type: Number, required: true },
  ratingFrom: { type: String, required: true },
  ratingTo: { type: String, required: true },
  brokerage: { type: String, required: true },
  time: { type: String, required: true },
})

const chartData = ref<ChartData | null>(null)
const chartOptions = ref<ChartOptions | null>(null)

onMounted(() => {
  chartData.value = setChartData()
  chartOptions.value = setChartOptions()
})

const setChartData = (): ChartData => {
  const documentStyle = getComputedStyle(document.documentElement)

  return {
    labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
    datasets: [
      {
        label: 'Dataset 1',
        fill: false,
        borderColor: documentStyle.getPropertyValue('--p-cyan-500'),
        yAxisID: 'y',
        tension: 0.4,
        data: [65, 59, 80, 81, 56, 55, 10],
      },
    ],
  }
}

const setChartOptions = (): ChartOptions => {
  const documentStyle = getComputedStyle(document.documentElement)
  const textColor = documentStyle.getPropertyValue('--p-text-color')
  const textColorSecondary = documentStyle.getPropertyValue('--p-text-muted-color')
  const surfaceBorder = documentStyle.getPropertyValue('--p-content-border-color')

  return {
    stacked: false,
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
      y: {
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
      y1: {
        type: 'linear',
        display: true,
        position: 'right',
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

const modalTexts = computed(() => {
  return {
    price: 'Price:',
    rating: 'Rating:',
    analysisBy: 'Analysis by:',
  }
})

const getAbsSeverity = (from: number, to: number) => {
  const diff = Number(priceAbsDiff(from, to))
  if (diff > 0) {
    return 'success'
  } else if (diff < 0) {
    return 'danger'
  } else {
    return 'secondary'
  }
}

const getPercSeverity = (from: number, to: number) => {
  const diff = Number(pricePercDiff(from, to))
  if (diff > 0) {
    return 'success'
  } else if (diff < 0) {
    return 'danger'
  } else {
    return 'secondary'
  }
}

const getPerc = computed(() => {
  const percDiff = pricePercDiff(props.targetFrom, props.targetTo)
  return `${Number(percDiff).toFixed(2)}%`
})
</script>
<template>
  <div class="flex flex-col">
    <div class="mb-1 flex items-center gap-2">
      <span class="text-2xl font-bold text-black"> {{ props.ticker }} </span>
      <span class="rounded-full bg-slate-500 px-2 py-0.5 text-xs text-white">
        {{ props.action }}
      </span>
    </div>
    <div class="mb-3 text-sm text-gray-600">{{ props.company }}</div>
  </div>

  <div class="mt-2 mb-4 flex flex-grow gap-4">
    <div class="flex-1 rounded-lg bg-white p-4 shadow-md outline">
      <div class="flex flex-col">
        <span class="font-semibold text-black"> {{ modalTexts.price }}</span>
        <div class="flex flex-grow items-center gap-4">
          <span class="w-full flex-1 text-2xl font-bold text-black">${{ props.targetTo }}</span>
          <div class="flex flex-1 flex-col gap-2">
            <span class="text-gray-500 line-through">${{ props.targetFrom }}</span>
            <div class="flex gap-2">
              <Tag :value="getPerc" :severity="getPercSeverity(props.targetFrom, props.targetTo)">
              </Tag>
              <Tag
                :value="`$${priceAbsDiff(props.targetFrom, props.targetTo)}`"
                :severity="getAbsSeverity(props.targetFrom, props.targetTo)"
              ></Tag>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 rounded-lg bg-white p-4 shadow-md outline">
      <div class="flex flex-col">
        <span class="font-semibold text-black"> {{ modalTexts.rating }}</span>
        <section class="my-2 flex flex-col">
          <span class="text-2xl font-bold text-black capitalize">{{ props.ratingTo }}</span>
          <span class="text-gray-500 capitalize line-through">{{ props.ratingFrom }}</span>
        </section>
      </div>
    </div>
  </div>

  <div class="flex-1 rounded-lg bg-white p-4 shadow-md outline">
    <div class="flex flex-col">
      <span class="font-semibold text-black">{{ modalTexts.analysisBy }}</span>
      <div class="my-2 flex flex-row justify-between">
        <span class="text-2xl text-black">{{ props.brokerage }}</span>
        <div class="flex items-center gap-2">
          <i class="pi pi-calendar text-black"></i>
          <span class="text-black">{{ formatDate(props.time) }}</span>
        </div>
      </div>
    </div>
  </div>
  <div class="card"></div>
</template>
