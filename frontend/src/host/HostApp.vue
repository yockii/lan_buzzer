<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="flex justify-between items-center px-6 py-4 bg-slate-800">
      <div class="text-sm text-slate-400">{{ displayUrl }}</div>
      <div class="flex items-center gap-4">
        <div v-if="qrCodeUrl" class="relative">
          <img :src="qrCodeUrl" alt="QR Code" class="w-20 h-20 bg-white rounded" />
          <div class="absolute -bottom-1 left-0 bg-slate-700 text-xs px-2 py-1 rounded">
            手机扫码
          </div>
        </div>
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
const serverInfo = ref<{ serverUrl: string; localIP: string } | null>(null)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const displayUrl = computed(() => {
  if (serverInfo.value) {
    return serverInfo.value.serverUrl
  }
  return 'localhost:3000'
})

const qrCodeUrl = computed(() => {
  if (serverInfo.value) {
    // 使用 goqr.me API 生成二维码
    return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(serverInfo.value.serverUrl)}`
  }
  return null
})

let ws: WebSocketClient | null = null

onMounted(async () => {
  // 获取服务器信息
  try {
    const response = await fetch('/api/info')
    serverInfo.value = await response.json()
  } catch (error) {
    console.error('Failed to get server info:', error)
  }

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
