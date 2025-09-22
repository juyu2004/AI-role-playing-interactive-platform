import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'roleSelect',
      component: () => import('../views/RoleSelectView.vue')
    },
    {
      path: '/chat/:roleId',
      name: 'chat',
      component: () => import('@/views/ChatView.vue')
    }
  ],
})

export default router
