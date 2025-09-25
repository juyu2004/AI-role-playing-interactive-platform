<template>
  <div class="profile-container">
    <!-- 头部导航 -->
    <div class="top-navigation">
      <div class="logo-container">
        <div class="logo">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="rgba(255,255,255,0.9)"/>
            <path d="M2 17L12 22L22 17" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
            <path d="M2 12L12 17L22 12" stroke="rgba(255,255,255,0.9)" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="logo-text">个人中心</span>
        </div>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="profile-content">
      <!-- 侧边栏导航 -->
      <div class="profile-sidebar">
        <div class="sidebar-menu">
          <button 
            class="menu-item" 
            :class="{ active: activeSection === 'personal-info' }"
            @click="activeSection = 'personal-info'"
          >
            <svg class="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 12C14.7614 12 17 9.76142 17 7C17 4.23858 14.7614 2 12 2C9.23858 2 7 4.23858 7 7C7 9.76142 9.23858 12 12 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M20.42 14.82C21.6427 16.1828 22.2477 17.9812 22.0259 19.8705C21.8042 21.7599 20.7732 23.4527 19.1431 24.57C17.513 25.6873 15.4533 26.1322 13.4286 25.7857C11.4038 25.4392 9.55878 24.3161 8.25 22.75V14.82" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            个人信息
          </button>
          <button 
            class="menu-item" 
            :class="{ active: activeSection === 'favorite-roles' }"
            @click="activeSection = 'favorite-roles'"
          >
            <svg class="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            收藏角色
          </button>
          <button 
            class="menu-item" 
            :class="{ active: activeSection === 'friends' }"
            @click="activeSection = 'friends'"
          >
            <svg class="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="9" cy="7" r="4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M22 21v-2a4 4 0 0 0-3-3.87" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            个人好友
          </button>
          <button 
            class="menu-item" 
            :class="{ active: activeSection === 'interests' }"
            @click="activeSection = 'interests'"
          >
            <svg class="menu-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <polyline points="22 4 12 14.01 9 11.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            爱好领域
          </button>
        </div>
      </div>

      <!-- 主内容区域 -->
      <div class="profile-main">
        <!-- 个人信息部分 -->
        <div v-if="activeSection === 'personal-info'" class="info-section">
          <h2 class="section-title">个人信息</h2>
          <div class="user-info">
            <div class="avatar-container">
              <img 
                :src="userInfo.avatarUrl" 
                alt="用户头像"
                class="user-avatar"
              />
              <button class="change-avatar-btn">更换头像</button>
            </div>
            <div class="basic-info">
              <div class="info-item">
                <label>用户名：</label>
                <span>{{ userInfo.username }}</span>
              </div>
              <div class="info-item">
                <label>邮箱：</label>
                <span>{{ userInfo.email }}</span>
              </div>
              <div class="info-item">
                <label>注册时间：</label>
                <span>{{ formatDate(userInfo.registerDate) }}</span>
              </div>
              <div class="info-item">
                <label>上次登录：</label>
                <span>{{ formatDate(userInfo.lastLogin) }}</span>
              </div>
            </div>
          </div>
          <button class="edit-profile-btn">编辑个人资料</button>
        </div>

        <!-- 收藏角色部分 -->
        <div v-if="activeSection === 'favorite-roles'" class="roles-section">
          <h2 class="section-title">收藏角色</h2>
          <div class="roles-grid">
            <div v-for="role in favoriteRoles" :key="role.id" class="role-card">
              <div class="role-avatar" :style="{ backgroundImage: `url(${role.avatarUrl})` }"></div>
              <h3 class="role-name">{{ role.name }}</h3>
              <p class="role-category">{{ role.category }}</p>
              <button class="chat-btn" @click="startChatWithRole(role.id)">开始对话</button>
              <button class="unfavorite-btn" @click="unfavoriteRole(role.id)">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z" fill="currentColor"/>
                </svg>
                取消收藏
              </button>
            </div>
          </div>
          <p v-if="favoriteRoles.length === 0" class="empty-message">您还没有收藏任何角色</p>
        </div>

        <!-- 个人好友部分 -->
        <div v-if="activeSection === 'friends'" class="friends-section">
          <h2 class="section-title">个人好友</h2>
          <div class="friends-list">
            <div v-for="friend in friends" :key="friend.id" class="friend-item">
              <img :src="friend.avatarUrl" alt="好友头像" class="friend-avatar" />
              <div class="friend-info">
                <h3 class="friend-name">{{ friend.username }}</h3>
                <p class="friend-status">{{ friend.isOnline ? '在线' : '离线' }}</p>
              </div>
              <div class="friend-actions">
                <button class="message-btn">发消息</button>
                <button class="remove-friend-btn">移除好友</button>
              </div>
            </div>
          </div>
          <p v-if="friends.length === 0" class="empty-message">您还没有添加任何好友</p>
          <button class="add-friend-btn">添加好友</button>
        </div>

        <!-- 爱好领域部分 -->
        <div v-if="activeSection === 'interests'" class="interests-section">
          <h2 class="section-title">爱好领域</h2>
          <div class="interests-tags">
            <div 
              v-for="interest in allInterests" 
              :key="interest.id" 
              class="interest-tag"
              :class="{ selected: userInterests.includes(interest.id) }"
              @click="toggleInterest(interest.id)"
            >
              {{ interest.name }}
            </div>
          </div>
          <button class="save-interests-btn">保存设置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import type { Role } from '@/types/api';

const router = useRouter();

// 当前激活的部分
const activeSection = ref('personal-info');

// 用户信息（模拟数据）
const userInfo = reactive({
  username: '用户123',
  email: 'user@example.com',
  avatarUrl: 'https://picsum.photos/id/64/200/200',
  registerDate: '2023-01-15T08:30:00Z',
  lastLogin: '2023-10-20T14:45:00Z'
});

// 收藏的角色（模拟数据）
const favoriteRoles = ref<Role[]>([
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
    id: '5',
    name: '居里夫人',
    category: '科学家',
    avatarUrl: 'https://picsum.photos/id/5/200/200',
    description: '首位获得两次诺贝尔奖的科学家'
  }
]);

// 好友列表（模拟数据）
const friends = ref([
  {
    id: '1',
    username: '张三',
    avatarUrl: 'https://picsum.photos/id/237/200/200',
    isOnline: true
  },
  {
    id: '2',
    username: '李四',
    avatarUrl: 'https://picsum.photos/id/177/200/200',
    isOnline: false
  },
  {
    id: '3',
    username: '王五',
    avatarUrl: 'https://picsum.photos/id/342/200/200',
    isOnline: true
  }
]);

// 所有可选的爱好领域
const allInterests = ref([
  { id: '1', name: '文学' },
  { id: '2', name: '历史' },
  { id: '3', name: '科学' },
  { id: '4', name: '艺术' },
  { id: '5', name: '音乐' },
  { id: '6', name: '电影' },
  { id: '7', name: '体育' },
  { id: '8', name: '技术' },
  { id: '9', name: '旅行' },
  { id: '10', name: '美食' }
]);

// 用户已选择的爱好领域
const userInterests = ref<string[]>(['1', '3', '8']);

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};

// 开始与角色对话
const startChatWithRole = (roleId: string) => {
  router.push({ name: 'chat', params: { roleId } });
};

// 取消收藏角色
const unfavoriteRole = (roleId: string) => {
  favoriteRoles.value = favoriteRoles.value.filter(role => role.id !== roleId);
};

// 切换爱好领域选择
const toggleInterest = (interestId: string) => {
  const index = userInterests.value.indexOf(interestId);
  if (index > -1) {
    userInterests.value.splice(index, 1);
  } else {
    userInterests.value.push(interestId);
  }
};
</script>

<style scoped>
/* 个人中心容器 */
.profile-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
}

/* 头部导航 */
.top-navigation {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

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

/* 主内容区域 */
.profile-content {
  display: flex;
  gap: 2rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

/* 侧边栏导航 */
.profile-sidebar {
  width: 200px;
  flex-shrink: 0;
}

.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: 500;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(5px);
}

.menu-item.active {
  background: white;
  color: #667eea;
}

.menu-icon {
  flex-shrink: 0;
}

/* 主内容 */
.profile-main {
  flex: 1;
  min-height: 600px;
}

/* 通用区域样式 */
.section-title {
  color: white;
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid rgba(255, 255, 255, 0.2);
}

.empty-message {
  color: rgba(255, 255, 255, 0.7);
  text-align: center;
  padding: 2rem;
  font-size: 1.1rem;
}

/* 个人信息部分 */
.user-info {
  display: flex;
  gap: 2rem;
  margin-bottom: 2rem;
}

.avatar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.user-avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.3);
  object-fit: cover;
}

.change-avatar-btn {
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.9rem;
}

.change-avatar-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.basic-info {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  padding: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
  font-weight: 500;
}

.info-item span {
  color: white;
  font-size: 1rem;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.edit-profile-btn {
  padding: 0.75rem 2rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.edit-profile-btn:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

/* 收藏角色部分 */
.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.role-card {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s ease;
}

.role-card:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-5px);
}

.role-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background-size: cover;
  background-position: center;
  border: 3px solid rgba(255, 255, 255, 0.3);
}

.role-name {
  color: white;
  font-size: 1.2rem;
  font-weight: 600;
  text-align: center;
}

.role-category {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.9rem;
  text-align: center;
}

.chat-btn {
  padding: 0.5rem 1.5rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 100%;
}

.chat-btn:hover {
  background: #f0f0f0;
}

.unfavorite-btn {
  padding: 0.4rem 1rem;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.unfavorite-btn:hover {
  background: rgba(255, 0, 0, 0.2);
}

/* 好友列表部分 */
.friends-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 2rem;
}

.friend-item {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 1rem;
  transition: all 0.3s ease;
}

.friend-item:hover {
  background: rgba(255, 255, 255, 0.15);
}

.friend-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.friend-info {
  flex: 1;
}

.friend-name {
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.friend-status {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.9rem;
}

.friend-status:before {
  content: '';
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 0.5rem;
  background-color: currentColor;
}

.friend-status[text="在线"] {
  color: #4ade80;
}

.friend-actions {
  display: flex;
  gap: 0.5rem;
}

.message-btn, .remove-friend-btn {
  padding: 0.4rem 1rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.message-btn {
  background: white;
  color: #667eea;
}

.message-btn:hover {
  background: #f0f0f0;
}

.remove-friend-btn {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.remove-friend-btn:hover {
  background: rgba(255, 0, 0, 0.2);
}

.add-friend-btn {
  padding: 0.75rem 2rem;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.add-friend-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 爱好领域部分 */
.interests-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.interest-tag {
  padding: 0.75rem 1.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 25px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 1rem;
  font-weight: 500;
}

.interest-tag:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.interest-tag.selected {
  background: white;
  color: #667eea;
  border-color: white;
}

.save-interests-btn {
  padding: 0.75rem 2rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.save-interests-btn:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .profile-content {
    flex-direction: column;
  }
  
  .profile-sidebar {
    width: 100%;
  }
  
  .sidebar-menu {
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 0.5rem;
  }
  
  .menu-item {
    white-space: nowrap;
    flex: 0 0 auto;
  }
  
  .basic-info {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 1rem;
  }
  
  .user-info {
    flex-direction: column;
    text-align: center;
  }
  
  .roles-grid {
    grid-template-columns: 1fr;
  }
  
  .friend-item {
    flex-direction: column;
    text-align: center;
  }
  
  .friend-actions {
    width: 100%;
    justify-content: center;
  }
}
</style>