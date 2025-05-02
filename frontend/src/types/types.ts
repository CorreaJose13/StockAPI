export interface Stock {
  ticker: string
  target_from: number
  target_to: number
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
}

export interface StockWithScore extends Stock {
  score: number
}

export interface MetricsResponse {
  total_stocks: number
  positive_change: number
  negative_change: number
  no_change: number
}

export interface StockResponse {
  stocks: Stock[]
  length: number
}

export interface StockWithScoreResponse {
  stocks: StockWithScore[]
}

export interface ErrorResponse {
  message: string
}
