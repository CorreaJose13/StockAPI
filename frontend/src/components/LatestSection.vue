<script setup lang="ts">
import type { Stock } from '@/types/types'
import ViewAllCTA from './ViewAllCTA.vue'
import StockCard from './StockCard.vue'
import { computed } from 'vue'
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
  icon: {
    type: String,
  },
})
const limitedStocks = computed(() => {
  return props.stocks.slice(0, 5)
})
</script>
<template>
  <div class="flex flex-col gap-4">
    <ViewAllCTA :title="title" :link="link" />
    <section class="grid w-full grid-cols-5 gap-4">
      <StockCard
        v-for="stock in limitedStocks"
        :key="stock.ticker"
        :ticker="stock.ticker"
        :action="stock.action"
        :company="stock.company"
        :targetFrom="stock.target_from"
        :targetTo="stock.target_to"
        :ratingFrom="stock.rating_from"
        :ratingTo="stock.rating_to"
        :brokerage="stock.brokerage"
        :time="stock.time"
      />
    </section>
  </div>
</template>
