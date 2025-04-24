<script setup lang="ts">
import { computed } from 'vue'
import InfoCard from '@/components/InfoCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'

const useStockData = computed(() => {
  const stocks = [
    {
      id: 1,
      ticker: 'AAPL',
      targetFrom: 180.5,
      targetTo: 215.75,
      company: 'Apple Inc.',
      action: 'upgraded by',
      brokerage: 'Morgan Stanley',
      ratingFrom: 'Hold',
      ratingTo: 'Buy',
      time: '2025-01-15T09:30:00.000Z',
    },
    {
      id: 2,
      ticker: 'MSFT',
      targetFrom: 320.0,
      targetTo: 375.25,
      company: 'Microsoft Corporation',
      action: 'reiterated by',
      brokerage: 'Goldman Sachs',
      ratingFrom: 'Buy',
      ratingTo: 'Buy',
      time: '2025-01-14T14:45:00.000Z',
    },
    {
      id: 3,
      ticker: 'TSLA',
      targetFrom: 195.0,
      targetTo: 165.5,
      company: 'Tesla, Inc.',
      action: 'downgraded by',
      brokerage: 'JP Morgan',
      ratingFrom: 'Neutral',
      ratingTo: 'Sell',
      time: '2025-01-13T11:15:00.000Z',
    },
  ]
  return stocks
})
</script>
<template>
  <div class="min-h-screen bg-gray-100">
    <div class="flex flex-col w-full max-w-screen-xl mx-auto">
      <section class="flex flex-col gap-1 w-full py-4">
        <h1 class="text-3xl font-bold text-start text-black">Dashboard</h1>
        <p class="text-start text-gray-600">
          Welcome to the dashboard! Here you can find all the information you need.
        </p>
      </section>
      <section class="flex flex-col w-full md:flex-row justify-between">
        <InfoCard title="Total Users" :value="12" color="text-black" description="Total number " />
        <InfoCard
          title="Total Users"
          :value="12"
          color="text-green-500"
          description="Total number "
        />
        <InfoCard
          title="Total Users"
          :value="12"
          color="text-red-600"
          description="Total number "
        />
      </section>
      <DashboardSection
        title="Recommendations"
        link="/stock-recommendations"
        :stocks="useStockData"
      />
      <DashboardSection title="Stocks" link="/all-stocks" :stocks="useStockData" />
    </div>
  </div>
</template>
