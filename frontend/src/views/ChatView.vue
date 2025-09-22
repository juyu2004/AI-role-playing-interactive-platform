<template>
  <div class="chat-container">
    <!-- 顶部导航栏 -->
    <div class="chat-header">
      <button class="back-button" @click="goBack" aria-label="返回">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7"></path>
        </svg>
      </button>
      <div class="role-info">
        <img
          :src="currentRole?.avatarUrl || `https://picsum.photos/id/${parseInt(roleId) % 50}/200/200`"
          :alt="currentRole?.name || '角色头像'"
          class="role-avatar"
        >
        <div class="role-details">
          <h2 class="role-name">{{ currentRole?.name || '加载中...' }}</h2>
          <p class="role-status">{{ isTyping ? '正在输入...' : currentRole?.category }}</p>
        </div>
      </div>
      <!-- 将图标按钮改为文本按钮 -->
      <button class="clear-history-button" @click="clearChatHistory" aria-label="清除历史">
        清除历史
      </button>
    </div>

    <!-- 对话消息区域 -->
    <main class="chat-messages" ref="chatMessagesRef">
      <div v-if="loadingMessages" class="loading-messages">
        <div class="loading-spinner"></div>
        <p>加载对话历史...</p>
      </div>

      <div v-else-if="error" class="error-messages">
        <p>{{ error }}</p>
        <button class="retry-button" @click="loadChatHistory">重试</button>
      </div>

      <div v-else-if="messages.length === 0" class="empty-messages">
        <p>开始与{{ currentRole?.name }}的对话吧！</p>
      </div>

      <div v-for="message in messages" :key="message.id" class="message-wrapper" :class="{ 'user-message': message.isUser }">
        <div v-if="!message.isUser" class="avatar-small">
          <img
            :src="currentRole?.avatarUrl || 'https://picsum.photos/id/' + (parseInt(roleId) % 50) + '/200/200'"
            :alt="currentRole?.name"
          />
        </div>
        <div class="message-bubble" :class="{ 'user': message.isUser }">{{ message.text }}</div>
        <div v-if="message.isUser" class="avatar-small">
          <div class="user-avatar">
            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 12C14.7614 12 17 9.76142 17 7C17 4.23858 14.7614 2 12 2C9.23858 2 7 4.23858 7 7C7 9.76142 9.23858 12 12 12Z" stroke="#667eea" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M20.42 14.82C21.6427 16.1828 22.2477 17.9812 22.0259 19.8705C21.8042 21.7599 20.7732 23.4527 19.1431 24.57C17.513 25.6873 15.4533 26.1322 13.4286 25.7857C11.4038 25.4392 9.55878 24.3161 8.25 22.75V14.82" stroke="#667eea" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
        </div>
      </div>

      <!-- 正在输入提示 -->
      <div v-if="isTyping" class="typing-indicator">
        <div class="avatar-small">
          <img
            :src="currentRole?.avatarUrl || 'https://picsum.photos/id/' + (parseInt(roleId) % 50) + '/200/200'"
            :alt="currentRole?.name"
          />
        </div>
        <div class="typing-dots">
          <div class="dot"></div>
          <div class="dot"></div>
          <div class="dot"></div>
        </div>
      </div>
    </main>

    <!-- 底部输入区域 -->
    <footer class="chat-input-container">
      <div class="input-wrapper">
        <input
          type="text"
          v-model="inputMessage"
          placeholder="输入消息..."
          class="chat-input"
          @keyup.enter="sendMessage"
          :disabled="isTyping"
        />
        <button
          class="send-button"
          @click="sendMessage"
          :disabled="!inputMessage.trim() || isTyping"
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M22 2L11 13M22 2L15 21L11 13M22 2H2" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import ApiService from '@/services/api'
import type { Role, ChatRequest } from '@/types/api'

const router = useRouter()
const route = useRoute()
const chatMessagesRef = ref<HTMLElement>()
const roleId = route.params.roleId as string

// 状态管理
const currentRole = ref<Role | null>(null)
const messages = ref<Array<{ id: string, text: string, isUser: boolean, timestamp: string, audioUrl?: string | null }>>([])
const inputMessage = ref('')
const isTyping = ref(false)
const loadingMessages = ref(false)
const error = ref('')

// 本地存储键名
const STORAGE_KEY = `chat_history_${roleId}`

// 保存对话历史到本地存储
const saveChatHistory = () => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      roleId,
      roleName: currentRole.value?.name,
      messages: messages.value,
      lastUpdated: new Date().toISOString()
    }))
    console.log('对话历史已保存到本地存储')
  } catch (err) {
    console.error('保存对话历史失败:', err)
  }
}

// 加载角色信息
const loadRoleInfo = async () => {
  try {
    // 从API获取角色详情
    const response = await ApiService.getRoleDetail(roleId)

    // 根据OpenAPI规范，API直接返回Role对象，而不是BaseResponse包装的对象
    // 因此我们需要修改判断逻辑
    if (response && typeof response === 'object') {
      // 检查是否是直接的Role对象
      if ('id' in response && 'name' in response) {
        currentRole.value = response as unknown as Role
        console.log('成功加载角色信息:', currentRole.value)
      }
      // 同时保留对BaseResponse格式的兼容，以防后端将来改变格式
      else if (response.success !== undefined && response.data) {
        currentRole.value = response.data
        console.log('成功加载角色信息(从BaseResponse):', currentRole.value)
      } else {
        throw new Error('无效的角色数据格式')
      }
    } else {
      throw new Error('无效的响应格式')
    }
  } catch (err) {
    console.error('加载角色信息失败:', err)
    error.value = '加载角色信息失败'
    // 错误时使用模拟数据
    const mockRole: Role = {
      id: roleId,
      name: ['夏洛克·福尔摩斯', '花木兰', '爱因斯坦', '莎士比亚', '居里夫人', '李白'][parseInt(roleId) % 6],
      category: ['文学角色', '历史人物', '科学家', '文学家', '科学家', '诗人'][parseInt(roleId) % 6],
      avatarUrl: `https://picsum.photos/id/${parseInt(roleId) % 50}/200/200`,
      description: '正在对话中的角色'
    }
    currentRole.value = mockRole
  }
}

// 加载对话历史
const loadChatHistory = async () => {
  loadingMessages.value = true
  error.value = ''
  try {
    // 从本地存储加载对话历史
    const storedHistory = localStorage.getItem(STORAGE_KEY)
    if (storedHistory) {
      try {
        const parsedHistory = JSON.parse(storedHistory)
        if (parsedHistory.messages && Array.isArray(parsedHistory.messages)) {
          messages.value = parsedHistory.messages
          console.log('从本地存储加载对话历史成功:', messages.value.length, '条消息')
          return
        }
      } catch (parseError) {
        console.error('解析本地存储的对话历史失败:', parseError)
      }
    }
    // 如果本地没有存储的历史或解析失败，使用空消息列表
    messages.value = []
    console.log('没有找到本地存储的对话历史，使用空消息列表')
  } catch (err) {
    console.error('加载对话历史失败:', err)
    error.value = '加载对话历史失败'
    messages.value = []
  } finally {
    loadingMessages.value = false
  }
}

// 清除本地对话历史
const clearChatHistory = () => {
  if (confirm('确定要结束当前对话历史吗？对话记录不可恢复')) {
    try {
      localStorage.removeItem(STORAGE_KEY)
      messages.value = []
      console.log('对话历史已清除')
    } catch (err) {
      console.error('清除对话历史失败:', err)
    }
  }
}

// 发送消息
const sendMessage = async () => {
  const message = inputMessage.value.trim()
  console.log('发送消息:', message)
  if (!message || isTyping.value) return

  // 添加用户消息到列表
  const userMessage = {
    id: `user-${Date.now()}`,
    text: message,
    isUser: true,
    timestamp: new Date().toISOString()
  }
  messages.value.push(userMessage)
  inputMessage.value = ''

  // 保存到本地存储
  saveChatHistory()

  // 滚动到底部
  scrollToBottom()

  // 模拟正在输入
  isTyping.value = true

  try {
    // 构建聊天请求
    const chatRequest: ChatRequest = {
      roleId,
      text: message
    }

    // 调用API发送消息
    const response = await ApiService.sendChat(chatRequest)
    console.log('发送消息响应:', response)

    // 处理API响应，根据实际返回格式提取文本内容
    let replyText = '抱歉，我暂时无法回答这个问题。'
    let replyAudioUrl: string | null = null

    // 检查响应格式
    if (response && typeof response === 'object') {
      // 情况1: 响应是直接的ChatResponse对象 (根据OpenAPI规范)
      if ('text' in response) {
        replyText = response.text as string
        replyAudioUrl = 'audioUrl' in response ? (response as Record<string, unknown>).audioUrl as string | null || null : null
      }
      // 情况2: 响应是BaseResponse格式 (兼容之前的实现)
      else if (response.data && typeof response.data === 'object') {
        replyText = response.data.text || replyText
        replyAudioUrl = response.data.audioUrl || null
      }
    }

    // 添加角色回复到列表
    const roleReply = {
      id: `role-${Date.now()}`,
      text: replyText,
      isUser: false,
      timestamp: new Date().toISOString(),
      audioUrl: replyAudioUrl
    }
    messages.value.push(roleReply)

    // 保存到本地存储
    saveChatHistory()

    // 如果有audioUrl，可以在这里处理音频播放逻辑
    if (roleReply.audioUrl) {
      console.log('角色回复包含音频:', roleReply.audioUrl)
      // 实际项目中可以在这里实现音频播放
    }
  } catch (err) {
    console.error('发送消息失败:', err)
    // API失败时使用模拟回复
    const mockReply = {
      id: `role-${Date.now()}`,
      text: '我是' + (currentRole.value?.name || 'AI助手') + '，很高兴与你对话！',
      isUser: false,
      timestamp: new Date().toISOString(),
      audioUrl: null
    }
    messages.value.push(mockReply)

    // 保存到本地存储
    saveChatHistory()
  } finally {
    // 停止正在输入提示
    isTyping.value = false
    // 滚动到底部
    scrollToBottom()
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatMessagesRef.value) {
      chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
    }
  })
}

// 返回上一页
const goBack = () => {
  router.push('/')
}

// 组件挂载时加载数据
onMounted(() => {
  loadRoleInfo()
  loadChatHistory()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f8fafc;
}

/* 顶部导航栏 */
.chat-header {
  display: flex;
  align-items: center;
  padding: 1rem;
  background: white;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.back-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 50%;
  transition: background-color 0.3s ease;
  color: #64748b;
}

.back-button:hover {
  background: #f1f5f9;
}

.role-info {
  display: flex;
  align-items: center;
  margin-left: 1rem;
  flex: 1;
}

.role-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.role-details {
  margin-left: 1rem;
}

.role-name {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #1e293b;
}

.role-status {
  margin: 0;
  font-size: 0.9rem;
  color: #64748b;
}

/* 优化清除历史按钮样式 */
.clear-history-button {
  background: #1e293b;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: background-color 0.3s ease;
}

.clear-history-button:hover {
  background: #667eea;
  color: white;
}

.clear-history-button:hover {
  background: #054a8f;
}

.clear-history-button:hover {
  color: #c3cfe2;
}
.role-info {
  display: flex;
  align-items: center;
  margin-left: 1rem;
  flex: 1;
}

.role-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 1rem;
  border: 2px solid #667eea;
}

.role-name {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.role-status {
  font-size: 0.875rem;
  color: #64748b;
  margin: 0;
}

/* 对话消息区域 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: linear-gradient(135deg, #ababe0 0%, #e2c3d9 100%);
}

.loading-messages,
.error-messages,
.empty-messages {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #64748b;
  text-align: center;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(102, 126, 234, 0.3);
  border-radius: 50%;
  border-top-color: #667eea;
  animation: spin 1s ease-in-out infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.retry-button {
  margin-top: 1rem;
  padding: 0.5rem 1.5rem;
  border: none;
  border-radius: 25px;
  background: #667eea;
  color: white;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-button:hover {
  background: #5a67d8;
  transform: translateY(-2px);
}

/* 消息样式 */
.message-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  max-width: 75%;
}

.message-wrapper.user-message {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.avatar-small {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.avatar-small img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.user-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message-bubble {
  padding: 0.75rem 1rem;
  border-radius: 18px;
  font-size: 0.95rem;
  line-height: 1.4;
  word-wrap: break-word;
}

.message-bubble:not(.user) {
  background: white;
  color: #1e293b;
  border-bottom-left-radius: 6px;
}

.message-bubble.user {
  background: #667eea;
  color: white;
  border-bottom-right-radius: 6px;
}

/* 正在输入提示 */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  max-width: 75%;
}

.typing-dots {
  display: flex;
  gap: 4px;
  padding: 0.5rem 1rem;
  background: white;
  border-radius: 18px;
  border-bottom-left-radius: 6px;
}

.typing-dots .dot {
  width: 8px;
  height: 8px;
  background: #667eea;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out both;
}

.typing-dots .dot:nth-child(1) { animation-delay: -0.32s; }
.typing-dots .dot:nth-child(2) { animation-delay: -0.16s; }

@keyframes typing {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

/* 底部输入区域 */
.chat-input-container {
  padding: 1rem;
  background: white;
  border-top: 1px solid #e2e8f0;
  box-shadow: 0 -1px 3px rgba(0, 0, 0, 0.1);
}

.input-wrapper {
  display: flex;
  gap: 0.75rem;
  max-width: 800px;
  margin: 0 auto;
}

.chat-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 25px;
  font-size: 0.95rem;
  outline: none;
  transition: border-color 0.3s ease;
}

.chat-input:focus {
  border-color: #667eea;
}

.chat-input:disabled {
  background: #f1f5f9;
  cursor: not-allowed;
}

.send-button {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 50%;
  background: #667eea;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.send-button:hover:not(:disabled) {
  background: #5a67d8;
  transform: scale(1.05);
}

.send-button:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .message-wrapper {
    max-width: 85%;
  }

  .chat-input-container {
    padding: 0.75rem;
  }

  .role-name {
    font-size: 1.1rem;
  }
}
</style>
