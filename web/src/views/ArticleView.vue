<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { renderMarkdown } from '../utils/markdown'

interface Article {
  id: number
  title: string
  slug: string
  summary: string
  content: string
  status: string
}
const htmlContent = computed(() => renderMarkdown(article.value?.content || ''))

const route = useRoute()
const article = ref<Article | null>(null)
const errorMessage = ref('')

async function fetchArticle(slug: string) {
  article.value = null
  errorMessage.value = ''

  try {
    const response = await fetch(`/api/articles/${encodeURIComponent(slug)}`)

    if (response.status === 404) {
      errorMessage.value = '文章不存在或尚未发布。'
      return
    }

    if (!response.ok) throw new Error(`HTTP ${response.status}`)
    article.value = await response.json()
  
  } catch (error) {
    console.error('Error fetching article:', error)
    errorMessage.value = '暂时无法加载文章。'
  }
}

onMounted(() => fetchArticle(String(route.params.slug)))
watch(() => route.params.slug, (slug) => fetchArticle(String(slug)))
</script>

<template>
  <main class="shell">
    <RouterLink to="/">返回文章列表</RouterLink>

    <p v-if="!article && !errorMessage" class="intro">正在加载文章...</p>
    <p v-else-if="errorMessage" class="intro">{{ errorMessage }}</p>

    <article v-else-if="article">
      <p class="eyebrow">ARTICLE</p>
      <div class="markdown-body" v-html="htmlContent"></div>
    </article>
  </main>
</template>
