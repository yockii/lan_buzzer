<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <!-- Top bar with toggle and reset -->
    <div class="px-4 py-2 flex items-center gap-3">
      <!-- Toggle button -->
      <button
        @click="topBarCollapsed = !topBarCollapsed"
        class="text-slate-400 hover:text-white text-sm px-2 py-1 transition-colors"
        :title="topBarCollapsed ? '展开顶部栏' : '收起顶部栏'"
      >
        {{ topBarCollapsed ? '◀' : '▶' }}
      </button>

      <!-- Reset button (only in quiz mode) -->
      <button
        v-if="gameMode === 'quiz'"
        @click="handleReset"
        class="text-slate-500 hover:text-slate-300 text-xs px-2 py-1 transition-colors"
        title="重置题库"
      >
        重置题库
      </button>

      <!-- Collapsible content -->
      <div
        v-show="!topBarCollapsed"
        class="flex-1 flex justify-between items-center"
      >
        <div class="flex items-center gap-4">
          <select
            v-if="allIPs.length > 1"
              v-model="selectedIP"
              class="bg-slate-700 text-white text-sm px-3 py-2 rounded border border-slate-600 focus:outline-none focus:border-slate-500"
          >
            <option v-for="ip in allIPs" :key="ip" :value="ip">
              {{ ip }}
            </option>
          </select>
          <div v-else class="text-sm text-slate-400">{{ displayUrl }}</div>
        </div>
        <div v-if="qrCodeUrl" class="flex flex-col items-center gap-1">
          <img :src="qrCodeUrl" alt="QR Code" class="w-20 h-20 bg-white rounded" />
          <div class="bg-slate-700 text-xs px-2 py-1 rounded">
            手机扫码
          </div>
        </div>
      </div>
    </div>

    <ControlPanel
      :game-state="gameState"
      :game-mode="gameMode"
      @start="handleStart"
      @next="handleNext"
    />

    <!-- Game content area with players on sides -->
    <div class="flex-1 flex gap-4 px-6 overflow-hidden">
      <!-- Left players -->
      <div class="w-1/4 min-w-0">
        <PlayerList
          :players="leftPlayers"
          :is-quiz-mode="gameMode === 'quiz'"
          :ws="ws"
          :player-answers="playerAnswers"
          @remove-player="handleRemovePlayer"
        />
      </div>

      <!-- Center: quiz or buzzer display -->
      <div class="flex-1">
        <QuizDisplay v-if="gameMode === 'quiz'" :ws="ws" @quizQuestion="handleQuizQuestion" />
        <BuzzerDisplay
          v-else
          :game-state="gameState"
          :winner="winner"
        />
      </div>

      <!-- Right players -->
      <div class="w-1/4 min-w-0">
        <PlayerList
          :players="rightPlayers"
          :is-quiz-mode="gameMode === 'quiz'"
          :ws="ws"
          :player-answers="playerAnswers"
          @remove-player="handleRemovePlayer"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Player, Question, QuizAnswerUpdate } from '@/shared/types'
import BuzzerDisplay from './components/BuzzerDisplay.vue'
import QuizDisplay from './components/QuizDisplay.vue'
import ControlPanel from './components/ControlPanel.vue'
import PlayerList from './components/PlayerList.vue'

const gameState = ref<'waiting' | 'ready' | 'locked'>('waiting')
const gameMode = ref<'buzzer' | 'quiz'>('buzzer')
const topBarCollapsed = ref(false)
const players = ref<Player[]>([])
const winner = ref<Player | null>(null)
const playerAnswers = ref<Map<string, QuizAnswerUpdate>>(new Map())
const serverInfo = ref<{ serverUrl: string; localIP: string; allIPs: string[] } | null>(null)
const selectedIP = ref<string>('')

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const allIPs = computed(() => {
  if (serverInfo.value?.allIPs) {
    return serverInfo.value.allIPs
  }
  return []
})

const displayUrl = computed(() => {
  if (selectedIP.value) {
    return `http://${selectedIP.value}:3000`
  }
  if (serverInfo.value) {
    return serverInfo.value.serverUrl
  }
  return 'localhost:3000'
})

const qrCodeUrl = computed(() => {
  if (selectedIP.value) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(displayUrl.value)}`
  }
  if (serverInfo.value) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(serverInfo.value.serverUrl)}`
  }
  return null
})

let ws: WebSocketClient | null = null

const handleQuizQuestion = (payload: Question) => {
  console.log('[DEBUG HostApp] quiz_question received, clearing answers')
  gameMode.value = 'quiz'
  // Clear previous answers when new question starts (use new Map for reactivity)
  playerAnswers.value = new Map()
  console.log('[DEBUG HostApp] answers cleared, new Map size:', playerAnswers.value.size)
}

onMounted(async () => {
  // 获取服务器信息
  try {
    const response = await fetch('/api/info')
    const info = await response.json()
    serverInfo.value = info
    // 默认选择服务器推荐的IP
    if (info.localIP) {
      selectedIP.value = info.localIP
    }
    // 根据服务器模式设置游戏模式
    if (info.mode === 'quiz') {
      gameMode.value = 'quiz'
    }
  } catch (error) {
    console.error('Failed to get server info:', error)
  }

  ws = new WebSocketClient(wsUrl.value)

  ws.on('state_changed', (payload) => {
    gameState.value = payload.state
    // Only reset mode and answers on explicit game reset (not between quiz questions)
    if (payload.state === 'waiting' && gameMode.value === 'buzzer') {
      playerAnswers.value.clear()
    }
  })

  ws.on('player_list', (payload) => {
    players.value = payload.players
  })

  ws.on('buzz_result', (payload) => {
    if (payload.winner) {
      winner.value = payload.winner
    } else {
      winner.value = null
    }
  })


  ws.on('quiz_answer_update', (payload: QuizAnswerUpdate) => {
    playerAnswers.value.set(payload.playerId, payload)
  })

  ws.on('quiz_reset', () => {
    // Clear all answers when question bank is reset
    playerAnswers.value = new Map()
  })

  ws.on('open', () => {
    // 主持人端也发送join消息，但标记为主持人
    ws?.send('join', {
      name: '__HOST__',
      deviceType: 'desktop'
    })
  })

  ws.connect()
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
})

const handleStart = () => {
  if (ws) {
    if (gameMode.value === 'quiz') {
      ws.send('quiz_start', {})
    } else {
      ws.send('start_game', {})
    }
  }
}

const handleNext = () => {
  if (ws) {
    if (gameMode.value === 'quiz') {
      // Check if all answers have been judged
      const hasPending = Array.from(playerAnswers.value.values()).some(
        answer => answer.status === 'pending'
      )
      if (hasPending) {
        if (!confirm('有选手答案尚未判定，确定要进入下一题吗？')) {
          return
        }
      }
      ws.send('quiz_next', {})
    } else {
      ws.send('reset_game', {})
      winner.value = null
    }
  }
}

const handleReset = () => {
  if (ws && confirm('确定要重置题库吗？所有题目将重新开始！')) {
    ws.send('quiz_reset', {})
    // Clear local state
    playerAnswers.value = new Map()
  }
}

const handleRemovePlayer = (playerId: string) => {
  if (ws && confirm('确定要移除这位选手吗？')) {
    ws.send('remove_player', { playerId })
    // Clean up player's quiz answer
    playerAnswers.value.delete(playerId)
  }
}

// Split players into left and right groups
const leftPlayers = computed(() => {
  const mid = Math.ceil(players.value.length / 2)
  return players.value.slice(0, mid)
})

const rightPlayers = computed(() => {
  const mid = Math.ceil(players.value.length / 2)
  return players.value.slice(mid)
})
</script>
