import type { StockResponse, Stock } from '@/types/types'
import { CACHE_TIMEOUT } from '@/constants/constants'
import { defineStore } from 'pinia'
import { API_URL } from '@/config/config'
import axios from 'axios'

const INITIAL_PAGE = 1
const INITIAL_LIMIT = 10

export const useLatestStore = defineStore('latest', {
  state: () => ({
    status: 0,
    total: 0,
    stocks: <Stock[]>[],
    loading: false,
    error: null,
    lastFetched: 0,
  }),

  getters: {
    hasData: (state) => state.stocks.length > 0,
    isStale: (state) => {
      if (!state.lastFetched) return true

      return Date.now() - state.lastFetched > CACHE_TIMEOUT
    },
  },

  actions: {
    async fetchData(page: number, limit: number) {
      if (this.loading) return

      this.loading = true
      this.error = null

      try {
        const response = await axios.get(`${API_URL}/stocks`, {
          params: {
            page,
            limit,
          },
        })

        const data = response.data as StockResponse
        this.status = response.status
        this.total = data.length
        this.stocks = data.stocks
        this.lastFetched = Date.now()
      } catch (err: any) {
        this.error = err
        console.error('Error fetching metrics:', err)
      } finally {
        this.loading = false
      }
    },

    async fetchInitialDataIfNeeded() {
      if (!this.hasData || this.isStale) {
        await this.fetchData(INITIAL_PAGE, INITIAL_LIMIT)
      }
      return { stocks: this.stocks, total: this.total }
    },
  },
})
