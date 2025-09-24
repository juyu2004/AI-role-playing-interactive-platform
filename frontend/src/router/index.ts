import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'roleSelect', component: () => import('../views/home/RoleSelectView.vue') },
    { path: '/chat/:roleId', name: 'chat', component: () => import('@/views/chat/ChatView.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/auth/LoginView.vue') },
    { path: '/register', name: 'register', component: () => import('@/views/auth/RegisterView.vue') },
    { path: '/profile', name: 'profile', component: () => import('@/views/profile/authProfile.vue') },
    // 添加角色详情页面路由
    { path: '/role/:roleId', name: 'roleProfile', component: () => import('@/views/profile/roleProfile.vue') }
  ],
})

export default router
