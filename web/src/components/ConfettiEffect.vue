<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'

type Particle = {
  x: number
  y: number
  width: number
  height: number
  opacity: number
  velocityX: number
  velocityY: number
  rotation: number
  rotationVelocity: number
  color: string
}

const canvas = ref<HTMLCanvasElement | null>(null)
const particles: Particle[] = []
const colors = ['#e63946', '#f4a261', '#e9c46a', '#2a9d8f', '#457b9d', '#7b2cbf', '#ff4d8d']
let context: CanvasRenderingContext2D | null = null
let frame = 0
let reduceMotion = false

function resize() {
  if (!canvas.value) return
  const scale = window.devicePixelRatio || 1
  canvas.value.width = window.innerWidth * scale
  canvas.value.height = window.innerHeight * scale
  canvas.value.style.width = `${window.innerWidth}px`
  canvas.value.style.height = `${window.innerHeight}px`
  context = canvas.value.getContext('2d')
  context?.setTransform(scale, 0, 0, scale, 0, 0)
}

function animate() {
  if (!context) return
  context.clearRect(0, 0, window.innerWidth, window.innerHeight)
  for (let index = particles.length - 1; index >= 0; index -= 1) {
    const particle = particles[index]
    particle.x += particle.velocityX
    particle.y += particle.velocityY
    particle.velocityY += 0.11
    particle.rotation += particle.rotationVelocity
    particle.opacity -= 0.012
    if (particle.opacity <= 0) {
      particles.splice(index, 1)
      continue
    }
    context.save()
    context.translate(particle.x, particle.y)
    context.rotate(particle.rotation)
    context.globalAlpha = particle.opacity
    context.fillStyle = particle.color
    context.fillRect(-particle.width / 2, -particle.height / 2, particle.width, particle.height)
    context.restore()
  }
  frame = particles.length ? requestAnimationFrame(animate) : 0
}

function burst(event: PointerEvent) {
  if (reduceMotion) return
  for (let count = 0; count < 72; count += 1) {
    const angle = Math.random() * Math.PI * 2
    const speed = 2 + Math.random() * 8
    particles.push({
      x: event.clientX,
      y: event.clientY,
      width: 5 + Math.random() * 8,
      height: 3 + Math.random() * 5,
      opacity: 1,
      velocityX: Math.cos(angle) * speed,
      velocityY: Math.sin(angle) * speed - 3,
      rotation: Math.random() * Math.PI,
      rotationVelocity: (Math.random() - 0.5) * 0.5,
      color: colors[Math.floor(Math.random() * colors.length)],
    })
  }
  if (!frame) frame = requestAnimationFrame(animate)
}

function launchCelebrationWave(amount: number) {
  for (let count = 0; count < amount; count += 1) {
    const edgeBias = Math.random()
    const x = edgeBias < 0.3
      ? Math.random() * window.innerWidth * 0.25
      : edgeBias > 0.7
        ? window.innerWidth * 0.75 + Math.random() * window.innerWidth * 0.25
        : Math.random() * window.innerWidth
    particles.push({
      x,
      y: window.innerHeight + Math.random() * 24,
      width: 5 + Math.random() * 10,
      height: 3 + Math.random() * 7,
      opacity: 1,
      velocityX: (window.innerWidth / 2 - x) * 0.003 + (Math.random() - 0.5) * 5,
      velocityY: -(8 + Math.random() * 11),
      rotation: Math.random() * Math.PI,
      rotationVelocity: (Math.random() - 0.5) * 0.55,
      color: colors[Math.floor(Math.random() * colors.length)],
    })
  }
  if (!frame) frame = requestAnimationFrame(animate)
}

function celebrate() {
  if (reduceMotion) return
  launchCelebrationWave(150)
  window.setTimeout(() => launchCelebrationWave(100), 180)
}

onMounted(() => {
  reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  resize()
  window.addEventListener('resize', resize)
  window.addEventListener('pointerdown', burst)
  window.addEventListener('comment-celebration', celebrate)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(frame)
  window.removeEventListener('resize', resize)
  window.removeEventListener('pointerdown', burst)
  window.removeEventListener('comment-celebration', celebrate)
})
</script>

<template>
  <canvas ref="canvas" class="particle-canvas" aria-hidden="true" />
</template>
