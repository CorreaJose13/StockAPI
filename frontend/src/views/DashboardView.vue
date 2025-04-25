<script setup lang="ts">
import { useStockData } from '@/composables/stocks'
import InfoCard from '@/components/InfoCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
const stocks = useStockData()
</script>
<template>
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
      <InfoCard title="Total Users" :value="12" color="text-red-600" description="Total number " />
    </section>
    <DashboardSection title="Recommendations" link="/stock-recommendations" :stocks="stocks" />
    <DashboardSection title="Stocks" link="/all-stocks" :stocks="stocks" />
  </div>
</template>
