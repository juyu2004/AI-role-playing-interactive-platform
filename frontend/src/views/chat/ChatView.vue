<template>
  <div class="chat-container">
    <!-- 左侧历史对话列表 -->
    <div class="chat-sidebar">
      <div class="sidebar-header">
        <h3>对话历史</h3>
        <button class="new-chat-button" @click="startNewChat" aria-label="新建对话">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14"></path>
          </svg>
        </button>
      </div>

      <div class="chat-history-list">
        <div
          v-for="chat in chatHistoryList"
          :key="chat.id"
          class="chat-history-item"
          :class="{ 'active': chat.id === currentChatId }"
          @click="switchToChat(chat.id)"
        >
          <div class="chat-history-avatar">
            <img
              :src="chat.roleAvatar || `https://picsum.photos/id/${parseInt(chat.roleId) % 50}/200/200`"
              :alt="chat.roleName"
            >
          </div>
          <div class="chat-history-content">
            <div class="chat-history-title">
              <span class="role-name">{{ chat.roleName }}</span>
              <span class="last-message-time">{{ formatLastMessageTime(chat.lastUpdated) }}</span>
            </div>
            <p class="last-message-preview">{{ getLastMessagePreview(chat) }}</p>
          </div>
          <button
            class="delete-chat-button"
            @click.stop="deleteChat(chat.id)"
            aria-label="删除对话"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3,6 5,6 21,6"></polyline>
              <path d="M19,6v14a2,2,0,0,1-2,2H7a2,2,0,0,1-2-2V6m3,0V4a2,2,0,0,1,2-2h4a2,2,0,0,1,2,2V6"></path>
              <line x1="10" y1="11" x2="10" y2="17"></line>
              <line x1="14" y1="11" x2="14" y2="17"></line>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-main">
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
      <!-- 修改聊天消息区域 -->
      <main class="chat-messages custom-scrollbar" ref="chatMessagesRef">
        <div v-if="loadingMessages" class="loading-messages">
          <div class="loading-spinner"></div>
          <p>加载对话历史...</p>
        </div>

        <div v-else-if="error" class="error-messages">
          <p>{{ error }}</p>
          // 修改前
          <!-- <button class="retry-button" @click="loadChatHistory">重试</button> -->

          // 修改后
          <button class="retry-button" @click="() => loadChatHistory()">重试</button>
        </div>

        <div v-else-if="messages.length === 0" class="empty-messages">
          <p>开始与{{ currentRole?.name }}的对话吧！</p>
        </div>

        <div v-for="message in messages" :key="message.id" class="message-wrapper" :class="{ 'user-message': message.isUser }">
          <!-- 角色头像 - 只在非用户消息时显示 -->
          <div v-if="!message.isUser" class="avatar-small">
            <div
              class="role-avatar"
              :style="{
                backgroundImage: `url(${currentRole?.avatarUrl || 'https://picsum.photos/id/' + (parseInt(roleId) % 50) + '/200/200'})`
              }"
              :aria-label="currentRole?.name"
            >
            </div>
          </div>
          <!-- 消息气泡 -->
          <div class="message-bubble" :class="{ 'user': message.isUser }">
            <p>{{ message.text }}</p>
            <!-- 音频消息播放控件 -->
            <div v-if="message.audioUrl" class="audio-message">
              <button
                class="play-button"
                @click="toggleAudioPlayback(message.id, message.audioUrl!)"
                :class="{ 'playing': currentlyPlayingAudio === message.id }"
              >
                <svg v-if="currentlyPlayingAudio !== message.id" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="6" y="4" width="4" height="16"></rect>
                  <rect x="14" y="4" width="4" height="16"></rect>
                </svg>
              </button>
              <div class="audio-progress">
                <div class="audio-timeline" @click="seekAudio(message.id, $event)">
                  <div class="audio-progress-bar" :style="{ width: getAudioProgress(message.id) + '%' }"></div>
                </div>
                <span class="audio-duration">{{ getAudioDurationFormatted(message.id) }}</span>
              </div>
            </div>
          </div>
          <!-- 用户头像 - 只在用户消息时显示 -->
          <div v-if="message.isUser" class="avatar-small">
            <div class="user-avatar">
              <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 12C14.7614 12 17 9.76142 17 7C17 4.23858 14.7614 2 12 2C9.23858 2 7 4.23858 7 7C7 9.76142 9.23858 12 12 12Z" stroke="#818cf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M20.42 14.82C21.6427 16.1828 22.2477 17.9812 22.0259 19.8705C21.8042 21.7599 20.7732 23.4527 19.1431 24.57C17.513 25.6873 15.4533 26.1322 13.4286 25.7857C11.4038 25.4392 9.55878 24.3161 8.25 22.75V14.82" stroke="#818cf8" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 正在输入提示 -->
        <div v-if="isTyping" class="typing-indicator">
          <div class="avatar-small">
            <div
              class="role-avatar"
              :style="{
                backgroundImage: `url(${currentRole?.avatarUrl || 'https://picsum.photos/id/' + (parseInt(roleId) % 50) + '/200/200'})`
              }"
              :aria-label="currentRole?.name"
            >
            </div>
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
          <!-- 语音录制按钮 - 改为点击切换模式 -->
          <button
            class="voice-button"
            @click="toggleRecording"
            :disabled="isTyping"
            :class="{ 'recording': isRecording }"
            aria-label="语音消息"
          >
            <svg v-if="!isRecording" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 14.3a1 1 0 0 1-1-1V9a1 1 0 0 1 2 0v4.3a1 1 0 0 1-1 1z" fill="currentColor"/>
              <path d="M9 11a1 1 0 0 1-1-1V6a1 1 0 0 1 2 0v4a1 1 0 0 1-1 1z" fill="currentColor"/>
              <path d="M15 11a1 1 0 0 1-1-1V6a1 1 0 0 1 2 0v4a1 1 0 0 1-1 1z" fill="currentColor"/>
              <circle cx="12" cy="19" r="2" fill="currentColor"/>
              <circle cx="12" cy="5" r="3" stroke="currentColor" stroke-width="2"/>
            </svg>
            <!-- 录制中显示停止图标 -->
            <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="8" y="5" width="8" height="14" rx="2" fill="currentColor"/>
            </svg>
          </button>
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
        <!-- 录制状态显示 -->
        <div v-if="isRecording" class="recording-indicator">
          <div class="recording-dot"></div>
          <span>正在录音，再次点击结束并发送</span>
          <!-- 添加录音时长显示 -->
          <span class="recording-duration">{{ recordingDuration }}</span>
        </div>
        <!-- 未录制状态提示 -->
        <!-- <div v-else class="recording-hint">
          <span>点击开始录音</span>
        </div> -->
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, onUnmounted} from 'vue'
import { useRouter, useRoute } from 'vue-router'
import ApiService from '../../services/api'
import type { Role, ChatRequest } from '../../types/api'

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

// 语音相关状态
const isRecording = ref(false)
const mediaRecorder = ref<MediaRecorder | null>(null)
const audioChunks = ref<Blob[]>([])
const audioElements = ref<Record<string, HTMLAudioElement>>({})
const audioProgress = ref<Record<string, number>>({})
const audioDurations = ref<Record<string, number>>({})
const currentlyPlayingAudio = ref<string | null>(null)

// 添加录音时长状态
const recordingDuration = ref('00:00')
const recordingStartTime = ref<number | null>(null)
const recordingTimer = ref<number | null>(null)

// 历史对话相关状态
const currentChatId = ref<string>(`chat_${roleId}_${Date.now()}`)
const chatHistoryList = ref<Array<{
  id: string
  roleId: string
  roleName: string
  roleAvatar: string | null
  lastUpdated: string
  messages: typeof messages.value
}>>([])

// 生成唯一的对话ID
const generateChatId = () => {
  return `chat_${roleId}_${Date.now()}`
}

// 本地存储键名
const getStorageKey = (chatId?: string) => {
  return chatId || currentChatId.value
}

// 保存对话历史到本地存储
const saveChatHistory = () => {
  try {
    // 保存当前对话
    localStorage.setItem(getStorageKey(), JSON.stringify({
      id: currentChatId.value,
      roleId,
      roleName: currentRole.value?.name || '',
      roleAvatar: currentRole.value?.avatarUrl || null,
      messages: messages.value,
      lastUpdated: new Date().toISOString()
    }))

    // 更新历史对话列表
    updateChatHistoryList()

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

    // 检查响应格式
    if (response && typeof response === 'object') {
      // 首先检查是否是BaseResponse格式（根据ApiService定义）
      if (response.success !== undefined && response.data) {
        currentRole.value = response.data
        console.log('成功加载角色信息(从BaseResponse):', currentRole.value)
      }
      // 然后检查是否是直接的Role对象
      else if ('id' in response && 'name' in response) {
        currentRole.value = response as unknown as Role
        console.log('成功加载角色信息:', currentRole.value)
      } else {
        throw new Error('无效的角色数据格式')
      }
    } else {
      throw new Error('无效的响应格式')
    }
  } catch (err) {
    console.error('加载角色信息失败:', err)
    // error.value = '加载角色信息失败'
    // 错误时使用模拟数据，但确保与选择的角色ID对应
    // 定义完整的模拟角色列表
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

    // 查找与传入的roleId相匹配的模拟角色
    const selectedMockRole = mockRoles.find(role => role.id === roleId);

    // 如果找到匹配的角色，则使用该角色；否则使用默认角色
    if (selectedMockRole) {
      currentRole.value = selectedMockRole;
    } else {
      // 退回到原来的取模逻辑作为备选方案
      currentRole.value = {
        id: roleId,
        name: ['夏洛克·福尔摩斯', '花木兰', '爱因斯坦', '莎士比亚', '居里夫人', '李白'][parseInt(roleId) % 6],
        category: ['文学角色', '历史人物', '科学家', '文学家', '科学家', '诗人'][parseInt(roleId) % 6],
        avatarUrl: `https://picsum.photos/id/${parseInt(roleId) % 50}/200/200`,
        description: '正在对话中的角色'
      }
    }
    console.log('使用模拟角色数据:', currentRole.value)
  }
}

// 加载对话历史
const loadChatHistory = async (chatId?: string) => {
  const targetChatId = chatId || currentChatId.value
  loadingMessages.value = true
  error.value = ''
  try {
    // 从本地存储加载对话历史
    const storedHistory = localStorage.getItem(getStorageKey(targetChatId))
    if (storedHistory) {
      try {
        const parsedHistory = JSON.parse(storedHistory)
        if (parsedHistory.messages && Array.isArray(parsedHistory.messages)) {
          messages.value = parsedHistory.messages
          currentChatId.value = targetChatId
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

// 加载所有历史对话列表
const loadChatHistoryList = () => {
  try {
    const historyList: typeof chatHistoryList.value = []

    // 遍历localStorage查找所有相关的对话历史
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(`chat_${roleId}_`)) {
        try {
          const storedHistory = localStorage.getItem(key)
          if (storedHistory) {
            const parsedHistory = JSON.parse(storedHistory)
            if (parsedHistory.roleId === roleId) {
              historyList.push({
                id: parsedHistory.id || key,
                roleId: parsedHistory.roleId,
                roleName: parsedHistory.roleName || '',
                roleAvatar: parsedHistory.roleAvatar || null,
                lastUpdated: parsedHistory.lastUpdated || new Date().toISOString(),
                messages: parsedHistory.messages || []
              })
            }
          }
        } catch (error) {
          console.error('解析历史对话失败:', error)
        }
      }
    }

    // 按最后更新时间排序
    historyList.sort((a, b) => {
      return new Date(b.lastUpdated).getTime() - new Date(a.lastUpdated).getTime()
    })

    chatHistoryList.value = historyList
    console.log('加载历史对话列表成功:', historyList.length, '条记录')
  } catch (err) {
    console.error('加载历史对话列表失败:', err)
  }
}

// 更新历史对话列表
const updateChatHistoryList = () => {
  loadChatHistoryList()
}

// 切换到指定对话
const switchToChat = (chatId: string) => {
  if (chatId !== currentChatId.value) {
    loadChatHistory(chatId)
  }
}

// 开始新对话
const startNewChat = () => {
  currentChatId.value = generateChatId()
  messages.value = []
  // 自动保存新对话记录
  saveChatHistory()
}

// 删除对话
const deleteChat = (chatId: string) => {
  if (confirm('确定要删除这条对话历史吗？删除后不可恢复。')) {
    try {
      localStorage.removeItem(getStorageKey(chatId))
      updateChatHistoryList()

      // 如果删除的是当前对话，创建新对话
      if (chatId === currentChatId.value) {
        startNewChat()
      }

      console.log('对话历史已删除')
    } catch (err) {
      console.error('删除对话历史失败:', err)
    }
  }
}

// 清除当前对话历史
const clearChatHistory = () => {
  if (confirm('确定要清除当前对话历史吗？清除后不可恢复。')) {
    try {
      messages.value = []
      saveChatHistory()
      console.log('当前对话历史已清除')
    } catch (err) {
      console.error('清除对话历史失败:', err)
    }
  }
}

// 获取最后一条消息预览
const getLastMessagePreview = (chat: typeof chatHistoryList.value[0]): string => {
  if (!chat.messages || chat.messages.length === 0) {
    return '新对话'
  }

  const lastMessage = chat.messages[chat.messages.length - 1]
  let preview = lastMessage.text

  // 如果是语音消息，显示特殊标识
  if (preview === '[语音消息]') {
    preview = '[语音]'
  }

  // 限制预览长度
  if (preview.length > 40) {
    preview = preview.substring(0, 40) + '...'
  }

  return preview
}

// 格式化最后消息时间
const formatLastMessageTime = (timestamp: string): string => {
  const date = new Date(timestamp)
  const now = new Date()
  const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))

  if (diffInMinutes < 1) {
    return '刚刚'
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  } else if (diffInMinutes < 1440) {
    const hours = Math.floor(diffInMinutes / 60)
    return `${hours}小时前`
  } else {
    const days = Math.floor(diffInMinutes / 1440)
    return `${days}天前`
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
    console.log('API响应:', response)

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
      // 实现音频播放
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

// 保留第687行开始的第二个更加完整的startRecording函数实现
// 开始录音
const startRecording = async () => {
  try {
    // 请求用户授权使用麦克风
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })

    // 创建MediaRecorder实例
    const mediaRecorderInstance = new MediaRecorder(stream)
    mediaRecorder.value = mediaRecorderInstance
    audioChunks.value = []

    // 监听数据可用事件
    mediaRecorderInstance.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.value.push(event.data)
      }
    }

    // 监听录制停止事件
    mediaRecorderInstance.onstop = async () => {
      // 停止所有音轨
      stream.getTracks().forEach(track => track.stop())

      // 停止计时
      if (recordingTimer.value) {
        clearInterval(recordingTimer.value)
        recordingTimer.value = null
      }

      // 重置录制时间
      recordingDuration.value = '00:00'
      recordingStartTime.value = null

      // 创建音频Blob
      const audioBlob = new Blob(audioChunks.value, { type: 'audio/wav' })

      // 上传音频文件
      uploadAndSendAudio(audioBlob)
    }

    // 开始录制
    mediaRecorderInstance.start()
    isRecording.value = true

    // 记录开始时间并启动计时器
    recordingStartTime.value = Date.now()
    updateRecordingDuration()
    recordingTimer.value = window.setInterval(updateRecordingDuration, 1000)

    console.log('开始录音')
  } catch (err) {
    console.error('开始录音失败:', err)
    alert('无法访问麦克风，请确保已授予权限')
  }
}

// 上传并发送音频
const uploadAndSendAudio = async (audioBlob: Blob) => {
  try {
    // 创建FormData对象
    const formData = new FormData()
    formData.append('audio', audioBlob, `recording-${Date.now()}.wav`)
    // 显示正在输入状态
    isTyping.value = true
    // 上传音频
    const uploadResponse = await ApiService.uploadAudio(formData)
    let audioUrl: string | null = null

    if (uploadResponse && typeof uploadResponse === 'object') {
      if ('url' in uploadResponse) {
        audioUrl = uploadResponse.url as string
      } else if (uploadResponse.data && uploadResponse.data.url) {
        audioUrl = uploadResponse.data.url
      }
    }

    // 构建聊天请求（假设API支持音频输入）
    const chatRequest: ChatRequest = {
      roleId,
      text: '[语音消息]'
    }

    // 调用API发送消息
    const response = await ApiService.sendChat(chatRequest)
    console.log('发送音频消息响应:', response)

    // 处理API响应
    let replyText = '抱歉，我暂时无法回答这个问题。'
    let replyAudioUrl: string | null = null

    if (response && typeof response === 'object') {
      if ('text' in response) {
        replyText = response.text as string
        replyAudioUrl = 'audioUrl' in response ? (response as Record<string, unknown>).audioUrl as string | null || null : null
      } else if (response.data && typeof response.data === 'object') {
        replyText = response.data.text || replyText
        replyAudioUrl = response.data.audioUrl || null
      }
    }

    // 添加用户语音消息到列表
    const userMessage = {
      id: `user-${Date.now()}`,
      text: '[语音消息]',
      isUser: true,
      timestamp: new Date().toISOString(),
      audioUrl: audioUrl
    }
    messages.value.push(userMessage)

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

  } catch (err) {
    console.error('发送音频消息失败:', err)
    // 添加失败的语音消息
    const userMessage = {
      id: `user-${Date.now()}`,
      text: '[语音消息发送失败]',
      isUser: true,
      timestamp: new Date().toISOString(),
      audioUrl: null
    }
    messages.value.push(userMessage)

    // 使用模拟回复
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

// 切换音频播放
const toggleAudioPlayback = async (messageId: string, audioUrl: string) => {
  try {
    // 如果正在播放其他音频，先停止
    if (currentlyPlayingAudio.value && currentlyPlayingAudio.value !== messageId) {
      stopAudioPlayback(currentlyPlayingAudio.value)
    }

    // 检查是否已有音频元素
    if (!audioElements.value[messageId]) {
      // 创建新的音频元素
      const audio = new Audio(audioUrl)
      audioElements.value[messageId] = audio

      // 监听音频元数据加载完成事件
      audio.addEventListener('loadedmetadata', () => {
        audioDurations.value[messageId] = audio.duration
        audioProgress.value[messageId] = 0
      })

      // 监听时间更新事件
      audio.addEventListener('timeupdate', () => {
        if (audio.duration > 0) {
          audioProgress.value[messageId] = (audio.currentTime / audio.duration) * 100
        }
      })

      // 监听播放结束事件
      audio.addEventListener('ended', () => {
        audioProgress.value[messageId] = 100
        currentlyPlayingAudio.value = null
      })

      // 监听错误事件
      audio.addEventListener('error', (err) => {
        console.error('音频播放错误:', err)
        alert('音频播放失败')
        currentlyPlayingAudio.value = null
      })
    }

    const audio = audioElements.value[messageId]

    // 切换播放/暂停状态
    if (currentlyPlayingAudio.value === messageId) {
      // 暂停播放
      audio.pause()
      currentlyPlayingAudio.value = null
    } else {
      // 开始播放
      await audio.play()
      currentlyPlayingAudio.value = messageId
    }
  } catch (err) {
    console.error('音频播放控制失败:', err)
    alert('音频播放控制失败')
    currentlyPlayingAudio.value = null
  }
}

// 停止音频播放
const stopAudioPlayback = (messageId: string) => {
  if (audioElements.value[messageId]) {
    audioElements.value[messageId].pause()
    if (currentlyPlayingAudio.value === messageId) {
      currentlyPlayingAudio.value = null
    }
  }
}

// 音频进度条拖动
const seekAudio = (messageId: string, event: MouseEvent) => {
  if (!audioElements.value[messageId]) return

  const audio = audioElements.value[messageId]
  const timeline = event.currentTarget as HTMLElement
  const rect = timeline.getBoundingClientRect()
  const pos = (event.clientX - rect.left) / rect.width

  if (audio.duration > 0) {
    audio.currentTime = pos * audio.duration
    audioProgress.value[messageId] = pos * 100
  }
}

// 获取音频进度
const getAudioProgress = (messageId: string): number => {
  return audioProgress.value[messageId] || 0
}

// 获取格式化的音频时长
const getAudioDurationFormatted = (messageId: string): string => {
  const duration = audioDurations.value[messageId] || 0
  const minutes = Math.floor(duration / 60)
  const seconds = Math.floor(duration % 60)
  return `${minutes}:${seconds < 10 ? '0' + seconds : seconds}`
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
  loadChatHistoryList()
})

// 切换录音状态 - 改为点击切换模式
const toggleRecording = async () => {
  if (isTyping.value) return

  if (isRecording.value) {
    // 如果正在录音，则停止录音
    stopRecording()
  } else {
    // 如果未录音，则开始录音
    await startRecording()
  }
}

// 停止录音
const stopRecording = () => {
  if (!isRecording.value || !mediaRecorder.value) return

  // 确保状态先更新
  isRecording.value = false

  // 停止录制
  try {
    if (mediaRecorder.value.state !== 'inactive') {
      mediaRecorder.value.stop()
      console.log('停止录音')
    }
  } catch (error) {
    console.error('停止录音时出错:', error)
  }
}

// 更新录音时长显示
const updateRecordingDuration = () => {
  if (!recordingStartTime.value) {
    recordingDuration.value = '00:00'
    return
  }

  const durationInSeconds = Math.floor((Date.now() - recordingStartTime.value) / 1000)
  const minutes = Math.floor(durationInSeconds / 60)
  const seconds = durationInSeconds % 60

  recordingDuration.value = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}

// 组件卸载时清理资源
onUnmounted(() => {
  // 停止所有音频播放
  Object.keys(audioElements.value).forEach(messageId => {
    if (audioElements.value[messageId]) {
      audioElements.value[messageId].pause()
      audioElements.value[messageId].src = ''
    }
  })

  // 清理录音资源
  if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
    mediaRecorder.value.stop()
  }

  // 清理计时器
  if (recordingTimer.value) {
    clearInterval(recordingTimer.value)
    recordingTimer.value = null
  }
})
</script>

<style scoped src="../chat/chat.css"></style>

