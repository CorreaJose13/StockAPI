<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { getStocksList } from '@/composables/stocks'
import type { Stock } from '@/types/types'
import { useDebounceFn } from '@vueuse/core'

const stocks = ref<Stock[]>([])
const totalStocks = ref(0)
const loading = ref(false)

const currentPage = ref(1)
const limit = ref(10)
const field = ref('')
const order = ref(0)
const searchQuery = ref('')

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
  currentPage.value = 1 // Resetear a la primera página al buscar
  loadStocks(currentPage.value, limit.value, field.value, order.value, query)
}, 300)

watch(searchQuery, (newValue) => {
  debouncedSearch(newValue)
})

const getRatingSeverity = (rating: string) => {
  switch (rating) {
    case 'sell':
      return 'danger'
    case 'buy':
      return 'success'
    case 'outperform':
      return 'info'
    case 'underperform':
      return 'warn'
    case 'hold':
      return 'secondary'
    default:
      return 'secondary'
  }
}
const getTargetSeverity = (targetFrom: number, targetTo: number) => {
  const targetDiff = targetTo - targetFrom
  if (targetDiff > 0) {
    return 'success'
  } else if (targetDiff < 0) {
    return 'danger'
  } else {
    return 'secondary'
  }
}

const stocksTableTexts = computed(() => {
  return {
    title: 'Stock Data Overview',
    description:
      'Quickly view stock tickers, companies, brokerage actions, ratings, and target prices.',
    placeholder: 'Search by ticker, company, or brokerage',
  }
})

const loadStocks = async (
  page: number,
  limit: number,
  field?: string,
  order?: number,
  query?: string,
) => {
  let orderString = order === -1 ? 'desc' : 'asc'
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
    <section class="flex flex-col gap-1 w-full py-4">
      <h1 class="text-3xl font-bold text-start text-white">{{ stocksTableTexts.title }}</h1>
      <p class="text-start text-slate-200">
        {{ stocksTableTexts.description }}
      </p>
    </section>
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
      >
        <Column
          field="ticker"
          header="TICKER"
          class="text-sm font-bold text-black"
          style="width: 5%"
          sortable
        >
          <template #body="{ data }">
            {{ data.ticker }}
          </template>
        </Column>
        <Column
          field="company"
          header="COMPANY"
          class="text-sm text-black"
          style="width: 20%"
          sortable
        >
          <template #body="{ data }">
            {{ data.company }}
          </template>
          ></Column
        >
        <Column
          field="brokerage"
          header="BROKERAGE"
          class="text-sm text-black"
          style="width: 15%"
          sortable
        >
          <template #body="{ data }">
            {{ data.brokerage }}
          </template>
          ></Column
        >
        <Column field="action" header="ACTION" class="text-sm text-black capitalize">
          <template #body="{ data }">
            <span class="">
              {{ data.action }}
            </span>
          </template>
          ></Column
        >
        <Column field="rating" header="RATING" class="text-sm text-black" style="width: 25%">
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
        <Column field="target" header="TARGET" class="text-sm text-black" style="width: 15%">
          <template #body="{ data }">
            <div class="flex flex-row gap-2 items-center">
              <Tag class="capitalize" severity="secondary">$ {{ data.target_from }}</Tag>
              <i
                :class="[
                  data.target_from < data.target_to
                    ? 'pi pi-arrow-up text-green-500'
                    : data.target_from > data.target_to
                      ? 'pi pi-arrow-down text-red-500'
                      : 'pi pi-arrow-right text-gray-500',
                ]"
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
        <Column field="time" header="DATE" class="text-sm text-black" style="width: 10%" sortable>
          <template #body="{ data }">
            {{
              new Date(data.time).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric',
              })
            }}
          </template>
        </Column>
      </DataTable>
    </section>
  </div>
</template>
