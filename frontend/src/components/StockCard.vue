<script setup lang="ts">
import { computed } from 'vue'
import { priceAbsDiff, formatDate } from '@/utils/stock'
import { ref } from 'vue'
import StockModal from './StockModal.vue'

const visible = ref(false)

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

const modalStyle = ref({
  root: {
    background: '{slate.200}',
  },
})

const getHeaderTag = computed(() => {
  const diff = priceAbsDiff(props.targetFrom, props.targetTo)
  if (diff > 0) {
    return {
      text: 'Up',
      severity: 'success',
    }
  } else if (diff < 0) {
    return {
      text: 'Down',
      severity: 'danger',
    }
  } else {
    return {
      text: 'Neutral',
      severity: 'secondary',
    }
  }
})

const cardSections = computed(() => [
  {
    label: 'Price:',
    fromValue: `$${props.targetFrom}`,
    toValue: `$${props.targetTo}`,
  },
  {
    label: 'Rating:',
    fromValue: props.ratingFrom,
    toValue: props.ratingTo,
    capitalize: true,
  },
])
</script>
<template>
  <div
    class="bg-white rounded-lg shadow-md p-4 min-w-3xs flex flex-col gap-1"
    @click="visible = true"
  >
    <div class="flex items-center gap-2">
      <span class="font-bold text-lg text-black"> {{ props.ticker }} </span>
      <Tag :severity="getHeaderTag.severity">
        {{ getHeaderTag.text }}
      </Tag>
    </div>
    <span class="text-gray-700 text-sm mb-3">{{ props.company }}</span>

    <div v-for="(section, index) in cardSections" :key="index" class="flex flex-col gap-1">
      <span class="text-black font-semibold">{{ section.label }}</span>
      <div class="flex items-center gap-1">
        <span class="text-gray-500 line-through" :class="{ capitalize: section.capitalize }">
          {{ section.fromValue }}
        </span>
        <i class="pi pi-arrow-right" style="font-size: 0.75rem"></i>
        <span class="text-black" :class="{ capitalize: section.capitalize }">
          {{ section.toValue }}
        </span>
      </div>
    </div>

    <div class="text-xs text-gray-500">{{ props.brokerage }} • {{ formatDate(props.time) }}</div>
  </div>
  <Dialog v-model:visible="visible" modal header="Stock details" class="min-w-xl" :dt="modalStyle">
    <StockModal
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
