<template>
  <div class="quiz-display bg-slate-800 rounded-2xl p-8 mb-6">
    <div v-if="!currentQuestion" class="text-center py-12">
      <div class="text-slate-400 text-lg">
        {{ allQuestionsDone ? '✅ 所有题目已完成！' : '点击"开始答题"按钮开始' }}
      </div>
    </div>

    <div v-else>
      <div class="text-2xl font-bold mb-8 text-center">
        {{ currentQuestion.question }}
      </div>

      <div
        v-if="currentQuestion.options && currentQuestion.options.length > 0"
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Question } from '@/shared/types'

const emit = defineEmits(['quizQuestion'])

const props = defineProps<{
  ws: any
}>()

const currentQuestion = ref<Question | null>(null)
const allQuestionsDone = ref(false)

if (props.ws) {
  props.ws.on('quiz_question', (payload: Question) => {
    console.log('quiz_question', payload)
    currentQuestion.value = payload
    allQuestionsDone.value = false
    emit('quizQuestion', payload)
  })

  // Listen for when no more questions are available
  props.ws.on('quiz_no_questions', () => {
    currentQuestion.value = null
    allQuestionsDone.value = true
  })
}
</script>
