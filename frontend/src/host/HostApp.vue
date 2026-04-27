<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="flex justify-between items-center px-6 py-4">
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

onMounted(async () => {
  // 获取服务器信息
  try {
    const response = await fetch('/api/info')
    serverInfo.value = await response.json()
    // 默认选择服务器推荐的IP
    if (serverInfo.value.localIP) {
      selectedIP.value = serverInfo.value.localIP
    }
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
