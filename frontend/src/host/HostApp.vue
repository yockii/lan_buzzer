<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="flex justify-between items-center px-6 py-4 bg-slate-800">
      <div class="text-sm text-slate-400">{{ serverUrl }}</div>
      <div class="bg-slate-700 px-4 py-2 rounded-lg text-sm">
        📱 二维码（手机扫码加入）
      </div>
    </div>

    <BuzzerDisplay
      :game-state="gameState"
      :winner="winner"
    />

    <ControlPanel
      :game-state="gameState"
      @start="handleStart"
      @reset="handleReset"
    />

    <PlayerList
      :players="players"
      @remove-player="handleRemovePlayer"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Player } from '@/shared/types'
import BuzzerDisplay from './components/BuzzerDisplay.vue'
import ControlPanel from './components/ControlPanel.vue'
import PlayerList from './components/PlayerList.vue'

const gameState = ref<'waiting' | 'ready' | 'locked'>('waiting')
const players = ref<Player[]>([])
const winner = ref<Player | null>(null)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const serverUrl = computed(() => {
  return window.location.host
})

let ws: WebSocketClient | null = null

onMounted(() => {
  ws = new WebSocketClient(wsUrl.value)

  ws.on('state_changed', (payload) => {
    gameState.value = payload.state
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

  ws.connect()
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
})

const handleStart = () => {
  if (ws) {
    ws.send('start_game', {})
  }
}

const handleReset = () => {
  if (ws) {
    ws.send('reset_game', {})
    winner.value = null
  }
}

const handleRemovePlayer = (playerId: string) => {
  if (ws && confirm('确定要移除这位选手吗？')) {
    ws.send('remove_player', { playerId })
  }
}
</script>
