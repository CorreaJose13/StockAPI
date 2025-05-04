<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getStocksAnalysis } from '@/composables/stocks'
import type { StockWithScore, Stock } from '@/types/types'
import { getRatingSeverity, getTargetArrow, getTargetSeverity } from '@/utils/stock'

const stocksScore = ref<StockWithScore[]>([])
const loading = ref(false)
const selectedStock = ref<Stock | null>(null)
const showModal = ref(false)

const currentPage = ref(1)
const limit = ref(10)

const onRowSelect = (event: any) => {
  selectedStock.value = event.data
  showModal.value = true
}

const modalStyle = ref({
  root: {
    background: '{slate.200}',
  },
})

const indexedStocksScore = computed(() => {
  return stocksScore.value.map((stock, idx) => ({
    ...stock,
    index: `#${idx + 1}`,
  }))
})

const insightTexts = computed(() => {
  return {
    title: 'Market Insights',
    description: 'Explore market trends and expert analysis in the stock market.',
    insightsTitle: 'Suggested Stocks',
    sectorsTitle: 'Sector Highlights',
  }
})

const loadAnalysis = async () => {
  try {
    const result = await getStocksAnalysis()
    stocksScore.value = result.stocks
  } catch (error) {
    console.error('Error fetching stocks:', error)
  }
}

onMounted(async () => {
  await loadAnalysis()
})

const tableColumns = computed(() => [
  {
    field: 'ranking',
    header: 'Rank',
    class: 'text-sm font-bold text-black',
    style: 'width: 5%',
    sortable: false,
    template: (data: StockWithScore) => data.index,
  },
  {
    field: 'ticker',
    header: 'Ticker',
    class: 'text-sm font-bold text-black',
    style: 'width: 5%',
    sortable: false,
    template: (data: StockWithScore) => data.ticker,
  },
  {
    field: 'company',
    header: 'Company',
    class: 'text-sm text-black',
    style: 'width: 20%',
    sortable: false,
    template: (data: StockWithScore) => data.company,
  },
  {
    field: 'brokerage',
    header: 'Analyst',
    class: 'text-sm text-black',
    style: 'width: 15%',
    sortable: false,
    template: (data: StockWithScore) => data.brokerage,
  },
  {
    field: 'action',
    header: 'Action',
    class: 'text-sm text-black capitalize',
    style: '',
    sortable: false,
    template: (data: StockWithScore) => data.action,
  },
])

const rowClass = () => {
  return 'cursor-pointer'
}
</script>
<template>
  <div class="max-w-screen-2xl max-h-screen mx-auto p-4">
    <ViewHeader :title="insightTexts.title" :description="insightTexts.description" />
    <section class="flex flex-col gap-2 w-full py-4 mt-4">
      <div class="flex flex-row gap-2 items-center">
        <i class="pi pi-briefcase text-white"></i>
        <h2 class="text-xl font-bold text-start text-white">{{ insightTexts.insightsTitle }}</h2>
      </div>
    </section>
    <DataTable
      :value="indexedStocksScore"
      :loading="loading"
      paginator
      :first="(currentPage - 1) * limit"
      :rows="limit"
      :rowsPerPageOptions="[10, 20]"
      scrollable
      scrollHeight="40rem"
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
        :style="col.style"
        :sortable="col.sortable"
      >
        <template #body="{ data }">
          <span class="capitalize">{{ col.template(data) }}</span>
        </template>
      </Column>
      <Column field="rating" header="Rating" class="text-sm text-black" style="width: 25%">
        <template #body="{ data }">
          <div class="flex flex-row gap-2 items-center">
            <Tag
              class="capitalize text-sm"
              :value="data.rating_from"
              :severity="getRatingSeverity(data.rating_from)"
            />
            <i class="pi pi-arrow-right text-gray-500" style="font-size: 0.75rem"></i>
            <Tag
              class="capitalize text-sm"
              :value="data.rating_to"
              :severity="getRatingSeverity(data.rating_to)"
            />
          </div>
        </template>
        ></Column
      >
      <Column field="target" header="Price" class="text-sm text-black" style="width: 15%">
        <template #body="{ data }">
          <div class="flex flex-row gap-2 items-center">
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
    </DataTable>
    <Dialog
      v-model:visible="showModal"
      modal
      header="Stock details"
      class="min-w-xl"
      :dt="modalStyle"
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
