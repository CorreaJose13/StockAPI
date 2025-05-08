import type { MetricsResponse } from '@/types/types'
import { CACHE_TIMEOUT } from '@/constants/constants'
import { API_URL } from '@/config/config'
import { defineStore } from 'pinia'

import axios from 'axios'

export const useMetricsStore = defineStore('metrics', {
  state: () => ({
    status: 0,
    total: 0,
    upgrade: 0,
    downgrade: 0,
    remain: 0,
    loading: false,
    error: null,
    lastFetched: 0,
  }),

  getters: {
    hasData: (state) => state.total !== 0,
    isStale: (state) => {
      if (!state.lastFetched) return true

      return Date.now() - state.lastFetched > CACHE_TIMEOUT
    },
  },

  actions: {
    async fetchData() {
      if (this.loading) return

      this.loading = true
      this.error = null

      try {
        const response = await axios.get(`${API_URL}/metrics`)
        const data = response.data as MetricsResponse
        this.status = response.status
        this.total = data.total_stocks
        this.upgrade = data.positive_change
        this.downgrade = data.negative_change
        this.remain = data.no_change
        this.lastFetched = Date.now()
      } catch (err: any) {
        this.error = err
        console.error('Error fetching metrics:', err)
      } finally {
        this.loading = false
      }
    },

    async fetchIfNeeded() {
      if (!this.hasData || this.isStale) {
        await this.fetchData()
      }
      return {
        total: this.total,
        upgrade: this.upgrade,
        downgrade: this.downgrade,
        remain: this.remain,
      }
    },
  },
})
