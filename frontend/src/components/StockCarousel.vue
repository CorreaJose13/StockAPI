<script setup lang="ts">
import type { StockWithScore } from '@/types/types'
import StockCard from '@/components/StockCard.vue'
import ViewAllCTA from '@/components/ViewAllCTA.vue'
import { computed, ref } from 'vue'

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
    type: Array as () => StockWithScore[],
    required: true,
  },
})

const responsiveOptions = ref([
  {
    breakpoint: '1280px',
    numVisible: 3,
    numScroll: 1,
  },
  {
    breakpoint: '1024px',
    numVisible: 3,
    numScroll: 1,
  },
  {
    breakpoint: '768px',
    numVisible: 2,
    numScroll: 1,
  },
  {
    breakpoint: '640px',
    numVisible: 1,
    numScroll: 1,
  },
])

const stocksSlice = computed(() => {
  return props.stocks.slice(0, 10)
})
</script>
<template>
  <section class="flex flex-col gap-4">
    <ViewAllCTA :title="props.title" :link="props.link" />
    <div class="card">
      <Carousel
        :value="stocksSlice"
        :numVisible="3"
        :numScroll="1"
        :responsive-options="responsiveOptions"
      >
        <template #item="slotProps">
          <div class="h-full p-4">
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
          </div>
        </template>
      </Carousel>
    </div>
  </section>
</template>
<style scoped>
:deep(.p-carousel-indicator-button) {
  background: black;
  border: 1px solid white;
}
:deep(.p-carousel-indicator-active .p-carousel-indicator-button) {
  background: white;
}
:deep(.p-carousel-next-button) {
  background: white !important;
}
:deep(.p-carousel-prev-button) {
  background: white !important;
}
</style>
