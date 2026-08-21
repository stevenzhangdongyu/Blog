<script setup lang="ts">
import { onMounted, ref } from 'vue'

interface Article {
  id: number
  title: string
  slug: string
  summary: string
  status: string
}

interface ArticlesResponse {
  articles: Article[]
}

const apiStatus = ref('正在检查...')
const articles = ref<Article[]>([])
const homeIntro = ref('只有一种成功，那就是按照自己的意愿过完一生。')

async function fetchHealth() {
  try {
    const response = await fetch('/api/health')
    if (!response.ok) throw new Error(`HTTP ${response.status}`)

    const result: { status: string } = await response.json()
    apiStatus.value = result.status
  } catch (error) {
    console.error('Error fetching health:', error)
    apiStatus.value = '连接失败'
  }
}

async function fetchArticles() {
  try {
    const response = await fetch('/api/articles')
    if (!response.ok) throw new Error(`HTTP ${response.status}`)

    const result: ArticlesResponse = await response.json()
    articles.value = result.articles
  } catch (error) {
    console.error('Error fetching articles:', error)
  }
}

async function fetchSiteSettings() {
  try {
    const response = await fetch('/api/site-settings')
    if (!response.ok) throw new Error(`HTTP ${response.status}`)

    const result: { homeIntro: string } = await response.json()
    if (result.homeIntro) homeIntro.value = result.homeIntro
  } catch (error) {
    console.error('Error fetching site settings:', error)
  }
}

onMounted(() => {
  fetchHealth()
  fetchArticles()
  fetchSiteSettings()
  const preloadArticleView = () => void import('./ArticleView.vue')
  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(preloadArticleView, { timeout: 3000 })
  } else {
    globalThis.setTimeout(preloadArticleView, 1000)
  }
})
</script>

<template>
  <main class="shell">
    <p class="eyebrow">PERSONAL BLOG LAB</p>
    <h1 class="home-title">写出我心</h1>
    <p class="intro">{{ homeIntro }}</p>

    <section aria-labelledby="articles-title">
      <p id="articles-title" class="label">ARTICLES</p>
      <ul class="article-list">
        <li v-for="article in articles" :key="article.id">
          <RouterLink class="article-link" :to="`/articles/${article.slug}`">
            <span class="article-title">{{ article.title }}</span>
            <span class="article-summary">{{ article.summary }}</span>
          </RouterLink>
        </li>
      </ul>
    </section>
  </main>
</template>
