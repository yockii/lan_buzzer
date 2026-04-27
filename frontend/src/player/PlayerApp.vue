<template>
  <NameInput v-if="currentStep === 'input'" />
  <BuzzerButton v-else :player-name="playerName" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NameInput from './components/NameInput.vue'
import BuzzerButton from './components/BuzzerButton.vue'

const currentStep = ref('input')
const playerName = ref('')

onMounted(() => {
  const storedName = sessionStorage.getItem('playerName')
  if (storedName) {
    playerName.value = storedName
    currentStep.value = 'buzzer'
  }
})
</script>
