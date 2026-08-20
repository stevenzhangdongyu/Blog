
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ConfettiEffect from './components/ConfettiEffect.vue'

const backgroundReady = ref(false)
const backgroundImage = ref('')

onMounted(() => {
  const image = new Image()
  image.onload = () => {
    backgroundImage.value = `linear-gradient(rgba(255, 250, 245, 0.78), rgba(255, 250, 245, 0.88)), url(${image.src})`
    backgroundReady.value = true
  }
  image.onerror = () => {
    // Text remains usable if the optional decorative image is unavailable.
    backgroundReady.value = true
  }
  image.src = '/blog-background.jpg'
})
</script>

<template>
  <div class="page-background" :class="{ 'is-ready': backgroundReady }" :style="{ backgroundImage }" aria-hidden="true" />
  <div v-if="!backgroundReady" class="background-loader" role="status" aria-live="polite">
    <span class="loader-ring" aria-hidden="true"></span>
    <span>背景加载中</span>
  </div>
  <ConfettiEffect />
  <RouterView />
</template>
