import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import StocksTableView from '@/views/StocksTableView.vue'
import InsightsView from '@/views/InsightsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
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
