<template>
  <div class="bg-slate-800 rounded-2xl p-6">
    <h3 class="text-xl font-bold mb-4">选手列表 ({{ players.length }})</h3>

    <div class="space-y-2">
      <div
        v-for="player in enhancedPlayers"
        :key="player.id"
        class="flex items-center justify-between p-3 bg-slate-700 rounded-lg"
      >
        <div class="flex items-center gap-3">
          <span v-if="player.isWinner" class="text-2xl">👑</span>
          <span
            class="w-3 h-3 rounded-full mr-2"
            :style="{ backgroundColor: player.color }"
          ></span>
          <span class="font-semibold">{{ player.name }}</span>
          <span class="text-xs text-slate-400">({{ player.deviceType === 'mobile' ? '手机' : '电脑' }})</span>
        </div>

        <div class="flex items-center gap-2">
          <!-- Quiz mode: show answer -->
          <template v-if="isQuizMode && player.answer">
            <span
              v-if="player.answerStatus === 'pending'"
              @click="handleJudge(player.id, 'correct')"
              class="px-3 py-1 bg-yellow-500 hover:bg-yellow-600 rounded cursor-pointer transition-colors"
            >
              {{ player.answer }} (待确认)
            </span>
            <span
              v-else-if="player.answerStatus === 'correct'"
              @click="handleJudge(player.id, 'incorrect')"
              class="px-3 py-1 bg-green-500 hover:bg-green-600 rounded cursor-pointer transition-colors"
            >
              ✓{{ player.answer }}
            </span>
            <span
              v-else
              @click="handleJudge(player.id, 'pending')"
              class="px-3 py-1 bg-red-500 hover:bg-red-600 rounded cursor-pointer transition-colors"
            >
              ✗{{ player.answer }}
            </span>
          </template>

          <!-- Buzzer mode: remove button -->
          <button
            v-if="!isQuizMode"
            @click="$emit('remove-player', player.id)"
            class="px-3 py-1 bg-red-500 hover:bg-red-600 rounded text-sm transition-colors"
          >
            移除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Player } from '@/shared/types'

const props = defineProps<{
  players: Player[]
  isQuizMode: boolean
  ws: any
}>()

defineEmits<{
  removePlayer: [playerId: string]
}>()

// Enhanced players with quiz info
const enhancedPlayers = computed(() => {
  return props.players.map(p => ({
    ...p,
    answer: '',
    answerStatus: '' as 'pending' | 'correct' | 'incorrect',
    isWinner: false
  }))
})

const handleJudge = (playerId: string, status: string) => {
  if (props.ws) {
    props.ws.send('quiz_judge', {
      playerId: playerId,
      correct: status === 'correct'
    })
  }
}
</script>
