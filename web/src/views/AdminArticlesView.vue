<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { renderMarkdown } from '../utils/markdown'

interface Article {
  id: number
  title: string
  slug: string
  summary: string
  content: string
  status: string
}
interface Quote { id: number; content: string }

const router = useRouter()
const articles = ref<Article[]>([])
const editing = ref<Article | null>(null)
const errorMessage = ref('')
const homeIntro = ref('')
const settingsMessage = ref('')
const settingsError = ref('')
const settingsSaving = ref(false)
const quotes = ref<Quote[]>([])
const newQuote = ref('')
const quoteMessage = ref('')
const editor = ref<HTMLTextAreaElement | null>(null)
const previewHtml = computed(() => renderMarkdown(editing.value?.content ?? ''))

async function loadArticles() {
  const response = await fetch('/api/admin/articles', { credentials: 'include' })
  if (response.status === 401) return router.push('/admin/login')
  const result = await response.json()
  articles.value = result.articles
}

async function loadSettings() {
  const response = await fetch('/api/site-settings')
  if (!response.ok) return
  const result: { homeIntro: string } = await response.json()
  homeIntro.value = result.homeIntro
}

async function loadQuotes() {
  const response = await fetch('/api/admin/quotes', { credentials: 'include' })
  if (response.status === 401) return router.push('/admin/login')
  if (!response.ok) return
  quotes.value = (await response.json()).quotes
}

async function addQuote() {
  const content = newQuote.value.trim()
  if (!content) return
  const response = await fetch('/api/admin/quotes', { method: 'POST', headers: { 'Content-Type': 'application/json' }, credentials: 'include', body: JSON.stringify({ content }) })
  if (response.ok) { newQuote.value = ''; quoteMessage.value = '名言已添加'; await loadQuotes() }
}

async function editQuote(quote: Quote) {
  const content = window.prompt('修改名言内容', quote.content)?.trim()
  if (!content || content === quote.content) return
  const response = await fetch(`/api/admin/quotes/${quote.id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, credentials: 'include', body: JSON.stringify({ content }) })
  if (response.ok) { quoteMessage.value = '名言已更新'; await loadQuotes() }
}

async function removeQuote(id: number) {
  if (!window.confirm('确定删除这句名言吗？')) return
  const response = await fetch(`/api/admin/quotes/${id}`, { method: 'DELETE', credentials: 'include' })
  if (response.ok) { quoteMessage.value = '名言已删除'; await loadQuotes() }
}

async function saveSettings() {
  settingsMessage.value = ''
  settingsError.value = ''
  settingsSaving.value = true
  try {
    const response = await fetch('/api/admin/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ homeIntro: homeIntro.value }),
    })
    const result = await response.json().catch(() => ({}))
    if (!response.ok) {
      settingsError.value = result.error || '首页文案保存失败'
      return
    }
    homeIntro.value = result.homeIntro
    settingsMessage.value = '首页文案已保存'
  } catch {
    settingsError.value = '首页文案保存失败，请稍后再试'
  } finally {
    settingsSaving.value = false
  }
}

function newArticle() {
  editing.value = { id: 0, title: '', slug: '', summary: '', content: '', status: 'draft' }
}

function createSlug(title: string) {
  return title
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

function updateTitle() {
  if (!editing.value) return
  if (!editing.value.slug) editing.value.slug = createSlug(editing.value.title)
}

async function insertMarkdown(before: string, after = '', placeholder = '文本') {
  if (!editing.value) return
  const input = editor.value
  const start = input?.selectionStart ?? editing.value.content.length
  const end = input?.selectionEnd ?? start
  const selected = editing.value.content.slice(start, end) || placeholder
  editing.value.content = `${editing.value.content.slice(0, start)}${before}${selected}${after}${editing.value.content.slice(end)}`
  await nextTick()
  input?.focus()
  input?.setSelectionRange(start + before.length, start + before.length + selected.length)
}

function insertCodeBlock() {
	return insertMarkdown('```go\n', '\n```', 'fmt.Println("Hello")')
}

function handleEditorKeydown(event: KeyboardEvent) {
  if (event.key !== 'Tab' || !editing.value) return
  event.preventDefault()
  const input = event.target as HTMLTextAreaElement
  const start = input.selectionStart
  editing.value.content = `${editing.value.content.slice(0, start)}  ${editing.value.content.slice(start)}`
  nextTick(() => input.setSelectionRange(start + 2, start + 2))
}

async function save() {
  if (!editing.value) return
  errorMessage.value = ''
  const isNew = editing.value.id === 0
  const response = await fetch(isNew ? '/api/admin/articles' : `/api/admin/articles/${editing.value.id}`, {
    method: isNew ? 'POST' : 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(editing.value),
  })
  if (!response.ok) {
    const result = await response.json().catch(() => ({}))
    errorMessage.value = result.error || '保存失败'
    return
  }
  editing.value = null
  await loadArticles()
}

async function remove(id: number) {
  if (!confirm('确定删除这篇文章吗？')) return
  await fetch(`/api/admin/articles/${id}`, { method: 'DELETE', credentials: 'include' })
  await loadArticles()
}

async function logout() {
  await fetch('/api/admin/logout', { method: 'POST', credentials: 'include' })
  await router.push('/admin/login')
}

onMounted(() => {
  loadArticles()
  loadSettings()
  loadQuotes()
})
</script>

<template>
  <main class="admin-shell">
    <header class="admin-header">
      <div><p class="eyebrow">ADMIN STUDIO</p><h1>文章管理</h1></div>
      <div class="admin-actions"><button @click="newArticle">新建文章</button><button class="secondary" @click="logout">退出</button></div>
    </header>
    <section v-if="editing" class="editor">
      <label>标题 <input v-model="editing.title" @input="updateTitle" /></label>
      <label>Slug <input v-model="editing.slug" placeholder="留空则按标题自动生成" /></label>
      <label>摘要 <textarea v-model="editing.summary" rows="3" /></label>
      <section class="markdown-workbench" aria-label="Markdown 编辑器">
        <div class="markdown-pane">
          <div class="editor-toolbar" aria-label="Markdown 工具栏">
            <button type="button" title="一级标题" @click="insertMarkdown('# ', '', '标题')">H1</button>
            <button type="button" title="二级标题" @click="insertMarkdown('## ', '', '标题')">H2</button>
            <button type="button" title="加粗" @click="insertMarkdown('**', '**')">B</button>
            <button type="button" title="斜体" @click="insertMarkdown('*', '*')">I</button>
            <button type="button" title="链接" @click="insertMarkdown('[', '](https://)')">链接</button>
            <button type="button" title="行内代码" @click="insertMarkdown('`', '`', 'code')">&lt;/&gt;</button>
	            <button type="button" title="代码块" @click="insertCodeBlock">代码块</button>
            <button type="button" title="引用" @click="insertMarkdown('> ', '', '引用内容')">引用</button>
          </div>
          <textarea ref="editor" v-model="editing.content" rows="22" spellcheck="false" @keydown="handleEditorKeydown" />
        </div>
        <div class="markdown-pane preview-pane">
          <p class="pane-label">实时预览</p>
          <div class="markdown-body" v-html="previewHtml" />
        </div>
      </section>
      <label>状态 <select v-model="editing.status"><option value="draft">草稿</option><option value="published">发布</option></select></label>
      <p v-if="errorMessage" class="form-error">{{ errorMessage }}</p>
      <button @click="save">保存</button>
      <button class="secondary" @click="editing = null">取消</button>
    </section>
    <template v-else>
      <section class="site-settings" aria-labelledby="site-settings-title">
        <div class="site-settings-heading">
          <div><p class="label">HOME PAGE</p><h2 id="site-settings-title">首页文案</h2></div>
          <button :disabled="settingsSaving || !homeIntro.trim()" @click="saveSettings">{{ settingsSaving ? '保存中...' : '保存文案' }}</button>
        </div>
        <textarea v-model="homeIntro" maxlength="500" rows="3" aria-label="首页文案" placeholder="输入首页展示的文案" />
        <div class="settings-feedback" aria-live="polite">
          <p v-if="settingsError" class="form-error">{{ settingsError }}</p>
          <p v-else-if="settingsMessage" class="form-success">{{ settingsMessage }}</p>
          <span>{{ homeIntro.length }}/500</span>
        </div>
      </section>
      <section class="site-settings" aria-labelledby="quotes-title">
        <div class="site-settings-heading"><div><p class="label">ROTATING QUOTES</p><h2 id="quotes-title">首页名言</h2></div><span>{{ quotes.length }} 条</span></div>
        <div class="quote-editor"><textarea v-model="newQuote" maxlength="500" rows="2" placeholder="添加一句首页名言" /><button :disabled="!newQuote.trim()" @click="addQuote">添加名言</button></div>
        <ul class="quote-admin-list"><li v-for="quote in quotes" :key="quote.id"><span>{{ quote.content }}</span><span class="admin-actions"><button class="secondary" @click="editQuote(quote)">编辑</button><button class="danger" @click="removeQuote(quote.id)">删除</button></span></li></ul>
        <p v-if="quoteMessage" class="form-success">{{ quoteMessage }}</p>
      </section>
      <section class="admin-list">
        <article v-for="article in articles" :key="article.id" class="admin-row">
          <div><h2>{{ article.title }}</h2><p>{{ article.slug }} · {{ article.status }}</p></div>
          <div class="admin-actions"><button class="secondary" @click="editing = { ...article }">编辑</button><button class="danger" @click="remove(article.id)">删除</button></div>
        </article>
        <p v-if="articles.length === 0">还没有文章。</p>
      </section>
    </template>
  </main>
</template>
