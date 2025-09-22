<template>
  <div class="role-select-container">
    <!-- 新的头部导航 -->
    <div class="top-navigation">
      <!-- 左侧Logo -->
      <div class="logo-container">
        <div class="logo">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="rgba(255,255,255,0.9)"/>
            <path d="M2 17L12 22L22 17" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
            <path d="M2 12L12 17L22 12" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="logo-text">AI角色对话</span>
        </div>
      </div>

      <!-- 右侧用户控制区域 -->
      <div class="user-controls">
        <!-- 未登录状态显示登录注册按钮 -->
        <template v-if="!isLoggedIn">
          <button class="login-btn" @click='handleLogin'>登录</button>
          <button class="register-btn" @click="handleRegister">注册</button>
        </template>

        <!-- 已登录状态显示用户头像和退出按钮 -->
        <template v-else>
          <div class="user-profile">
            <img
              :src="userAvatar || 'https://picsum.photos/id/64/40/40'"
              :alt="'用户头像'"
              class="user-avatar"
            />
            <button class="logout-btn" @click="handleLogout">退出登录</button>
          </div>
        </template>
      </div>
    </div>

    <!-- 原有的页面内容 -->
    <header class="header">
      <h1 class="title">Choose your Role to Talk</h1>
      <div class="search-container">
        <input
          type="text"
          v-model="searchQuery"
          placeholder="搜索角色..."
          class="search-input"
          @input="handleSearch"
        />
        <svg class="search-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M21 21L16.65 16.65M19 10.5C19 15.1944 15.1944 19 10.5 19C5.80558 19 2 15.1944 2 10.5C2 5.80558 5.80558 2 10.5 2C15.1944 2 19 5.80558 19 10.5Z" stroke="#6B7280" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
    </header>

    <!-- 角色卡片部分保持不变 -->
    <main class="roles-content">
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>加载角色中...</p>
      </div>

      <div v-else-if="error" class="error-container">
        <p>{{ error }}</p>
        <button class="retry-button" @click="fetchRoles">重试</button>
      </div>

      <div v-else-if="filteredRoles.length === 0" class="empty-container">
        <p>未找到匹配的角色</p>
      </div>

      <div v-else class="roles-grid">
        <div
          v-for="role in filteredRoles"
          :key="role.id"
          class="role-card"
          @click="selectRole(role.id)"
          @mouseenter="hoveredRoleId = role.id"
          @mouseleave="hoveredRoleId = ''"
        >
          <div class="avatar-container">
            <img
              :src="role.avatarUrl || 'https://picsum.photos/id/' + (parseInt(role.id) % 50) + '/200/200'"
              :alt="role.name"
              class="role-avatar"
            />
            <div class="avatar-overlay" :class="{ 'active': hoveredRoleId === role.id }">
              <div class="view-details">开始对话</div>
            </div>
          </div>
          <div class="role-info">
            <h3 class="role-name">{{ role.name }}</h3>
            <p class="role-category">{{ role.category }}</p>
            <p class="role-description">{{ role.description }}</p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ApiService from '@/services/api'
import type { Role } from '@/types/api'

const router = useRouter()
const roles = ref<Role[]>([])
const filteredRoles = ref<Role[]>([])
const searchQuery = ref('')
const loading = ref(false)
const error = ref('')
const hoveredRoleId = ref('')

// 新增用户状态相关变量
const isLoggedIn = ref(false) // 模拟用户登录状态
const userAvatar = ref('') // 模拟用户头像

// 模拟角色数据，当API不可用时使用
const mockRoles: Role[] = [
  {
    id: '1',
    name: '夏洛克·福尔摩斯',
    category: '文学角色',
    avatarUrl: 'https://picsum.photos/id/1/200/200',
    description: '英国著名侦探，善于观察和推理'
  },
  {
    id: '2',
    name: '花木兰',
    category: '历史人物',
    avatarUrl: 'https://picsum.photos/id/2/200/200',
    description: '代父从军的女英雄'
  },
  {
    id: '3',
    name: '爱因斯坦',
    category: '科学家',
    avatarUrl: 'https://picsum.photos/id/3/200/200',
    description: '相对论的创立者'
  },
  {
    id: '4',
    name: '莎士比亚',
    category: '文学家',
    avatarUrl: 'https://picsum.photos/id/4/200/200',
    description: '英国著名剧作家和诗人'
  },
  {
    id: '5',
    name: '居里夫人',
    category: '科学家',
    avatarUrl: 'https://picsum.photos/id/5/200/200',
    description: '首位获得两次诺贝尔奖的科学家'
  },
  {
    id: '6',
    name: '李白',
    category: '诗人',
    avatarUrl: 'https://picsum.photos/id/6/200/200',
    description: '唐代著名浪漫主义诗人'
  }
]

// 新增登录处理函数
const handleLogin = () => {
  // 实际项目中这里应该跳转到登录页面
  console.log('登录按钮点击')
  // 模拟登录成功
  router.push({ name: 'login' })
  isLoggedIn.value = true
  userAvatar.value = 'https://picsum.photos/id/64/40/40'
  // 跳转到角色选择页面
}

// 新增注册处理函数
const handleRegister = () => {
  console.log('注册按钮点击')
  // 模拟跳转到注册页面
  router.push({ name: 'register' })
}

// 新增退出登录处理函数
const handleLogout = () => {
  // 实际项目中这里应该清除用户登录状态
  console.log('退出登录按钮点击')
  isLoggedIn.value = false
  userAvatar.value = ''
}

// 获取角色列表
const fetchRoles = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await ApiService.getRoles()
    console.log('获取角色列表:', response)

    // 检查响应是否为数组，而不是检查success属性
    if (Array.isArray(response) && response.length > 0) {
      roles.value = response
    } else if (response.success && response.data && response.data.length > 0) {
      // 保留原有的处理逻辑，以便兼容可能的标准响应格式
      roles.value = response.data
    } else {
      // 如果API返回为空，使用模拟数据
      roles.value = mockRoles
    }
    // 初始过滤
    handleSearch()
  } catch (err) {
    console.error('获取角色列表失败:', err)
    // API请求失败时使用模拟数据
    roles.value = mockRoles
    error.value = '获取角色列表失败，使用模拟数据'
    // 初始过滤
    handleSearch()
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    filteredRoles.value = roles.value
    return
  }

  const query = searchQuery.value.toLowerCase().trim()
  filteredRoles.value = roles.value.filter(role =>
    role.name.toLowerCase().includes(query) ||
    role.category.toLowerCase().includes(query) ||
    role.description.toLowerCase().includes(query)
  )
}

// 选择角色并跳转到对话界面
const selectRole = (roleId: string) => {
  router.push({ name: 'chat', params: { roleId } })
}

// 组件挂载时获取角色列表
onMounted(() => {
  fetchRoles()
  // 模拟检查用户登录状态
  // 在实际项目中，这里应该从localStorage或API检查用户登录状态
  // isLoggedIn.value = localStorage.getItem('userToken') !== null
})
</script>

<style scoped src="../home/index.css"></style>

<style scoped>
/* 新增的导航栏样式 */
.top-navigation {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  margin-bottom: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* Logo样式 */
.logo-container {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.logo-text {
  color: white;
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 0.5px;
}

/* 用户控制区域 */
.user-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

/* 登录和注册按钮 */
.login-btn,
.register-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.login-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.login-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.register-btn {
  background: white;
  color: #667eea;
}

.register-btn:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

/* 用户资料区域 */
.user-profile {
  display: flex;
  align-items: center;
  gap: 1rem;
}

/* 用户头像 */
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.5);
  object-fit: cover;
}

/* 退出登录按钮 */
.logout-btn {
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .top-navigation {
    padding: 1rem;
  }

  .logo-text {
    font-size: 1.2rem;
  }

  .login-btn,
  .register-btn {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }

  .user-controls {
    gap: 0.5rem;
  }
}
</style>

