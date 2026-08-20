
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ConfettiEffect from './components/ConfettiEffect.vue'

const backgroundReady = ref(false)

onMounted(() => {
  const background = new Image()
  background.onload = () => {
    backgroundReady.value = true
  }
  background.src = '/blog-background.jpg'
})
</script>

<template>
  <div class="page-background" :class="{ 'is-ready': backgroundReady }" aria-hidden="true" />
  <ConfettiEffect />
  <RouterView />
</template>
