import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import StocksView from '@/views/StocksView.vue'
import RecommendationsView from '@/views/RecommendationsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/all-stocks',
      name: 'stocks',
      component: StocksView,
    },
    {
      path: '/stock-recommendations',
      name: 'recommendations',
      component: RecommendationsView,
    },
  ],
})

export default router
