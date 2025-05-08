<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Stock, StockWithScore } from '@/types/types'
import { getRatingSeverity, getTargetArrow, getTargetSeverity } from '@/utils/stock'
import { formatDateShort, modalDt } from '@/utils/stock'
import { useAnalysisStore } from '@/stores/analysis'

const analysisStore = useAnalysisStore()

const selectedStock = ref<Stock | null>(null)
const showModal = ref(false)

const limit = ref(10)

const onRowSelect = (event: any) => {
  selectedStock.value = event.data
  showModal.value = true
}

const indexedStocksScore = computed(() => {
  return analysisStore.stocks.map((stock, idx) => ({
    ...stock,
    index: `#${idx + 1}`,
  }))
})

const insightTexts = computed(() => {
  return {
    title: 'Market Insights',
    description: 'Explore our buy suggestions in the stock market.',
  }
})

onMounted(async () => {
  await analysisStore.fetchIfNeeded()
})

const tableColumns = computed(() => [
  {
    field: 'ranking',
    header: 'Rank',
    class: 'font-semibold',
    style: 'width: 5%',
    template: (data: StockWithScore) => data.index,
  },
  {
    field: 'ticker',
    header: 'Ticker',
    class: 'font-bold',
    style: 'width: 5%',
    template: (data: StockWithScore) => data.ticker,
  },
  {
    field: 'company',
    header: 'Company',
    style: 'width: 20%',
    template: (data: StockWithScore) => data.company,
  },
  {
    field: 'brokerage',
    header: 'Analyst',
    style: 'width: 15%',
    template: (data: StockWithScore) => data.brokerage,
  },
  {
    field: 'time',
    header: 'Date',
    style: 'width: 10%',
    sortable: true,
    template: (data: Stock) => formatDateShort(data.time),
  },
])

const rowClass = () => {
  return 'cursor-pointer'
}
</script>
<template>
  <div class="mx-auto max-w-screen-2xl">
    <div class="py-4">
      <ViewHeader :title="insightTexts.title" :description="insightTexts.description" />
    </div>
    <DataTable
      :value="indexedStocksScore"
      :loading="analysisStore.loading"
      paginator
      :rows="limit"
      :rowsPerPageOptions="[10, 20, 50]"
      scrollable
      scroll-height="40rem"
      :rowHover="true"
      :rowClass="rowClass"
      @row-click="onRowSelect"
    >
      <Column
        v-for="col in tableColumns"
        :key="col.field"
        :field="col.field"
        :header="col.header"
        :class="col.class"
        class="text-sm text-black dark:text-white"
        :style="col.style"
      >
        <template #body="{ data }">
          <span>{{ col.template(data) }}</span>
        </template>
      </Column>

      <Column field="target" header="Price" class="text-sm text-black">
        <template #body="{ data }">
          <div class="flex flex-row items-center gap-2">
            <Tag class="capitalize" severity="secondary">$ {{ data.target_from }}</Tag>
            <i
              :class="getTargetArrow(data.target_from, data.target_to)"
              style="font-size: 0.75rem"
            ></i>
            <Tag class="capitalize" :severity="getTargetSeverity(data.target_from, data.target_to)"
              >$ {{ data.target_to }}</Tag
            >
          </div></template
        >
      </Column>

      <Column field="rating" header="Rating" class="text-sm text-black">
        <template #body="{ data }">
          <div class="flex flex-row items-center gap-2">
            <Tag
              class="text-sm capitalize"
              :value="data.rating_from"
              :severity="getRatingSeverity(data.rating_from)"
            />
            <i class="pi pi-arrow-right text-gray-500" style="font-size: 0.75rem"></i>
            <Tag
              class="text-sm capitalize"
              :value="data.rating_to"
              :severity="getRatingSeverity(data.rating_to)"
            />
          </div>
        </template>
        ></Column
      >
    </DataTable>
    <Dialog
      v-model:visible="showModal"
      modal
      header="Stock details"
      :dt="modalDt"
      :dismissableMask="true"
    >
      <StockModal
        v-if="selectedStock"
        :ticker="selectedStock.ticker"
        :action="selectedStock.action"
        :company="selectedStock.company"
        :targetFrom="selectedStock.target_from"
        :targetTo="selectedStock.target_to"
        :ratingFrom="selectedStock.rating_from"
        :ratingTo="selectedStock.rating_to"
        :brokerage="selectedStock.brokerage"
        :time="selectedStock.time"
      />
    </Dialog>
  </div>
</template>
<style scoped>
:deep(.p-paginator) {
  border-radius: 0;
}
</style>
