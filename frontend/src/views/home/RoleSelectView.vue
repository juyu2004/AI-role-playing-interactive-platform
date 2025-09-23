<template>
  <div class="role-select-container">
    <!-- 头部导航：登陆+注册+logo -->
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

    <!-- body部分 -->
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

    <!-- 角色卡片+模糊搜索 -->
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
          v-for="role in paginatedRoles"
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

      <!-- 分页控件 -->
      <div v-if="totalPages > 1" class="pagination-controls">
        <button
          class="pagination-btn"
          @click="goToPrevPage"
          :disabled="currentPage === 1"
        >
          上一页
        </button>
        <div class="page-numbers">
          <button
            v-for="page in totalPages"
            :key="page"
            class="page-btn"
            :class="{ 'active': currentPage === page }"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
        </div>
        <button
          class="pagination-btn"
          @click="goToNextPage"
          :disabled="currentPage === totalPages"
        >
          下一页
        </button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
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

// 分页相关变量
const currentPage = ref(1)
const pageSize = ref(4)

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

// 获取角色列表 - 修复mock数据不显示的问题
const fetchRoles = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await ApiService.getRoles()
    console.log('获取角色列表:', response)

    // 改进的响应处理逻辑
    if (response && response.success && response.data && Array.isArray(response.data) && response.data.length > 0) {
      // 标准响应格式，且有数据
      roles.value = response.data
    } else if (Array.isArray(response)) {
      // 如果直接返回数组
      roles.value = response
    } else {
      // 如果API返回为空或格式不符合预期，使用模拟数据
      roles.value = mockRoles
      console.log('使用模拟角色数据')
      error.value = '' // 关键修复：不显示错误信息
    }
    // 初始过滤
    handleSearch()
    // 重置页码
    currentPage.value = 1
  } catch (err) {
    console.error('获取角色列表失败:', err)
    // API请求失败时使用模拟数据
    roles.value = mockRoles
    console.log('使用模拟角色数据:', roles.value)
    error.value = '' // 关键修复：不显示错误信息
    // 初始过滤
    handleSearch()
    // 重置页码
    currentPage.value = 1
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    filteredRoles.value = roles.value
  } else {
    const query = searchQuery.value.toLowerCase().trim()
    filteredRoles.value = roles.value.filter(role =>
      role.name.toLowerCase().includes(query) ||
      role.category.toLowerCase().includes(query) ||
      role.description.toLowerCase().includes(query)
    )
  }
  // 搜索时重置到第一页
  currentPage.value = 1
}

// 计算当前页显示的角色
const paginatedRoles = computed(() => {
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return filteredRoles.value.slice(startIndex, endIndex)
})

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(filteredRoles.value.length / pageSize.value)
})

// 跳转到指定页
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

// 下一页
const goToNextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// 上一页
const goToPrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
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
logo-container {
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

/* 分页控件样式 */
.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
  padding: 1rem;
}

.pagination-btn {
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pagination-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 0.5rem;
}

.page-btn {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.page-btn.active {
  background: white;
  color: #667eea;
  border-color: white;
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

  /* 分页控件响应式调整 */
  .pagination-controls {
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .pagination-btn {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }

  .page-btn {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }
}
</style>

