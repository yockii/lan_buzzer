<template>
  <div class="bg-slate-800 rounded-2xl p-4 h-full flex flex-col">
    <h3 class="text-lg font-bold mb-3">选手 ({{ players.length }})</h3>

    <div class="flex-1 overflow-y-auto space-y-2">
      <div
        v-for="player in enhancedPlayers"
        :key="player.id"
        class="flex flex-col gap-2 p-3 bg-slate-700 rounded-lg"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span v-if="player.isWinner" class="text-lg">👑</span>
            <span
              class="w-2 h-2 rounded-full"
              :style="{ backgroundColor: player.color }"
            ></span>
            <span class="font-semibold text-sm">{{ player.name }}</span>
          </div>
          <span class="text-xs text-slate-400">({{ player.deviceType === 'mobile' ? '手机' : '电脑' }})</span>
        </div>

        <div class="flex items-center gap-1">
          <!-- Quiz mode: show answer -->
          <template v-if="isQuizMode && player.answer">
            <span
              v-if="player.answerStatus === 'pending'"
              @click="handleJudge(player.id, 'correct')"
              class="px-2 py-1 bg-yellow-500 hover:bg-yellow-600 rounded cursor-pointer transition-colors text-xs flex-1 text-center truncate"
              title="点击标记为正确"
            >
              {{ player.answer }}
            </span>
            <span
              v-else-if="player.answerStatus === 'correct'"
              @click="handleJudge(player.id, 'incorrect')"
              class="px-2 py-1 bg-green-500 hover:bg-green-600 rounded cursor-pointer transition-colors text-xs flex-1 text-center truncate"
              title="点击修改为错误"
            >
              ✓{{ player.answer }}
            </span>
            <span
              v-else
              @click="handleJudge(player.id, 'correct')"
              class="px-2 py-1 bg-red-500 hover:bg-red-600 rounded cursor-pointer transition-colors text-xs flex-1 text-center truncate"
              title="点击修改为正确"
            >
              ✗{{ player.answer }}
            </span>
          </template>

          <!-- Buzzer mode: remove button -->
          <button
            v-if="!isQuizMode"
            @click="$emit('remove-player', player.id)"
            class="px-2 py-1 bg-red-500 hover:bg-red-600 rounded text-xs transition-colors flex-1"
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
import type { Player, QuizAnswerUpdate } from '@/shared/types'

const props = defineProps<{
  players: Player[]
  isQuizMode: boolean
  ws: any
  playerAnswers: Map<string, QuizAnswerUpdate>
}>()

defineEmits<{
  removePlayer: [playerId: string]
}>()

// Enhanced players with quiz info
const enhancedPlayers = computed(() => {
  // Find the earliest correct answer timestamp
  const correctAnswers = Array.from(props.playerAnswers.values())
    .filter(a => a.status === 'correct')
    .map(a => a.timestamp)
    .filter(t => t > 0)

  const earliestCorrectTimestamp = correctAnswers.length > 0
    ? Math.min(...correctAnswers)
    : Infinity

  return props.players.map(p => {
    const answer = props.playerAnswers.get(p.id)
    const isWinner = answer?.status === 'correct' &&
                     answer?.timestamp === earliestCorrectTimestamp

    return {
      ...p,
      answer: answer?.answer || '',
      answerStatus: answer?.status || 'pending',
      isWinner: isWinner || false
    }
  })
})

const handleJudge = (playerId: string, status: string) => {
  console.log('[DEBUG PlayerList] handleJudge called:', { playerId, status, ws: !!props.ws })
  if (props.ws) {
    const payload = {
      playerId: playerId,
      correct: status === 'correct'
    }
    console.log('[DEBUG PlayerList] Sending quiz_judge:', payload)
    props.ws.send('quiz_judge', payload)
    console.log('[DEBUG PlayerList] Message sent')
  } else {
    console.log('[DEBUG PlayerList] ws is null/undefined')
  }
}
</script>
