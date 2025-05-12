<script setup lang="ts">
import { priceAbsDiff, formatDateShort, modalDt } from '@/utils/stock'
import { getRatingSeverity, formatPrice, tagDt, validateImageSize } from '@/utils/stock'
import StockModal from './StockModal.vue'
import { computed, ref, onMounted } from 'vue'

const showModal = ref(false)

const imageError = ref(false)
const imageUrl = ref('')

const props = defineProps({
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

const getHeaderTag = computed(() => {
  const diff = priceAbsDiff(props.targetFrom, props.targetTo)
  if (diff > 0) {
    return {
      text: 'Upgrade',
      severity: 'success',
    }
  } else if (diff < 0) {
    return {
      text: 'Downgrade',
      severity: 'danger',
    }
  } else {
    return {
      text: 'Remain',
      severity: 'secondary',
    }
  }
})

const getPriceColor = computed(() => {
  const diff = priceAbsDiff(props.targetFrom, props.targetTo)
  if (diff > 0) {
    return 'text-green-500'
  } else if (diff < 0) {
    return 'text-red-500'
  } else {
    return 'text-gray-500'
  }
})

const priceSection = computed(() => {
  return {
    label: 'Price:',
    fromValue: formatPrice(props.targetFrom),
    toValue: formatPrice(props.targetTo),
  }
})

const ratingSection = computed(() => {
  return {
    label: 'Rating:',
    toValue: props.ratingTo,
  }
})

const timeSection = computed(() => {
  return {
    title: 'Last update:',
    value: formatDateShort(props.time),
  }
})

onMounted(async () => {
  try {
    const url = await validateImageSize(props.ticker)
    if (url) {
      imageUrl.value = url
    }
  } catch (error) {
    imageError.value = true
  }
})
</script>
<template>
  <div
    class="flex h-full min-w-3xs cursor-pointer flex-col justify-between rounded-lg bg-white p-4 transition duration-300 ease-in-out hover:-translate-y-1 hover:scale-105 dark:bg-stone-900"
    @click="showModal = true"
  >
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <img v-if="imageUrl && !imageError" :src="imageUrl" class="size-8" />
        <span class="text-xl font-bold text-black dark:text-white"> {{ props.ticker }} </span>
        <Tag :severity="getHeaderTag.severity" :dt="tagDt">
          {{ getHeaderTag.text }}
        </Tag>
      </div>
      <span class="text-md mb-3 w-3/4 text-gray-700 capitalize dark:text-slate-200">{{
        props.company
      }}</span>
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex flex-row items-center gap-2 text-lg">
        <span class="font-semibold text-black dark:text-white">{{ priceSection.label }}</span>
        <div class="flex items-center gap-2">
          <span class="text-slate-500 capitalize line-through dark:text-slate-300">
            {{ priceSection.fromValue }}
          </span>
          <i class="pi pi-arrow-right" style="font-size: 0.75rem"></i>
          <span class="text-black capitalize" :class="getPriceColor">
            {{ priceSection.toValue }}
          </span>
        </div>
      </div>

      <div class="flex flex-row items-center gap-2 text-lg">
        <span class="font-semibold text-black dark:text-white">{{ ratingSection.label }}</span>
        <Tag :severity="getRatingSeverity(props.ratingTo)" class="size-sm capitalize" :dt="tagDt">
          {{ ratingSection.toValue }}
        </Tag>
      </div>

      <div class="text-sm text-gray-500 dark:text-gray-300">
        {{ timeSection.title }} {{ timeSection.value }}
      </div>
    </div>
  </div>
  <Dialog
    v-model:visible="showModal"
    modal
    :header="`${props.ticker} details`"
    :dt="modalDt"
    :dismissableMask="true"
  >
    <StockModal
      :srcUrl="imageUrl"
      :ticker="props.ticker"
      :action="props.action"
      :company="props.company"
      :targetFrom="props.targetFrom"
      :targetTo="props.targetTo"
      :ratingFrom="props.ratingFrom"
      :ratingTo="props.ratingTo"
      :brokerage="props.brokerage"
      :time="props.time"
    ></StockModal>
  </Dialog>
</template>
