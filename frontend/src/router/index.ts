import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import StocksTableView from '@/views/StocksTableView.vue'
import InsightsView from '@/views/InsightsView.vue'

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
      component: StocksTableView,
    },
    {
      path: '/stock-insights',
      name: 'insights',
      component: InsightsView,
    },
  ],
})

export default router
