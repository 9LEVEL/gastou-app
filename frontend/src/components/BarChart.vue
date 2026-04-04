<template>
  <div class="bar-chart">
    <div v-for="(item, index) in items" :key="index" class="bar-chart-row">
      <div class="bar-chart-label text-sm truncate">{{ item.label }}</div>
      <div class="bar-chart-bars">
        <div class="bar-chart-track">
          <div
            class="bar-chart-fill bar-planned"
            :style="{
              width: maxValue > 0 ? (item.planejado / maxValue * 100) + '%' : '0%',
              background: item.cor + '44'
            }"
          ></div>
        </div>
        <div class="bar-chart-track">
          <div
            class="bar-chart-fill bar-real"
            :style="{
              width: maxValue > 0 ? (item.real / maxValue * 100) + '%' : '0%',
              background: item.cor
            }"
          ></div>
        </div>
      </div>
      <div class="bar-chart-values text-xs">
        <div class="text-secondary">R$ {{ format(item.planejado) }}</div>
        <div class="num" style="font-weight: 600;">R$ {{ format(item.real) }}</div>
      </div>
    </div>
    <div class="bar-chart-legend text-xs text-secondary" style="margin-top: 12px; display: flex; gap: 16px;">
      <div class="flex items-center gap-8">
        <span style="width: 12px; height: 8px; border-radius: 2px; background: rgba(0,0,0,0.15); display: inline-block;"></span>
        Planejado
      </div>
      <div class="flex items-center gap-8">
        <span style="width: 12px; height: 8px; border-radius: 2px; background: var(--color-primary); display: inline-block;"></span>
        Real
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface BarItem {
  label: string
  planejado: number
  real: number
  cor: string
}

const props = defineProps<{
  items: BarItem[]
}>()

const maxValue = computed(() => {
  let max = 0
  for (const item of props.items) {
    if (item.planejado > max) max = item.planejado
    if (item.real > max) max = item.real
  }
  return max
})

function format(value: number): string {
  return value.toFixed(2).replace('.', ',')
}
</script>

<style scoped>
.bar-chart-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.bar-chart-label {
  width: 80px;
  flex-shrink: 0;
  font-weight: 500;
}

.bar-chart-bars {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.bar-chart-track {
  height: 10px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 5px;
  overflow: hidden;
}

.bar-chart-fill {
  height: 100%;
  border-radius: 5px;
  transition: width 0.5s ease;
  min-width: 2px;
}

.bar-chart-values {
  width: 70px;
  flex-shrink: 0;
  text-align: right;
}
</style>
