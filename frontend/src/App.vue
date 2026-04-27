<template>
  <HostApp v-if="isHost" />
  <PlayerApp v-else />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import HostApp from './host/HostApp.vue'
import PlayerApp from './player/PlayerApp.vue'

const isHost = ref(false)

onMounted(() => {
  const userAgent = navigator.userAgent.toLowerCase()
  const isMobile = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(userAgent)
  const path = window.location.pathname

  isHost.value = !isMobile && (path === '/' || path === '/host')
})
</script>
