<template>
  <NameInput v-if="currentStep === 'input'" />
  <QuizAnswer v-else-if="currentStep === 'game'" :player-name="playerName" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NameInput from './components/NameInput.vue'
import QuizAnswer from './components/QuizAnswer.vue'

const currentStep = ref('input')
const playerName = ref('')

onMounted(() => {
  const storedName = sessionStorage.getItem('playerName')
  if (storedName) {
    playerName.value = storedName
    currentStep.value = 'game'
  }
})
</script>
