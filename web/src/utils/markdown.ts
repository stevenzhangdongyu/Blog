import DOMPurify from 'dompurify'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import MarkdownIt from 'markdown-it'
import { katex } from '@mdit/plugin-katex'
import 'katex/dist/katex.min.css'

const md = new MarkdownIt({
  html: false,
  linkify: true,
  typographer: true,

  highlight(str: string, lang: string): string {
    if (!lang || !hljs.getLanguage(lang)) {
      return ''
    }

    return `<pre class="hljs"><code>${
      hljs.highlight(str, {
        language: lang,
        ignoreIllegals: true,
      }).value
    }</code></pre>`
  },
})

md.use(katex)

export function renderMarkdown(text: string): string {
  if (!text) return ''

  const html = md.render(text)
  return DOMPurify.sanitize(html)
}

export default md