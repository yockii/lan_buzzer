<template>
  <div class="quiz-display bg-slate-800 rounded-2xl p-8 mb-6">
    <div v-if="!currentQuestion" class="text-center py-12">
      <button
        @click="handleStart"
        class="px-8 py-4 bg-blue-500 hover:bg-blue-600 rounded-lg text-xl font-semibold transition-colors"
      >
        开始答题
      </button>
    </div>

    <div v-else>
      <div class="text-2xl font-bold mb-8 text-center">
        {{ currentQuestion.question }}
      </div>

      <div
        v-if="currentQuestion.options.length > 0"
        class="grid gap-4 max-w-2xl mx-auto"
        :class="currentQuestion.options.length === 2 ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <button
          v-for="(option, index) in currentQuestion.options"
          :key="index"
          class="px-6 py-4 bg-slate-700 hover:bg-slate-600 rounded-lg text-lg text-left transition-colors"
        >
          {{ option }}
        </button>
      </div>

      <div v-else class="text-center text-slate-400">
        (问答题，选手输入答案)
      </div>

      <div class="mt-8 text-center">
        <button
          @click="handleNext"
          class="px-8 py-4 bg-green-500 hover:bg-green-600 rounded-lg text-xl font-semibold transition-colors"
        >
          下一题
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Question } from '@/shared/types'

const props = defineProps<{
  ws: WebSocketClient | null
}>()

const currentQuestion = ref<Question | null>(null)

const handleStart = () => {
  if (props.ws) {
    props.ws.send('quiz_start', {})
  }
}

const handleNext = () => {
  if (props.ws) {
    props.ws.send('quiz_next', {})
  }
}

if (props.ws) {
  props.ws.on('quiz_question', (payload: Question) => {
    currentQuestion.value = payload
  })
}
</script>
