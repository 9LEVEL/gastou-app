<template>
  <div class="category-group">
    <div class="category-group-header" @click="toggle">
      <div class="category-group-bar" :style="{ background: cor }"></div>
      <span class="category-group-name">{{ nome }}</span>
      <span class="category-group-count">{{ count }}</span>
      <svg
        class="category-group-chevron"
        :class="{ collapsed: !expanded }"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </div>
    <div
      class="category-group-body"
      :style="{ maxHeight: expanded ? '2000px' : '0px' }"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  nome: string
  cor: string
  count: number
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const expanded = computed(() => props.modelValue)

function toggle() {
  emit('update:modelValue', !expanded.value)
}
</script>
