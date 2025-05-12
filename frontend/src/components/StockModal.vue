<script setup lang="ts">
import { priceAbsDiff, pricePercDiff, formatDateLong, formatAction } from '@/utils/stock'
import { modalTagDt, modalTagDtXl, formatPrice, getRatingSeverity } from '@/utils/stock'
import { validateImageSize, brandImageUrl } from '@/utils/stock'
import { computed, onMounted, ref } from 'vue'
import StockChart from './StockChart.vue'

const props = defineProps({
  srcUrl: { type: String },
  ticker: { type: String, required: true },
  action: { type: String, required: true },
  company: { type: String, required: true },
  targetFrom: { type: Number, required: true },
  targetTo: { type: Number, required: true },
  ratingFrom: { type: String, required: true },
  ratingTo: { type: String, required: true },
  brokerage: { type: String, required: true },
  time: { type: String, required: true },
})

const imageUrl = ref('')
const imageError = ref(false)

const modalTexts = computed(() => {
  return {
    price: '💰 Price target',
    rating: '📝 Rating',
    action: 'Action:',
  }
})

const getNumberSeverity = (number: number) => {
  if (number > 0) {
    return 'success'
  } else if (number < 0) {
    return 'danger'
  } else {
    return 'secondary'
  }
}
const getAbsSeverity = (from: number, to: number) => {
  const absDiff = Number(priceAbsDiff(from, to))
  return getNumberSeverity(absDiff)
}

const getPercSeverity = (from: number, to: number) => {
  const percDiff = Number(pricePercDiff(from, to))
  return getNumberSeverity(percDiff)
}

const getPerc = computed(() => {
  const percDiff = pricePercDiff(props.targetFrom, props.targetTo)
  return `${Number(percDiff).toFixed(2)}%`
})

onMounted(async () => {
  if (!props.srcUrl) {
    try {
      const url = await validateImageSize(props.ticker)
      if (url) {
        imageUrl.value = url
      }
    } catch (error) {
      imageError.value = true
    }
  } else {
    imageUrl.value = props.srcUrl
  }
})
</script>
<template>
  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between">
      <div class="flex flex-row items-center gap-2">
        <img v-if="imageUrl && !imageError" :src="imageUrl" class="size-14" />
        <span class="text-2xl font-bold text-black dark:text-white"> {{ props.ticker }} </span>
      </div>
      <button
        class="text-md rounded-full bg-gray-400 px-4 py-1 font-semibold text-white capitalize"
      >
        {{ props.brokerage }}
      </button>
    </div>
    <span class="mb-3 w-3/4 text-lg text-gray-700 capitalize dark:text-slate-200">{{
      props.company
    }}</span>
  </div>

  <div class="grid grid-cols-2 gap-4">
    <div class="rounded-xl bg-stone-200 p-4">
      <div class="flex flex-col gap-2">
        <span class="text-md font-medium text-black"> {{ modalTexts.price }}</span>
        <section class="flex items-center gap-2">
          <span class="text-md text-gray-600 capitalize line-through">{{
            formatPrice(props.targetFrom)
          }}</span>
          <span class="text-2xl font-bold text-black capitalize">{{
            formatPrice(props.targetTo)
          }}</span>
          <div class="flex gap-2">
            <Tag
              :value="formatPrice(priceAbsDiff(props.targetFrom, props.targetTo))"
              :severity="getAbsSeverity(props.targetFrom, props.targetTo)"
              :dt="modalTagDt"
            ></Tag>
            <Tag
              :value="getPerc"
              :severity="getPercSeverity(props.targetFrom, props.targetTo)"
              :dt="modalTagDt"
            >
            </Tag>
          </div>
        </section>
      </div>
    </div>

    <div class="rounded-xl bg-stone-200 p-4">
      <div class="flex flex-col gap-2">
        <span class="text-md font-medium text-black"> {{ modalTexts.rating }}</span>
        <section class="flex items-center gap-2">
          <span class="text-md text-gray-600 capitalize line-through">{{ props.ratingFrom }}</span>
          <Tag
            :severity="getRatingSeverity(props.ratingTo)"
            class="size-sm capitalize"
            :dt="modalTagDtXl"
          >
            {{ props.ratingTo }}
          </Tag>
        </section>
      </div>
    </div>
  </div>

  <StockChart :ticker="props.ticker" class="my-6" />

  <div class="mt-6 flex items-center justify-between gap-2">
    <div class="flex flex-row gap-1">
      <p>{{ modalTexts.action }}</p>
      <span class="text-gray-600"> {{ formatAction(props.action) }}</span>
    </div>
    <span class="text-gray-600">🗓️ {{ formatDateLong(props.time) }}</span>
  </div>
</template>
