


<script setup lang="ts">
const nextMilestone = '读取并展示第一篇已发布文章'
import { ref, onMounted } from 'vue'
const apiStatus = ref("OK")
const articles = ref<Article[]>([])
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
const fetchData = async () => {
  try {
    const response = await fetch('/api/health')
    
    // fetch 不会自动抛出 HTTP 错误状态（如 404、500），需要手动判断 ok
    if (!response.ok) {
      throw new Error(`HTTP 错误！状态码: ${response.status}`)
    }
    
    const result = await response.json()
    apiStatus.value = result.status
    // apiStatus.value = result // 将解析后的数据赋给响应式变量
  } catch (err) {
    apiStatus.value = "Error"
  } finally {
    console.log("finished")
  }
}

const fetchArticles = async () => {
  try {
    const response = await fetch('/api/articles')
    
    if (!response.ok) {
      throw new Error(`HTTP 错误！状态码: ${response.status}`)
    }
    
    const result: ArticlesResponse = await response.json()
    articles.value = result.articles
  } catch (err) {
    console.error("Error fetching articles:", err)
  }
}

onMounted(() => {
  fetchData()
  fetchArticles()

})
</script>

<template>
  <main class="shell">
    <p class="eyebrow">PERSONAL BLOG LAB</p>
    <h1>从这里开始写。</h1>
    <p class="intro">
      页面已经启动，但博客功能仍是空白。下一步由你亲手打通数据库、Go API 和 Vue 页面。
    </p>
    <p>API 状态：{{ apiStatus }}</p>
    <section aria-labelledby="next-title">
      <p id="next-title" class="label">NEXT MILESTONE</p>
      <p class="milestone">{{ nextMilestone }}</p>
    </section>

    <section aria-labelledby="articles-title">
      <p id="articles-title" class="label">ARTICLES</p>
      <ul>
        <li v-for="article in articles" :key="article.id">
          <h2>{{ article.title }}</h2>
          <p>{{ article.summary }}</p>
          <p>Status: {{ article.status }}</p>
        </li>
      </ul>
    </section>
  </main>
</template>

