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
interface Comment {
  id: number
  authorName: string
  content: string
  createdAt: string
}
const htmlContent = computed(() => renderMarkdown(article.value?.content || ''))

const route = useRoute()
const article = ref<Article | null>(null)
const errorMessage = ref('')
const comments = ref<Comment[]>([])
const commentAuthor = ref('')
const commentContent = ref('')
const commentError = ref('')
const commentSubmitting = ref(false)

async function fetchArticle(slug: string) {
  article.value = null
  comments.value = []
  errorMessage.value = ''

  try {
    const response = await fetch(`/api/articles/${encodeURIComponent(slug)}`)

    if (response.status === 404) {
      errorMessage.value = '文章不存在或尚未发布。'
      return
    }

    if (!response.ok) throw new Error(`HTTP ${response.status}`)
    article.value = await response.json()
    await fetchComments(slug)
  
  } catch (error) {
    console.error('Error fetching article:', error)
    errorMessage.value = '暂时无法加载文章。'
  }
}

async function fetchComments(slug: string) {
  const response = await fetch(`/api/articles/${encodeURIComponent(slug)}/comments`)
  if (!response.ok) return
  const result = await response.json()
  comments.value = result.comments || []
}

async function submitComment() {
  commentError.value = ''
  commentSubmitting.value = true
  try {
    const response = await fetch(`/api/articles/${encodeURIComponent(String(route.params.slug))}/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ authorName: commentAuthor.value, content: commentContent.value, website: '' }),
    })
    const result = await response.json().catch(() => ({}))
    if (!response.ok) {
      commentError.value = result.error || '发表评论失败'
      return
    }
    comments.value.push(result)
    commentAuthor.value = ''
    commentContent.value = ''
    window.dispatchEvent(new Event('comment-celebration'))
  } catch {
    commentError.value = '暂时无法发表评论，请稍后再试。'
  } finally {
    commentSubmitting.value = false
  }
}

onMounted(() => fetchArticle(String(route.params.slug)))
watch(() => route.params.slug, (slug) => fetchArticle(String(slug)))
</script>

<template>
  <main class="shell">
    <RouterLink to="/">返回文章列表</RouterLink>

    <p v-if="!article && !errorMessage" class="intro inline-loading"><span class="loader-ring" aria-hidden="true"></span>正在加载文章...</p>
    <p v-else-if="errorMessage" class="intro">{{ errorMessage }}</p>

    <article v-else-if="article">
      <p class="eyebrow">ARTICLE</p>
      <div class="markdown-body" v-html="htmlContent"></div>
      <section class="comments" aria-labelledby="comments-title">
        <p id="comments-title" class="label">COMMENTS</p>
        <div v-if="comments.length" class="comment-list">
          <article v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-meta"><strong>{{ comment.authorName }}</strong><time>{{ new Date(comment.createdAt).toLocaleString('zh-CN') }}</time></div>
            <p>{{ comment.content }}</p>
          </article>
        </div>
        <p v-else class="comment-empty">还没有评论，留下第一条吧。</p>
        <form class="comment-form" @submit.prevent="submitComment">
          <label>昵称 <input v-model="commentAuthor" maxlength="50" required placeholder="怎么称呼你？" /></label>
          <label>评论 <textarea v-model="commentContent" maxlength="1000" rows="4" required placeholder="写下你的想法" /></label>
          <input class="honeypot" tabindex="-1" autocomplete="off" aria-hidden="true" />
          <p v-if="commentError" class="form-error">{{ commentError }}</p>
          <button :disabled="commentSubmitting">{{ commentSubmitting ? '提交中...' : '发表评论' }}</button>
        </form>
      </section>
    </article>
  </main>
</template>
