import { defineStore } from 'pinia'
import { API_URL } from '@/config/config'
import type { StockResponse, Stock } from '@/types/types'
import axios from 'axios'

const INITIAL_PAGE = 1
const INITIAL_LIMIT = 10
const MAX_CACHED_PAGES = 5
const CACHED_TIMEOUT = 60 * 60 * 1000 // 1 hour

type CacheKey = string

interface CachedData {
  stocks: Stock[]
  timestamp: number
}

export const useStocksStore = defineStore('stocks', {
  state: () => ({
    status: 0,
    pagesCache: new Map<CacheKey, CachedData>(),
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

      const staleTimeInMs = CACHED_TIMEOUT
      return Date.now() - state.lastFetched > staleTimeInMs
    },
  },

  actions: {
    generateCacheKey(
      page: number,
      limit: number,
      field?: string,
      order?: string,
      search?: string,
    ): CacheKey {
      return `${page}-${limit}-${field || ''}-${order || ''}-${search || ''}`
    },

    trimCache() {
      if (this.pagesCache.size > MAX_CACHED_PAGES) {
        const entries = Array.from(this.pagesCache.entries()).sort(
          (a, b) => a[1].timestamp - b[1].timestamp,
        )

        const entriesToRemove = entries.slice(0, entries.length - MAX_CACHED_PAGES)
        entriesToRemove.forEach((entry) => {
          this.pagesCache.delete(entry[0])
        })
      }
    },

    async fetchData(page: number, limit: number, field?: string, order?: string, search?: string) {
      const cacheKey = this.generateCacheKey(page, limit, field, order, search)
      const cachedData = this.pagesCache.get(cacheKey)
      const now = Date.now()

      if (cachedData && now - cachedData.timestamp < CACHED_TIMEOUT) {
        this.stocks = cachedData.stocks
        this.lastFetched = cachedData.timestamp
        return
      }

      if (this.loading) return

      this.loading = true
      this.error = null

      try {
        const response = await axios.get(`${API_URL}/stocks`, {
          params: {
            page,
            limit,
            field,
            order,
            search,
          },
        })

        const data = response.data as StockResponse
        this.status = response.status
        this.total = data.length
        this.stocks = data.stocks
        this.lastFetched = now

        this.pagesCache.set(cacheKey, {
          stocks: [...data.stocks],
          timestamp: now,
        })
        console.log('Caching data for key:', cacheKey)

        this.trimCache()
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
