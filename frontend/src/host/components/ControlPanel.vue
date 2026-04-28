<template>
  <div class="flex gap-4 justify-center mb-6">
    <button
      v-if="gameMode === 'buzzer'"
      @click="$emit('start')"
      :disabled="gameState !== 'waiting'"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="gameState === 'waiting' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-gray-600 cursor-not-allowed opacity-50'"
    >
      开始抢答
    </button>

    <button
      v-if="gameMode === 'quiz'"
      @click="$emit('start')"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors bg-green-500 hover:bg-green-600"
    >
      开始答题
    </button>

    <button
      @click="$emit('next')"
      :disabled="!canNext"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="canNext ? 'bg-gray-500 hover:bg-gray-600' : 'bg-gray-700 cursor-not-allowed opacity-50'"
    >
      {{ gameMode === 'quiz' ? '下一题' : '重置' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  gameState: 'waiting' | 'ready' | 'locked'
  gameMode: 'buzzer' | 'quiz'
}>()

const canNext = computed(() => {
  if (props.gameMode === 'buzzer') {
    return props.gameState !== 'waiting'
  }
  // Quiz mode: always can go to next question
  return true
})

defineEmits<{
  start: []
  next: []
}>()
</script>
