<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import type { Stock } from '@/types/types'
import StockCard from './StockCard.vue'
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  link: {
    type: String,
    required: true,
  },
  stocks: {
    type: Array as () => Stock[],
    required: true,
  },
})
const limitedStocks = computed(() => {
  return props.stocks.slice(0, 3)
})
</script>
<template>
  <div>
    <section class="flex flex-row justify-between w-full py-4">
      <h2 class="text-xl font-bold text-start text-black">{{ props.title }}</h2>
      <RouterLink
        :to="props.link"
        class="p-2 text-sm rounded-lg text-gray-500 font-semibold hover:text-gray-900 focus:text-gray-900 hover:bg-gray-300 focus:bg-gray-300 focus:outline-none focus:shadow-outline"
        >View all &gt;</RouterLink
      >
    </section>
    <section class="flex flex-col w-full md:flex-row justify-between">
      <StockCard
        v-for="stock in limitedStocks"
        :key="stock.id"
        :ticker="stock.ticker"
        :action="stock.action"
        :company="stock.company"
        :targetFrom="stock.targetFrom"
        :targetTo="stock.targetTo"
        :ratingFrom="stock.ratingFrom"
        :ratingTo="stock.ratingTo"
        :brokerage="stock.brokerage"
        :time="stock.time"
      />
    </section>
  </div>
</template>
