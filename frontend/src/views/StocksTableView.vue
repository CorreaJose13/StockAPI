<script setup lang="ts">
import { useStockData } from '@/composables/stocks'
import { computed } from 'vue'
import SearchIcon from '@/components/icons/SearchIcon.vue'

const stocks = useStockData()
const tableHeaders = computed(() => {
  return [
    'Ticker',
    'Company',
    'Brokerage',
    'Action',
    'Rating from',
    'Rating to',
    'Target from',
    'Target to',
  ]
})
const tableTexts = computed(() => {
  return {
    title: 'Projects with Invoices',
    subtitle: 'Overview of the current activities.',
    searchPlaceholder: 'Search for stock...',
  }
})
</script>
<template>
  <div class="max-w-screen-xl max-h-screen mx-auto p-4">
    <div class="w-full flex justify-between items-center mb-3 mt-1">
      <div>
        <h3 class="text-lg font-semibold text-black">{{ tableTexts.title }}</h3>
        <p class="text-gray-500">{{ tableTexts.subtitle }}</p>
      </div>
      <div class="ml-3">
        <div class="w-full max-w-md min-w-xs">
          <div class="relative">
            <input
              class="bg-white w-full pr-11 h-10 pl-3 py-2 placeholder:text-gray-400 text-black text-sm border border-gray-200 rounded transition duration-200 ease focus:outline-none focus:border-gray-500 hover:border-gray-400 shadow-sm focus:shadow-md"
              :placeholder="tableTexts.searchPlaceholder"
            />
            <SearchIcon />
          </div>
        </div>
      </div>
    </div>

    <div
      class="relative flex flex-col w-full h-full overflow-scroll text-gray-700 bg-white shadow-md rounded-lg bg-clip-border"
    >
      <table class="w-full text-left table-auto min-w-max">
        <thead>
          <tr>
            <th
              v-for="header in tableHeaders"
              :key="header"
              class="p-4 border-b border-slate-200 bg-slate-50"
            >
              <p class="text-sm font-normal leading-none uppercase text-gray-700">{{ header }}</p>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="stock in stocks"
            :key="stock.id"
            class="hover:bg-slate-50 border-b border-slate-200"
          >
            <td class="p-4 py-5">
              <span class="font-semibold text-sm text-black">{{ stock.ticker }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">{{ stock.company }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">{{ stock.brokerage }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">{{ stock.brokerage }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">{{ stock.ratingFrom }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">{{ stock.ratingTo }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">${{ stock.targetFrom }}</span>
            </td>
            <td class="p-4 py-5">
              <span class="block text-sm text-gray-500">${{ stock.targetTo }}</span>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="flex justify-between items-center px-4 py-3">
        <div class="text-sm text-slate-500">Showing <b>1-5</b> of 45</div>
        <div class="flex space-x-1">
          <button
            class="px-3 py-1 min-w-9 min-h-9 text-sm font-normal text-slate-500 bg-white border border-slate-200 rounded hover:bg-slate-50 hover:border-slate-400 transition duration-200 ease"
          >
            Prev
          </button>
          <button
            class="px-3 py-1 min-w-9 min-h-9 text-sm font-normal text-white bg-slate-800 border border-slate-800 rounded hover:bg-slate-600 hover:border-slate-600 transition duration-200 ease"
          >
            1
          </button>
          <button
            class="px-3 py-1 min-w-9 min-h-9 text-sm font-normal text-slate-500 bg-white border border-slate-200 rounded hover:bg-slate-50 hover:border-slate-400 transition duration-200 ease"
          >
            2
          </button>
          <button
            class="px-3 py-1 min-w-9 min-h-9 text-sm font-normal text-slate-500 bg-white border border-slate-200 rounded hover:bg-slate-50 hover:border-slate-400 transition duration-200 ease"
          >
            3
          </button>
          <button
            class="px-3 py-1 min-w-9 min-h-9 text-sm font-normal text-slate-500 bg-white border border-slate-200 rounded hover:bg-slate-50 hover:border-slate-400 transition duration-200 ease"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
