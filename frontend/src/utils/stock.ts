import { ref } from 'vue'

export const getRatingSeverity = (rating: string) => {
  switch (rating) {
    case 'sell':
      return 'danger'
    case 'buy':
      return 'success'
    case 'outperform':
      return 'success'
    case 'underperform':
      return 'warn'
    case 'hold':
      return 'info'
    default:
      return 'contrast'
  }
}

export const getTargetSeverity = (targetFrom: number, targetTo: number) => {
  const targetDiff = targetTo - targetFrom
  if (targetDiff > 0) {
    return 'success'
  } else if (targetDiff < 0) {
    return 'danger'
  } else {
    return 'secondary'
  }
}

export const getTargetArrow = (targetFrom: number, targetTo: number) => {
  return targetFrom < targetTo
    ? 'pi pi-arrow-up text-green-500'
    : targetFrom > targetTo
      ? 'pi pi-arrow-down text-red-500'
      : 'pi pi-arrow-right text-gray-500'
}

export const priceAbsDiff = (targetFrom: number, targetTo: number) => {
  const absDiff = targetTo - targetFrom
  return absDiff
}

export const pricePercDiff = (targetFrom: number, targetTo: number) => {
  if (targetFrom === 0) {
    return 'N/A'
  }
  const percDiff = (priceAbsDiff(targetFrom, targetTo) * 100) / targetFrom
  return percDiff
}

export const formatDateShort = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export const formatDateLong = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })
}

export const formatPrice = (price: number) => {
  return price.toLocaleString('en-US', {
    style: 'currency',
    currency: 'USD',
  })
}

export const modalDt = ref({
  root: {
    background: '{slate.200}',
  },
})

export const tagDt = ref({
  root: { fontSize: '1rem' },
})

export const tableDt = ref({
  root: {
    background: '{white}',
  },
})

export const rowClass = () => {
  return 'cursor-pointer'
}
