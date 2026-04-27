<template>
  <div class="w-full h-full bg-slate-900 text-white flex items-center justify-center p-8">
    <div class="text-6xl mb-8">🎯</div>
    <div class="text-3xl font-bold mb-4">输入你的名字</div>
    <div class="text-slate-400 mb-8">将在抢答时显示</div>

    <input
      v-model="name"
      @keyup.enter="handleSubmit"
      type="text"
      placeholder="例如：张三"
      class="w-full max-w-md px-6 py-4 text-xl bg-slate-800 border-2 border-slate-700 rounded-lg text-white text-center focus:outline-none focus:border-blue-500"
      maxlength="20"
    />

    <button
      @click="handleSubmit"
      :disabled="!name.trim()"
      class="w-full max-w-md mt-6 px-8 py-4 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-700 disabled:cursor-not-allowed rounded-lg text-xl font-semibold transition-colors"
    >
      进入抢答
    </button>

    <div v-if="error" class="mt-4 text-red-400">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const name = ref('')
const error = ref('')

const handleSubmit = () => {
  const trimmed = name.value.trim()

  if (!trimmed) {
    error.value = '请输入名字'
    return
  }

  if (trimmed.length < 2) {
    error.value = '名字至少需要2个字符'
    return
  }

  if (trimmed.length > 20) {
    error.value = '名字最多20个字符'
    return
  }

  sessionStorage.setItem('playerName', trimmed)
  window.location.href = '/player?step=buzzer'
}
</script>
