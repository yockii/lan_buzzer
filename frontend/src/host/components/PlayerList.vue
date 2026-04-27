<template>
  <div class="text-center py-4 text-slate-400">
    <span class="text-sm">已连接 {{ players.length }} 位选手</span>
    <div class="mt-2 flex flex-wrap justify-center gap-2">
      <span
        v-for="player in players"
        :key="player.id"
        class="inline-flex items-center bg-slate-800 px-3 py-1 rounded-full"
      >
        <span
          class="w-3 h-3 rounded-full mr-2"
          :style="{ backgroundColor: player.color }"
        ></span>
        {{ player.name }}
        <span class="ml-1 text-xs">{{ player.deviceType === 'desktop' ? '💻' : '📱' }}</span>
        <button
          @click="$emit('remove-player', player.id)"
          class="ml-2 text-slate-500 hover:text-red-400 transition-colors"
        >
          ×
        </button>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Player } from '@/shared/types'

defineProps<{
  players: Player[]
}>()

defineEmits<{
  'remove-player': [playerId: string]
}>()
</script>
