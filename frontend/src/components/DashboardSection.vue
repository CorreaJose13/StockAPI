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
  icon: {
    type: String,
  },
})
const limitedStocks = computed(() => {
  return props.stocks.slice(0, 5)
})
</script>
<template>
  <div>
    <section class="flex flex-row justify-between w-full py-4">
      <div class="flex flex-row gap-2 justify-center items-center">
        <i :class="['text-white', props.icon]"></i>
        <h2 class="text-xl font-bold text-start text-white">{{ props.title }}</h2>
      </div>
      <RouterLink
        :to="props.link"
        class="p-2 text-sm rounded-lg text-white font-semibold hover:text-black hover:bg-white"
      >
        <span class="flex flex-row gap-2 items-center justify-center">
          <span>See all</span>
          <i class="pi pi-arrow-right"></i>
        </span>
      </RouterLink>
    </section>
    <section class="flex flex-col w-full md:flex-row justify-between gap-4">
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
