import { createRouter, createWebHistory } from 'vue-router'
import ArticleView from '../views/ArticleView.vue'
import HomeView from '../views/HomeView.vue'
import AdminLoginView from '../views/AdminLoginView.vue'
import AdminArticlesView from '../views/AdminArticlesView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: HomeView,
    },
    {
      path: '/articles/:slug',
      component: ArticleView,
    },
    {
      path: '/admin/login',
      component: AdminLoginView,
    },
    {
      path: '/admin/articles',
      component: AdminArticlesView,
      meta: { requiresAdmin: true },
    },
  ],
})

router.beforeEach(async (to) => {
  if (!to.meta.requiresAdmin) return true
  const response = await fetch('/api/admin/session', { credentials: 'include' })
  return response.ok ? true : '/admin/login'
})

export default router
