
<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import ConfettiEffect from './components/ConfettiEffect.vue'

const backgroundReady = ref(false)
const backgroundLayers = ref(['', ''])
const visibleLayer = ref(0)
const wallpapers = [
  '/blog-background.jpg',
  '/wallpaper-school.jpg',
  '/wallpaper-yukino-wide.jpg',
  '/wallpaper-anime.jpg',
]
let wallpaperIndex = 0
let wallpaperTimer = 0

function setBackground(source: string) {
  const image = new Image()
  image.onload = () => {
    const nextLayer = visibleLayer.value === 0 ? 1 : 0
    backgroundLayers.value[nextLayer] = `linear-gradient(rgba(255, 250, 245, 0.78), rgba(255, 250, 245, 0.88)), url(${image.src})`
    visibleLayer.value = nextLayer
    backgroundReady.value = true
  }
  image.onerror = () => {
    backgroundReady.value = true
  }
  image.src = source
}

function preload(source: string) {
  const image = new Image()
  image.src = source
}

onMounted(() => {
  setBackground(wallpapers[wallpaperIndex])
  preload(wallpapers[(wallpaperIndex + 1) % wallpapers.length])
  wallpaperTimer = window.setInterval(() => {
    wallpaperIndex = (wallpaperIndex + 1) % wallpapers.length
    setBackground(wallpapers[wallpaperIndex])
    preload(wallpapers[(wallpaperIndex + 1) % wallpapers.length])
  }, 15000)
})

onBeforeUnmount(() => window.clearInterval(wallpaperTimer))
</script>

<template>
  <div
    v-for="(backgroundImage, index) in backgroundLayers"
    :key="index"
    class="page-background"
    :class="{ 'is-ready': backgroundReady && visibleLayer === index }"
    :style="{ backgroundImage }"
    aria-hidden="true"
  />
  <div v-if="!backgroundReady" class="background-loader" role="status" aria-live="polite">
    <span class="loader-ring" aria-hidden="true"></span>
    <span>背景加载中</span>
  </div>
  <ConfettiEffect />
  <RouterView />
</template>
