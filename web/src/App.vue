
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ConfettiEffect from './components/ConfettiEffect.vue'

const backgroundReady = ref(false)
const backgroundProgress = ref(0)
const backgroundHasProgress = ref(false)
const backgroundImage = ref('')

onMounted(async () => {
  try {
    const response = await fetch('/blog-background.jpg')
    if (!response.ok || !response.body) throw new Error('background unavailable')
    const total = Number(response.headers.get('content-length'))
    backgroundHasProgress.value = Number.isFinite(total) && total > 0
    const reader = response.body.getReader()
    const chunks: Uint8Array[] = []
    let loaded = 0
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      if (value) {
        chunks.push(value)
        loaded += value.length
        if (backgroundHasProgress.value) backgroundProgress.value = Math.min(99, Math.round((loaded / total) * 100))
      }
    }
    const blob = new Blob(chunks, { type: response.headers.get('content-type') || 'image/jpeg' })
    const image = new Image()
    image.onload = () => {
      backgroundProgress.value = 100
      backgroundImage.value = `linear-gradient(rgba(255, 250, 245, 0.78), rgba(255, 250, 245, 0.88)), url(${image.src})`
      backgroundReady.value = true
    }
    image.src = URL.createObjectURL(blob)
  } catch {
    backgroundProgress.value = 0
  }
})
</script>

<template>
  <div class="page-background" :class="{ 'is-ready': backgroundReady }" :style="{ backgroundImage }" aria-hidden="true" />
  <div v-if="!backgroundReady" class="background-loader" role="status" aria-live="polite">
    <span class="loader-ring" aria-hidden="true"></span>
    <span>{{ backgroundHasProgress && backgroundProgress ? `背景加载中 ${backgroundProgress}%` : '背景加载中' }}</span>
  </div>
  <ConfettiEffect />
  <RouterView />
</template>
