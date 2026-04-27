<template>
  <div class="flex gap-4 justify-center mb-6">
    <button
      @click="$emit('start')"
      :disabled="gameState !== 'waiting'"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="gameState === 'waiting' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-gray-600 cursor-not-allowed opacity-50'"
    >
      开始抢答
    </button>

    <button
      @click="$emit('reset')"
      :disabled="gameState === 'waiting'"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="gameState !== 'waiting' ? 'bg-gray-500 hover:bg-gray-600' : 'bg-gray-700 cursor-not-allowed opacity-50'"
    >
      下一题
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  gameState: 'waiting' | 'ready' | 'locked'
}>()

defineEmits<{
  start: []
  reset: []
}>()
</script>
