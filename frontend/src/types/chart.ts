export interface ChartData {
  open: number
  high: number
  low: number
  close: number
  volume: number
  date: string
}

export interface ChartResponse {
  time_series: ChartData[]
}
