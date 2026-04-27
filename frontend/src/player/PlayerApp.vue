<template>
  <NameInput v-if="currentStep === 'input'" />
  <QuizAnswer v-else-if="currentStep === 'game' && gameMode === 'quiz'" :player-name="playerName" />
  <BuzzerButton v-else-if="currentStep === 'game' && gameMode === 'buzzer'" :player-name="playerName" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import NameInput from './components/NameInput.vue'
import BuzzerButton from './components/BuzzerButton.vue'
import QuizAnswer from './components/QuizAnswer.vue'
import { WebSocketClient } from '@/shared/websocket'

const currentStep = ref('input')
const playerName = ref('')
const gameMode = ref<'buzzer' | 'quiz'>('buzzer')

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

let ws: WebSocketClient | null = null

onMounted(() => {
  const storedName = sessionStorage.getItem('playerName')
  if (storedName) {
    playerName.value = storedName
    currentStep.value = 'game'

    // Create WebSocket connection for mode detection
    ws = new WebSocketClient(wsUrl.value)

    ws.on('quiz_question', () => {
      gameMode.value = 'quiz'
    })

    ws.on('state_changed', () => {
      // Reset to buzzer mode when game state changes (e.g., reset)
      if (gameMode.value === 'quiz') {
        gameMode.value = 'buzzer'
      }
    })

    ws.connect()
  }
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
})
</script>
