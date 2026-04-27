<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="text-center py-6 bg-slate-800">
      <div class="text-2xl font-bold">👋 {{ playerName }}</div>
      <div class="mt-2">
        <span
          class="inline-block px-4 py-1 rounded-full text-sm"
          :class="connected ? 'bg-green-500' : 'bg-red-500'"
        >
          {{ connected ? '✓ 已连接' : '✗ 未连接' }}
        </span>
      </div>
    </div>

    <div class="flex-1 flex items-center justify-center px-8">
      <div class="w-full text-center py-16 bg-slate-800 rounded-2xl">
        <div class="text-4xl font-bold">
          {{ statusText }}
        </div>
      </div>
    </div>

    <div class="p-8">
      <button
        @click="handleBuzz"
        @keydown.space.prevent="handleBuzz"
        @keydown.enter.prevent="handleBuzz"
        :disabled="!canBuzz"
        class="w-full h-32 rounded-2xl font-bold text-3xl transition-all disabled:cursor-not-allowed"
        :class="canBuzz ? 'bg-red-500 hover:bg-red-600 active:scale-95' : 'bg-gray-700 opacity-50'"
      >
        抢答！
      </button>
      <div class="text-center mt-4 text-slate-400 text-sm">
        按空格键 或 点击按钮
      </div>
    </div>

    <div
      v-if="showEarlyBuzzModal"
      class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-8 z-50"
    >
      <div class="bg-slate-800 rounded-2xl p-8 max-w-md w-full text-center">
        <div class="text-5xl mb-4">⚠️</div>
        <div class="text-2xl font-bold mb-4">请等待主持人开始！</div>
        <div class="text-slate-400 mb-6">抢跑会受到时间惩罚</div>
        <button
          @click="closeEarlyBuzzModal"
          class="w-full px-8 py-4 bg-red-500 hover:bg-red-600 rounded-lg text-xl font-semibold transition-colors"
        >
          我知道了
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'

const props = defineProps<{
  playerName: string
}>()

const gameState = ref<'waiting' | 'ready' | 'locked'>('waiting')
const connected = ref(false)
const showEarlyBuzzModal = ref(false)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const statusText = computed(() => {
  switch (gameState.value) {
    case 'waiting':
      return '⏳ 等待主持人开始...'
    case 'ready':
      return '🚀 准备抢答！'
    case 'locked':
      return '🔒 本轮已结束'
  }
})

const canBuzz = computed(() => {
  return gameState.value === 'ready' && connected.value
})

let ws: WebSocketClient | null = null
const deviceType = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(navigator.userAgent)
  ? 'mobile'
  : 'desktop'

onMounted(() => {
  const storedName = sessionStorage.getItem('playerName') || props.playerName

  ws = new WebSocketClient(wsUrl.value)

  ws.on('state_changed', (payload) => {
    gameState.value = payload.state
  })

  ws.on('early_buzz_warning', () => {
    showEarlyBuzzModal.value = true
  })

  ws.on('open', () => {
    connected.value = true
    ws?.send('join', {
      name: storedName,
      deviceType: deviceType
    })
  })

  ws.connect()

  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
  window.removeEventListener('keydown', handleGlobalKeydown)
})

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (e.code === 'Space' || e.code === 'Enter') {
    if (canBuzz.value) {
      handleBuzz()
    }
  }
}

const handleBuzz = () => {
  if (!canBuzz.value || !ws) {
    return
  }

  ws.send('buzz', {
    timestamp: Date.now()
  })
}

const closeEarlyBuzzModal = () => {
  showEarlyBuzzModal.value = false
}
</script>
