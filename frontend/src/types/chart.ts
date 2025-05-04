export interface ChartDataset {
  label: string
  fill: boolean
  borderColor: string
  yAxisID: string
  tension: number
  data: number[]
}

export interface ChartData {
  labels: string[]
  datasets: ChartDataset[]
}

export interface ScaleOptions {
  type?: string
  display?: boolean
  position?: string
  ticks: {
    color: string
  }
  grid: {
    color: string
    drawOnChartArea?: boolean
  }
}

export interface ChartOptions {
  stacked: boolean
  maintainAspectRatio: boolean
  aspectRatio: number
  plugins: {
    legend: {
      labels: {
        color: string
      }
    }
  }
  scales: {
    x: ScaleOptions
    y: ScaleOptions
    y1: ScaleOptions
  }
}
