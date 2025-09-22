<template>
  <div class="role-select-container">
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
})
</script>

<style scoped>
.role-select-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
}

.header {
  max-width: 1200px;
  margin: 0 auto 2rem;
  text-align: center;
}

.title {
  font-size: 2.5rem;
  color: white;
  margin-top: 100px;
  margin-bottom: 1.5rem;
  font-weight: 700;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.search-container {
  margin-top: 30px;
  position: relative;
  max-width: 500px;
  margin: 0 auto;
}

.search-input {
  width: 100%;
  padding: 1rem 3rem;
  border: none;
  border-radius: 50px;
  font-size: 1rem;
  outline: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
}

.search-input:focus {
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
}

.roles-content {
  max-width: 1200px;
  margin: 0 auto;
}

.loading-container,
.error-container,
.empty-container {
  text-align: center;
  padding: 4rem;
  color: white;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.retry-button {
  margin-top: 1rem;
  padding: 0.5rem 1.5rem;
  border: none;
  border-radius: 25px;
  background: white;
  color: #667eea;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-button:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.role-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.role-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15);
}

.avatar-container {
  position: relative;
  width: 100%;
  height: 160px;
  overflow: hidden;
}

.role-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.role-info {
  padding: 1.2rem;
}

.role-name {
  font-size: 1.3rem;
  margin-bottom: 0.5rem;
  color: #333;
}

.role-category {
  display: inline-block;
  background: #e1e8ed;
  color: #657786;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.8rem;
  margin-bottom: 0.8rem;
}

.role-description {
  color: #666;
  line-height: 1.5;
  font-size: 0.9rem;
}

.role-card:hover .role-avatar {
  transform: scale(1.1);
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 1.5rem;
}

.avatar-overlay.active {
  opacity: 1;
}

.view-details {
  color: white;
  font-weight: 600;
  font-size: 1.1rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  background: rgba(102, 126, 234, 0.9);
  padding: 0.5rem 1.5rem;
  border-radius: 25px;
  transition: all 0.3s ease;
}

.role-info {
  padding: 1.5rem;
}

.role-name {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: #333;
}

.role-category {
  display: inline-block;
  background: #e1e8ed;
  color: #657786;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.85rem;
  margin-bottom: 1rem;
}

.role-description {
  color: #666;
  line-height: 1.6;
  font-size: 0.95rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .role-select-container {
    padding: 1rem;
  }

  .title {
    font-size: 2rem;
  }

  .roles-grid {
    grid-template-columns: 1fr;
  }
}
</style>
