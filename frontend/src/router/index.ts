import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'roleSelect',
      component: () => import('../views/home/RoleSelectView.vue')
    },
    {
      path: '/chat/:roleId',
      name: 'chat',
      component: () => import('@/views/chat/ChatView.vue')
    },
    // 登陆路由
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue')
    },
    // 注册路由
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue')
    }
  ],
})

export default router
