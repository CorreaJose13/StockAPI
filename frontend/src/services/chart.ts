import { API_URL } from '@/config/config'
import axios from 'axios'
import type { ChartResponse } from '@/types/chart'

export const fetchChartData = async (ticker: string) => {
  try {
    const response = await axios.get(`${API_URL}/chart`, {
      params: {
        ticker,
      },
    })
    const data = response.data as ChartResponse
    const timeSeries = data.time_series
    return timeSeries
  } catch (err: any) {
    console.error('Error fetching chart:', err)
  }
}
