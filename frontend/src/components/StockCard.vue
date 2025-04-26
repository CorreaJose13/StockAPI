<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
const visible = ref(false)
const targetDiff = computed(() => {
  const targetDiff = props.targetTo - props.targetFrom
  return targetDiff.toFixed(2)
})
const targetChangeVal = computed(() => {
  const targetChange = (targetDiff.value / props.targetFrom) * 100
  return targetChange.toFixed(2)
})
const targetChange = computed(() => {
  return targetChangeVal.value > 0 ? true : false
})

interface ChartDataset {
  label: string
  fill: boolean
  borderColor: string
  yAxisID: string
  tension: number
  data: number[]
}

interface ChartData {
  labels: string[]
  datasets: ChartDataset[]
}

interface ScaleOptions {
  type?: string
  display?: boolean
  position?: string
  ticks: {
    color: string
  }
  grid: {
    color: string
    drawOnChartArea?: boolean
  }
}

interface ChartOptions {
  stacked: boolean
  maintainAspectRatio: boolean
  aspectRatio: number
  plugins: {
    legend: {
      labels: {
        color: string
      }
    }
  }
  scales: {
    x: ScaleOptions
    y: ScaleOptions
    y1: ScaleOptions
  }
}

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
      {
        label: 'Dataset 2',
        fill: false,
        borderColor: documentStyle.getPropertyValue('--p-gray-500'),
        yAxisID: 'y1',
        tension: 0.4,
        data: [28, 48, 40, 19, 86, 27, 90],
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
</script>
<template>
  <div class="bg-white rounded-lg shadow-md p-4 min-w-3xs w-full max-w-sm" @click="visible = true">
    <div class="flex gap-2 items-center mb-1">
      <span class="font-bold text-lg text-black"> {{ props.ticker }} </span>
      <span class="bg-blue-500 text-white text-xs px-2 py-0.5 rounded-full">
        {{ props.action }}
      </span>
    </div>
    <div class="text-gray-600 text-sm mb-3">{{ props.company }}</div>

    <div class="mb-3">
      <div class="flex items-center">
        <span class="text-black font-semibold">Target Price:</span>
        <span class="text-gray-500 line-through ml-1">${{ props.targetFrom }}</span>
        <span class="ml-1 text-black">→</span>
        <span class="ml-1 font-bold text-black">${{ props.targetTo }}</span>
      </div>
    </div>

    <div class="mb-3">
      <div class="flex flex-col">
        <span class="text-black font-bold">Rating:</span>
        <div>
          <span class="text-gray-500 line-through"> {{ props.ratingFrom }} </span>
          <span class="ml-1 text-black">→</span>
          <span class="ml-1 font-bold text-black"> {{ props.ratingTo }} </span>
        </div>
      </div>
    </div>

    <div class="text-xs text-gray-500">{{ props.brokerage }} • {{ props.time }}</div>
  </div>

  <Dialog v-model:visible="visible" modal header="Edit Profile" class="min-w-xl">
    <template #header>
      <div class="flex flex-col">
        <div class="flex gap-2 items-center mb-1">
          <span class="font-bold text-2xl text-black"> {{ props.ticker }} </span>
          <span class="bg-blue-500 text-white text-xs px-2 py-0.5 rounded-full">
            {{ props.action }}
          </span>
        </div>
        <div class="text-gray-600 text-sm mb-3">{{ props.company }}</div>
      </div>
    </template>

    <div class="flex flex-grow gap-4 mt-2 mb-4">
      <div class="bg-white rounded-lg shadow-md p-4 flex-1 outline">
        <div class="flex flex-col">
          <span class="text-black font-semibold">Target Price:</span>
          <div class="flex flex-grow items-center gap-4">
            <span class="flex-1 w-full font-bold text-black text-2xl">${{ props.targetTo }}</span>
            <div class="flex flex-col flex-1 gap-2">
              <span class="text-gray-500 line-through">${{ props.targetFrom }}</span>
              <div class="flex gap-2">
                <Button
                  :label="targetDiff > 0 ? `+${targetDiff}` : `-${targetDiff}`"
                  variant="text"
                  size="small"
                  raised
                />
                <div v-if="targetChange">
                  <Button
                    icon="pi pi-arrow-up"
                    :label="`${targetChangeVal}%`"
                    variant="text"
                    size="small"
                    raised
                  />
                </div>
                <div v-else>
                  <Button
                    icon="pi pi-arrow-down"
                    :label="`${targetChangeVal}%`"
                    variant="text"
                    size="small"
                    raised
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-md p-4 flex-1 outline">
        <div class="flex flex-col">
          <span class="text-black font-semibold">Rating change:</span>
          <div class="my-2 flex flex-col">
            <span class="font-bold text-black text-2xl">{{ props.ratingTo }}</span>
            <span class="text-gray-500 line-through">{{ props.ratingFrom }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-md p-4 flex-1 outline">
      <div class="flex flex-col">
        <span class="text-black font-semibold">Analysis by:</span>
        <div class="my-2 flex flex-row justify-between">
          <span class="font-bold text-black text-2xl">{{ props.brokerage }}</span>
          <Button icon="pi pi-calendar" :label="`${props.time}`" variant="text" size="small" />
        </div>
      </div>
    </div>
    <div class="card">
      <Chart type="line" :data="chartData" :options="chartOptions" class="h-[30rem]" />
    </div>
  </Dialog>
</template>
