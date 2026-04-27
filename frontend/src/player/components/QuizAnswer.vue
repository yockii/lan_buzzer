<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="text-center py-6 bg-slate-800">
      <div class="text-2xl font-bold">{{ playerName }}</div>
      <div class="mt-2">
        <span
          class="inline-block px-4 py-1 rounded-full text-sm"
          :class="connected ? 'bg-green-500' : 'bg-red-500'"
        >
          {{ connected ? '✓ 已连接' : '✗ 未连接' }}
        </span>
      </div>
    </div>

    <!-- Waiting for question -->
    <div v-if="!currentQuestion" class="flex-1 flex items-center justify-center px-8">
      <div class="text-center text-slate-400">
        <div class="text-4xl mb-4">⏳</div>
        <div class="text-xl">等待主持人开始...</div>
      </div>
    </div>

    <!-- Answer submitted -->
    <div v-else-if="answerSubmitted" class="flex-1 flex items-center justify-center px-8">
      <div class="text-center text-slate-400">
        <div class="text-4xl mb-4">✓</div>
        <div class="text-xl">等待下一题</div>
      </div>
    </div>

    <!-- Show question and input -->
    <div v-else class="flex-1 flex flex-col items-center justify-center px-8 py-8">
      <div class="text-3xl font-bold mb-8 text-center">
        {{ currentQuestion.question }}
      </div>

      <!-- Multiple choice / true false -->
      <div
        v-if="currentQuestion.options.length > 0"
        class="w-full max-w-md grid gap-4"
        :class="currentQuestion.options.length === 2 ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <button
          v-for="(option, index) in currentQuestion.options"
          :key="index"
          @click="handleAnswer(option)"
          class="w-full px-6 py-4 bg-slate-800 hover:bg-slate-700 rounded-lg text-lg text-left transition-colors"
        >
          {{ option }}
        </button>
      </div>

      <!-- Open-ended -->
      <div v-else class="w-full max-w-md">
        <input
          v-model="openAnswer"
          @keyup.enter="handleAnswer(openAnswer)"
          type="text"
          placeholder="输入你的答案..."
          class="w-full px-6 py-4 text-xl bg-slate-800 border-2 border-slate-700 rounded-lg text-white text-center focus:outline-none focus:border-blue-500"
        />
        <button
          @click="handleAnswer(openAnswer)"
          :disabled="!openAnswer.trim()"
          class="w-full mt-4 px-8 py-4 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-700 disabled:cursor-not-allowed rounded-lg text-xl font-semibold transition-colors"
        >
          提交答案
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Question } from '@/shared/types'

const props = defineProps<{
  playerName: string
}>()

const currentQuestion = ref<Question | null>(null)
const answerSubmitted = ref(false)
const openAnswer = ref('')
const connected = ref(false)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

let ws: WebSocketClient | null = null
const deviceType = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(navigator.userAgent)
  ? 'mobile'
  : 'desktop'

onMounted(() => {
  const storedName = sessionStorage.getItem('playerName') || props.playerName

  ws = new WebSocketClient(wsUrl.value)

  ws.on('quiz_question', (payload: Question) => {
    currentQuestion.value = payload
    answerSubmitted.value = false
  })

  ws.on('quiz_next', () => {
    currentQuestion.value = null
    answerSubmitted.value = false
  })

  ws.on('quiz_answer_update', () => {
    answerSubmitted.value = true
  })

  ws.on('open', () => {
    connected.value = true
    ws?.send('join', {
      name: storedName,
      deviceType: deviceType
    })
  })

  ws.connect()
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
})

const handleAnswer = (answer: string) => {
  if (!answer.trim() || !ws) return

  ws.send('quiz_answer', { answer })
  answerSubmitted.value = true
  openAnswer.value = ''
}
</script>
