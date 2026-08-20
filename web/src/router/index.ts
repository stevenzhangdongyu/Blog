import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/articles/:slug',
      component: () => import('../views/ArticleView.vue'),
    },
    {
      path: '/admin/login',
      component: () => import('../views/AdminLoginView.vue'),
    },
    {
      path: '/admin/articles',
      component: () => import('../views/AdminArticlesView.vue'),
      meta: { requiresAdmin: true },
    },
  ],
})

router.beforeEach(async (to) => {
  document.documentElement.classList.add('is-navigating')
  if (!to.meta.requiresAdmin) return true
  const response = await fetch('/api/admin/session', { credentials: 'include' })
  return response.ok ? true : '/admin/login'
})

router.afterEach(() => document.documentElement.classList.remove('is-navigating'))
router.onError(() => document.documentElement.classList.remove('is-navigating'))

export default router
