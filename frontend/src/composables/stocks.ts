import type { StockResponse, Stock, MetricsResponse, StockWithScore } from '@/types/types'
import axios from 'axios'
import { API_URL } from '@/config/config'

export const getStocksList = async (
  page: number,
  limit: number,
  field?: string,
  order?: string,
  search?: string,
) => {
  try {
    const response = await axios.get(`${API_URL}/stocks`, {
      params: {
        field,
        order,
        search,
        page,
        limit,
      },
    })
    const data = response.data as StockResponse
    const stockList = data.stocks as Stock[]
    const totalStocks = data.length
    return { status: response.status, stocks: stockList, length: totalStocks }
  } catch (error) {
    console.error('Fetching stocks failed: ' + error)
    throw error
  }
}

export const getStocksMetrics = async () => {
  try {
    const response = await axios.get(`${API_URL}/metrics`, {})
    const data = response.data as MetricsResponse
    return {
      status: response.status,
      total: data.total_stocks,
      upgrade: data.positive_change,
      downgrade: data.negative_change,
      remain: data.no_change,
    }
  } catch (error) {
    console.error('Fetching metrics failed: ' + error)
    throw error
  }
}

export const getStocksAnalysis = async () => {
  try {
    const response = await axios.get(`${API_URL}/analyze`, {})
    const data = response.data as StockWithScore[]
    return {
      status: response.status,
      stocks: data,
    }
  } catch (error) {
    console.error('Fetching analysis failed: ' + error)
    throw error
  }
}
