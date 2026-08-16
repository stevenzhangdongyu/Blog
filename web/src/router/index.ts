import { createRouter, createWebHistory } from 'vue-router'
import ArticleView from '../views/ArticleView.vue'
import HomeView from '../views/HomeView.vue'

export default createRouter({
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
  ],
})
