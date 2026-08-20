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
  quotedText: string
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
const articleBody = ref<HTMLElement | null>(null)
const quoteMenu = ref({ visible: false, x: 0, y: 0, text: '' })
const highlightedQuote = ref('')

function captureSelection() {
  const selection = window.getSelection()
  const text = selection?.toString().trim() || ''
  if (!selection || !text || !articleBody.value || !articleBody.value.contains(selection.anchorNode)) {
    quoteMenu.value.visible = false
    return
  }
  const rect = selection.getRangeAt(0).getBoundingClientRect()
  quoteMenu.value = { visible: true, x: Math.max(12, rect.left + rect.width / 2 - 52), y: Math.max(12, rect.top - 44), text: text.slice(0, 300) }
}

function useSelectedQuote() {
  commentContent.value = commentContent.value || `关于“${quoteMenu.value.text}”：`
  quoteMenu.value.visible = false
  document.querySelector('.comment-form')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  window.setTimeout(() => document.querySelector<HTMLTextAreaElement>('.comment-form textarea')?.focus(), 350)
}

function jumpToQuote(text: string) {
  if (!text || !articleBody.value) return
  if (!articleBody.value.textContent?.includes(text)) {
    commentError.value = '原文已发生变化，暂时找不到这段引用。'
    return
  }
  highlightedQuote.value = text
  articleBody.value.scrollIntoView({ behavior: 'smooth', block: 'center' })
  window.setTimeout(() => { highlightedQuote.value = '' }, 2200)
}

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
      body: JSON.stringify({ authorName: commentAuthor.value, content: commentContent.value, quotedText: quoteMenu.value.text, website: '' }),
    })
    const result = await response.json().catch(() => ({}))
    if (!response.ok) {
      commentError.value = result.error || '发表评论失败'
      return
    }
    comments.value.push(result)
    commentAuthor.value = ''
    commentContent.value = ''
    quoteMenu.value.text = ''
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
      <div ref="articleBody" class="markdown-body" :class="{ 'quote-highlight-active': highlightedQuote }" @mouseup="captureSelection" @touchend="captureSelection" v-html="htmlContent"></div>
      <button v-if="quoteMenu.visible" class="quote-menu" :style="{ left: `${quoteMenu.x}px`, top: `${quoteMenu.y}px` }" type="button" @click="useSelectedQuote">引用评论</button>
      <section class="comments" aria-labelledby="comments-title">
        <p id="comments-title" class="label">COMMENTS</p>
        <div v-if="comments.length" class="comment-list">
          <article v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-meta"><strong>{{ comment.authorName }}</strong><time>{{ new Date(comment.createdAt).toLocaleString('zh-CN') }}</time></div>
            <button v-if="comment.quotedText" class="comment-quote" type="button" @click="jumpToQuote(comment.quotedText)">“{{ comment.quotedText }}”</button>
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
