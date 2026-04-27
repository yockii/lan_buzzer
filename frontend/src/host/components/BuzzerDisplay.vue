<template>
  <div class="flex-1 flex items-center justify-center px-8">
    <div class="w-full text-center py-16 bg-slate-800 rounded-2xl">
      <div v-if="winner" class="space-y-4">
        <div class="text-7xl font-bold" :style="{ color: winner.color }">
          {{ winner.name }}
        </div>
        <div class="text-2xl text-slate-300">
          🎉 抢到了！
        </div>
      </div>
      <div v-else class="space-y-4">
        <div class="text-5xl font-bold text-slate-400">
          {{ stateText }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Player } from '@/shared/types'

const props = defineProps<{
  gameState: 'waiting' | 'ready' | 'locked'
  winner: Player | null
}>()

const stateText = computed(() => {
  switch (props.gameState) {
    case 'waiting':
      return '等待开始'
    case 'ready':
      return '准备抢答...'
    case 'locked':
      return '已锁定'
  }
})
</script>
