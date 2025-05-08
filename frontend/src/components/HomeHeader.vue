<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMetricsStore } from '@/stores/metrics'
import graph from '@/assets/graph.jpg'

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
})

const isLoading = ref(false)
const metricsStore = useMetricsStore()

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
      color: 'text-sky-700',
      description: metricsPercentage.value.remain,
    },
  ]
})

onMounted(async () => {
  isLoading.value = true
  await metricsStore.fetchIfNeeded()
  isLoading.value = false
})
</script>
<template>
  <div class="flex h-48 w-full items-center justify-center">
    <div class="flex w-full max-w-screen-2xl flex-row items-center justify-between">
      <span>
        <h1 class="text-5xl font-bold text-white">{{ props.title }}</h1>
        <h2 class="text-xl text-white">{{ props.description }}</h2>
      </span>

      <section v-if="!isLoading" class="animate-fade-in-y flex flex-row gap-4">
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
    </div>
  </div>
</template>
