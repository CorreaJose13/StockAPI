<script setup lang="ts">
import Carousel from 'primevue/carousel'
import StockCard from '@/components/StockCard.vue'
import { RouterLink } from 'vue-router'
import type { StockWithScore } from '@/types/types'

const props = defineProps({
  icon: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  link: {
    type: String,
    required: true,
  },
  stocks: {
    type: Array as () => StockWithScore[],
    required: true,
  },
})
</script>

<template>
  <section class="flex flex-col gap-4 my-6">
    <div class="flex justify-between items-center">
      <div class="flex flex-row gap-2 items-center">
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
    </div>

    <Carousel :value="props.stocks" :numVisible="3" :numScroll="1" circular>
      <template #item="slotProps">
        <section class="flex flex-col w-full p-4">
          <StockCard
            :ticker="slotProps.data.ticker"
            :action="slotProps.data.action"
            :company="slotProps.data.company"
            :targetFrom="Number(slotProps.data.target_from || 0)"
            :targetTo="Number(slotProps.data.target_to || 0)"
            :ratingFrom="slotProps.data.rating_from"
            :ratingTo="slotProps.data.rating_to"
            :brokerage="slotProps.data.brokerage"
            :time="slotProps.data.time"
          />
        </section>
      </template>
    </Carousel>
  </section>
</template>
<style scoped>
:deep(.p-carousel-indicator-button) {
  background: white;
}
:deep(.p-carousel-indicator-active .p-carousel-indicator-button) {
  background: black;
}
:deep(.p-carousel-next-button) {
  background: white;
}
:deep(.p-carousel-prev-button) {
  background: white !important;
}
</style>
