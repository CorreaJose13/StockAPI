import { defineStore } from 'pinia'
import { API_URL } from '@/config/config'
import type { AnalysisResponse, StockWithScore } from '@/types/types'
import axios from 'axios'

export const useAnalysisStore = defineStore('analysis', {
  state: () => ({
    stocks: <StockWithScore[]>[],
    status: 0,
    loading: false,
    error: null,
    lastFetched: 0,
  }),

  getters: {
    hasData: (state) => state.stocks.length > 0,
    isStale: (state) => {
      if (!state.lastFetched) return true

      const staleTimeInMs = 60 * 60 * 1000 // 1 hour
      return Date.now() - state.lastFetched > staleTimeInMs
    },
  },

  actions: {
    async fetchData() {
      if (this.loading) return

      this.loading = true
      this.error = null

      try {
        const response = await axios.get(`${API_URL}/analyze`)
        const data = response.data as AnalysisResponse
        this.stocks = data.top_stocks
        this.status = response.status
        this.lastFetched = Date.now()
      } catch (err: any) {
        this.error = err
        console.error('Error fetching analysis:', err)
      } finally {
        this.loading = false
      }
    },

    async fetchIfNeeded() {
      if (!this.hasData || this.isStale) {
        await this.fetchData()
      }
      return this.stocks
    },
  },
})
