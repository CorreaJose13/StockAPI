<script setup lang="ts">
import { getRatingSeverity, getTargetSeverity, tableDt, rowClass } from '@/utils/stock'
import { getTargetArrow, formatDateShort, modalDt } from '@/utils/stock'
import { onMounted, ref, computed, watch } from 'vue'
import type { Stock } from '@/types/types'
import { useDebounceFn } from '@vueuse/core'
import { useStocksStore } from '@/stores/pagination'

const stocksStore = useStocksStore()

const currentPage = ref(1)
const limit = ref(10)
const field = ref('')
const order = ref(0)
const searchQuery = ref('')

const selectedStock = ref<Stock | null>(null)
const showModal = ref(false)

const onRowSelect = (event: any) => {
  selectedStock.value = event.data
  showModal.value = true
}

const onPage = (event: any) => {
  currentPage.value = event.page + 1
  limit.value = event.rows
  stocksStore.fetchData(currentPage.value, limit.value, field.value, formatOrder(order.value))
}

const onSort = (event: any) => {
  field.value = event.sortField
  order.value = event.sortOrder
  stocksStore.fetchData(currentPage.value, limit.value, field.value, formatOrder(order.value))
}

const formatOrder = (order: number) => {
  return order === -1 ? 'desc' : order === 1 ? 'asc' : ''
}

const debouncedSearch = useDebounceFn((query: string) => {
  currentPage.value = 1
  stocksStore.fetchData(
    currentPage.value,
    limit.value,
    field.value,
    formatOrder(order.value),
    query,
  )
}, 500)

watch(searchQuery, (newValue) => {
  debouncedSearch(newValue)
})

const stocksTableTexts = computed(() => {
  return {
    title: 'Ratings Overview',
    description:
      'Quickly view tickers, prices, and analyst ratings. Click on a stock to see more details.',
    placeholder: 'Search by ticker, company, or analyst',
  }
})

const tableColumns = computed(() => [
  {
    field: 'ticker',
    header: 'Ticker',
    class: 'font-bold',
    style: 'width: 5%',
    sortable: true,
    template: (data: Stock) => data.ticker,
  },
  {
    field: 'company',
    header: 'Company',
    style: 'width: 20%',
    sortable: true,
    template: (data: Stock) => data.company,
  },
  {
    field: 'brokerage',
    header: 'Analyst',
    style: 'width: 15%',
    sortable: true,
    template: (data: Stock) => data.brokerage,
  },
  {
    field: 'time',
    header: 'Date',
    style: 'width: 10%',
    sortable: true,
    template: (data: Stock) => formatDateShort(data.time),
  },
  {
    field: 'action',
    header: 'Action',
    style: 'width: 15%',
    sortable: false,
    template: (data: Stock) => data.action,
  },
])

onMounted(async () => {
  await stocksStore.fetchInitialDataIfNeeded()
})
</script>
<template>
  <div class="mx-auto max-w-screen-2xl">
    <div class="flex flex-row items-end justify-between py-4">
      <ViewHeader :title="stocksTableTexts.title" :description="stocksTableTexts.description" />
      <IconField class="py-4">
        <InputIcon>
          <i class="pi pi-search" />
        </InputIcon>
        <InputText v-model="searchQuery" :placeholder="stocksTableTexts.placeholder" class="w-xl" />
      </IconField>
    </div>
    <section>
      <DataTable
        :value="stocksStore.stocks"
        :totalRecords="stocksStore.total"
        :loading="stocksStore.loading"
        :lazy="true"
        paginator
        :first="(currentPage - 1) * limit"
        :rows="limit"
        :rowsPerPageOptions="[10, 20, 50, 100]"
        @page="onPage"
        scrollable
        scroll-height="40rem"
        :sortField="field"
        :sortOrder="order"
        @sort="onSort"
        :rowHover="true"
        :rowClass="rowClass"
        @row-click="onRowSelect"
        :dt="tableDt"
      >
        <Column
          v-for="col in tableColumns"
          :key="col.field"
          :field="col.field"
          :header="col.header"
          :class="col.class"
          class="text-sm text-black dark:text-white"
          :style="col.style"
          :sortable="col.sortable"
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
              <Tag
                class="capitalize"
                :severity="getTargetSeverity(data.target_from, data.target_to)"
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
    </section>
  </div>
</template>
<style scoped>
:deep(.p-paginator) {
  border-radius: 0;
}
</style>
