<script setup lang="ts">
import { ref, computed } from 'vue'
import { useStockData, getStocksByCategory, getStockCategories } from '@/composables/stocks'
import Carousel from 'primevue/carousel'
import StockCard from '@/components/StockCard.vue'

const stocks = useStockData()
const categories = getStockCategories()

const activeTab = ref(categories[0])

const tabs = ref(
  categories.map((category) => ({
    title: category,
    value: category,
  })),
)

const filteredStocks = (category: string) => {
  return getStocksByCategory(category)
}

const insightTexts = computed(() => {
  return {
    title: 'Market Insights',
    description: 'Explore market trends and expert analysis in the stock market.',
    insightsTitle: 'Trending Stocks',
    sectorsTitle: 'Sector Highlights',
  }
})
</script>
<template>
  <div class="max-w-screen-xl max-h-screen mx-auto p-4">
    <section class="flex flex-col gap-1 w-full py-4">
      <h1 class="text-3xl font-bold text-start text-black">{{ insightTexts.title }}</h1>
      <p class="text-start text-gray-600">
        {{ insightTexts.description }}
      </p>
    </section>
    <section class="flex flex-col gap-2 w-full py-4 mt-4">
      <h2 class="text-xl font-bold text-start text-black">{{ insightTexts.insightsTitle }}</h2>
      <Carousel :value="stocks" :numVisible="3" :numScroll="3" circular :autoplayInterval="5000">
        <template #item="slotProps">
          <section class="flex flex-col w-full md:flex-row justify-between p-4">
            <StockCard
              :ticker="slotProps.data.ticker"
              :action="slotProps.data.action"
              :company="slotProps.data.company"
              :targetFrom="slotProps.data.targetFrom"
              :targetTo="slotProps.data.targetTo"
              :ratingFrom="slotProps.data.ratingFrom"
              :ratingTo="slotProps.data.ratingTo"
              :brokerage="slotProps.data.brokerage"
              :time="slotProps.data.time"
            />
          </section>
        </template>
      </Carousel>
    </section>
    <section class="flex flex-col gap-2 w-full py-4 mt-4">
      <h2 class="text-xl font-bold text-start text-black">{{ insightTexts.sectorsTitle }}</h2>
      <div class="flex flex-row justify-center gap-4">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          class="text-gray-500 rounded-lg px-4 py-2 text-md font-medium"
          :class="
            activeTab === tab.value ? 'bg-blue-500 text-white' : 'text-gray-500 hover:text-gray-700'
          "
          @click="activeTab = tab.value"
        >
          {{ tab.title }}
        </button>
      </div>
      <div class="flex flex-grow justify-between mt-2">
        <StockCard
          v-for="stock in filteredStocks(activeTab)"
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
      </div>
    </section>
  </div>
</template>
<style scoped>
:deep(.p-carousel-indicator-button) {
  background: gray;
}
:deep(.p-carousel-indicator-active .p-carousel-indicator-button) {
  background: black;
}
</style>
