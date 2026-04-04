import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'lista',
      component: () => import('./views/ListaView.vue')
    },
    {
      path: '/compras',
      name: 'compras',
      component: () => import('./views/ComprasView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('./views/DashboardView.vue')
    },
    {
      path: '/config',
      name: 'config',
      component: () => import('./views/ConfigView.vue')
    }
  ]
})

export default router
