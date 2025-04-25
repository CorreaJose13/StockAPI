export const stocks = [
  {
    id: 1,
    ticker: 'AAPL',
    targetFrom: 180.5,
    targetTo: 215.75,
    company: 'Apple Inc.',
    action: 'upgraded by',
    brokerage: 'Morgan Stanley',
    ratingFrom: 'Hold',
    ratingTo: 'Buy',
    time: '2025-01-15T09:30:00.000Z',
  },
  {
    id: 2,
    ticker: 'MSFT',
    targetFrom: 320.0,
    targetTo: 375.25,
    company: 'Microsoft Corporation',
    action: 'reiterated by',
    brokerage: 'Goldman Sachs',
    ratingFrom: 'Buy',
    ratingTo: 'Buy',
    time: '2025-01-14T14:45:00.000Z',
  },
  {
    id: 3,
    ticker: 'TSLA',
    targetFrom: 195.0,
    targetTo: 165.5,
    company: 'Tesla, Inc.',
    action: 'downgraded by',
    brokerage: 'JP Morgan',
    ratingFrom: 'Neutral',
    ratingTo: 'Sell',
    time: '2025-01-13T11:15:00.000Z',
  },
]

export function useStockData() {
  return stocks
}
