<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { getStocksList } from '@/composables/stocks'
import type { Stock } from '@/types/types'
import { useDebounceFn } from '@vueuse/core'
import { getRatingSeverity, getTargetSeverity, getTargetArrow, formatDate } from '@/utils/stock'

const stocks = ref<Stock[]>([])
const totalStocks = ref(0)
const loading = ref(false)

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
  loadStocks(currentPage.value, limit.value, field.value, order.value)
}

const onSort = (event: any) => {
  field.value = event.sortField
  order.value = event.sortOrder
  loadStocks(currentPage.value, limit.value, field.value, order.value)
}

const debouncedSearch = useDebounceFn((query: string) => {
  currentPage.value = 1
  loadStocks(currentPage.value, limit.value, field.value, order.value, query)
}, 300)

const modalStyle = ref({
  root: {
    background: '{slate.200}',
  },
})

watch(searchQuery, (newValue) => {
  debouncedSearch(newValue)
})

const stocksTableTexts = computed(() => {
  return {
    title: 'Stock Ratings Overview',
    description:
      'Quickly view stock tickers, companies, brokerage actions, ratings, and target prices.',
    placeholder: 'Search by ticker, company, or brokerage',
  }
})

const tableColumns = computed(() => [
  {
    field: 'ticker',
    header: 'Ticker',
    class: 'text-sm font-bold text-black',
    style: 'width: 5%',
    sortable: true,
    template: (data: Stock) => data.ticker,
  },
  {
    field: 'company',
    header: 'Company',
    class: 'text-sm text-black',
    style: 'width: 20%',
    sortable: true,
    template: (data: Stock) => data.company,
  },
  {
    field: 'brokerage',
    header: 'Analyst',
    class: 'text-sm text-black',
    style: 'width: 15%',
    sortable: true,
    template: (data: Stock) => data.brokerage,
  },
  {
    field: 'time',
    header: 'Date',
    class: 'text-sm text-black',
    style: 'width: 10%',
    sortable: true,
    template: (data: Stock) => formatDate(data.time),
  },
  {
    field: 'action',
    header: 'Action',
    class: 'text-sm text-black capitalize',
    style: '',
    sortable: false,
    template: (data: Stock) => data.action,
  },
])

const rowClass = () => {
  return 'cursor-pointer'
}

const loadStocks = async (
  page: number,
  limit: number,
  field?: string,
  order?: number,
  query?: string,
) => {
  const orderString = order === -1 ? 'desc' : 'asc'
  loading.value = true
  try {
    const result = await getStocksList(page, limit, field, orderString, query)
    stocks.value = result.stocks
    totalStocks.value = result.length
  } catch (error) {
    console.error('Error fetching stocks:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStocks(currentPage.value, limit.value)
})
</script>
<template>
  <div class="max-w-screen-2xl max-h-screen mx-auto p-4">
    <ViewHeader :title="stocksTableTexts.title" :description="stocksTableTexts.description" />
    <section>
      <div class="flex justify-start mb-4">
        <IconField>
          <InputIcon>
            <i class="pi pi-search" />
          </InputIcon>
          <InputText
            v-model="searchQuery"
            :placeholder="stocksTableTexts.placeholder"
            class="w-xl"
          />
        </IconField>
      </div>
      <DataTable
        :value="stocks"
        :totalRecords="totalStocks"
        :loading="loading"
        :lazy="true"
        paginator
        :first="(currentPage - 1) * limit"
        :rows="limit"
        :rowsPerPageOptions="[10, 20, 50, 100]"
        @page="onPage"
        scrollable
        scrollHeight="40rem"
        :sortField="field"
        :sortOrder="order"
        @sort="onSort"
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
            <span>{{ col.template(data) }}</span>
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
              <Tag
                class="capitalize"
                :severity="getTargetSeverity(data.target_from, data.target_to)"
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
    </section>
  </div>
</template>
