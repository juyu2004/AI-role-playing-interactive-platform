<template>
  <div class="login-container">
    <!-- 返回按钮 -->
    <button class="back-button" @click="goBack">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M19 12H5M5 12L12 5M5 12L12 19" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <!-- 登录表单卡片 -->
    <div class="login-card">
      <!-- Logo 和标题 -->
      <div class="login-header">
        <div class="logo">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="#667eea"/>
            <path d="M2 17L12 22L22 17" stroke="#667eea" stroke-width="2" stroke-linecap="round"/>
            <path d="M2 12L12 17L22 12" stroke="#667eea" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <h2 class="login-title">欢迎回来</h2>
        <p class="login-subtitle">登录您的账号继续对话</p>
      </div>

      <!-- 登录表单 -->
      <form class="login-form" @submit.prevent="handleLogin">
        <!-- 用户名/邮箱输入框 -->
        <div class="form-group">
          <label for="username" class="form-label">用户名或邮箱</label>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="#9CA3AF" stroke-width="2"/>
              <path d="M12 16V12M12 8H12.01" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <input
              type="text"
              id="username"
              v-model="username"
              placeholder="请输入用户名或邮箱"
              class="form-input"
              required
            />
          </div>
        </div>

        <!-- 密码输入框 -->
        <div class="form-group">
          <div class="label-container">
            <label for="password" class="form-label">密码</label>
            <a href="#" class="forgot-password">忘记密码?</a>
          </div>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 15V17M12 7V9M5 12H19M6.17157 6.17157L4.75736 7.58579M19.2426 16.4142L17.8284 17.8284M19.2426 6.17157L17.8284 7.58579M6.17157 17.8284L4.75736 16.4142M17 12C17 14.7614 14.7614 17 12 17C9.23858 17 7 14.7614 7 12C7 9.23858 9.23858 7 12 7C14.7614 7 17 9.23858 17 12Z" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <input
              type="password"
              id="password"
              v-model="password"
              placeholder="请输入密码"
              class="form-input"
              required
            />
            <button
              type="button"
              class="toggle-password"
              @click="showPassword = !showPassword"
            >
              <svg v-if="showPassword" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M3 12C3 12 6 7 12 7C18 7 21 12 21 12C21 12 18 17 12 17C6 17 3 12 3 12Z" stroke="#9CA3AF" stroke-width="2"/>
                <path d="M12 13C12.5523 13 13 12.5523 13 12C13 11.4477 12.5523 11 12 11C11.4477 11 11 11.4477 11 12C11 12.5523 11.4477 13 12 13Z" fill="#9CA3AF"/>
              </svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M18.8571 5.14288C16.5917 3.17031 13.8599 1.86328 10.9672 1.38028C9.53066 1.12774 8.06738 1.01017 6.6041 1.04347C5.14082 1.07676 3.68935 1.25963 2.29354 1.58602C2.13392 1.64336 2 1.79659 2 1.95711V22.0429C2 22.2034 2.13392 22.3566 2.29354 22.414C3.68935 22.7404 5.14082 22.9232 6.6041 22.9565C8.06738 22.9898 9.53066 22.8723 10.9672 22.6197C13.8599 22.1367 16.5917 20.8297 18.8571 18.8571C20.8297 16.5917 22.1367 13.8599 22.6197 10.9672C22.8723 9.53066 22.9898 8.06738 22.9565 6.6041C22.9232 5.14082 22.7404 3.68935 22.414 2.29354C22.3566 2.13392 22.2034 2 22.0429 2H1.95714C1.79659 2 1.64336 2.13392 1.58602 2.29354C1.25963 3.68935 1.07676 5.14082 1.04347 6.6041C1.01017 8.06738 1.12774 9.53066 1.38028 10.9672C1.86328 13.8599 3.17031 16.5917 5.14288 18.8571L18.8571 5.14288Z" stroke="#9CA3AF" stroke-width="2"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>

        <!-- 登录按钮 -->
        <button type="submit" class="login-button" :disabled="isLoading">
          <span v-if="!isLoading">登录</span>
          <span v-else class="loading-spinner"></span>
        </button>

        <!-- 其他登录方式 -->
        <div class="other-login-options">
          <div class="divider"><span>或通过以下方式登录</span></div>
          <div class="social-login-buttons">
            <button type="button" class="social-button google">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M22.54 6.42A8.48 8.48 0 0 1 12 22.54A8.48 8.48 0 0 1 1.46 6.42C3.79 10.28 7.46 13 12 13C16.54 13 20.21 10.28 22.54 6.42Z" fill="#EA4335"/>
                <path d="M12 2.27A10.2 10.2 0 0 0 1.46 6.42A10.2 10.2 0 0 0 12 22.54A10.2 10.2 0 0 0 22.54 6.42" stroke="#FBBC05" stroke-width="2"/>
                <path d="M17.85 10.74L12 15.98L6.15 10.74" stroke="#4285F4" stroke-width="2"/>
                <path d="M12 2.27V15.98" stroke="#34A853" stroke-width="2"/>
              </svg>
              <span>Google</span>
            </button>
            <button type="button" class="social-button github">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M12 0C5.37 0 0 5.37 0 12C0 16.42 2.87 20.17 6.84 21.5C7.34 21.58 7.5 21.27 7.5 21V18.75C7.5 18.34 7.34 17.5 6.5 17.5C5.84 17.5 5.03 17.67 4.79 17.75C4.33 15.94 2.79 15.5 2.79 15.5C2.25 14.84 3 14.75 3 14.75C3.75 14.92 4.25 15.75 4.25 15.75C5 17.64 6.75 17.13 7.5 16.75C7.67 16.14 7.84 15.64 8 15.19C5.75 14.92 3.5 14.03 3.5 10.42C3.5 9.17 3.92 8.17 4.58 7.34C4.5 7.17 4.25 6.34 4.75 5.34C4.75 5.34 5.5 5.17 7.5 6.5C8.84 6.34 10.17 6.17 11.5 6.17C12.84 6.17 14.17 6.34 15.5 6.5C17.5 5.17 18.25 5.34 18.25 5.34C18.75 6.34 18.5 7.17 18.42 7.34C19.08 8.17 19.5 9.17 19.5 10.42C19.5 14.17 17.25 14.84 14.92 15.25C15.17 15.75 15.42 16.34 15.42 17.25V21C15.42 21.27 15.58 21.58 16.08 21.5C20.03 20.17 22.92 16.42 22.92 12C22.92 5.37 17.55 0 12 0Z" fill="white"/>
              </svg>
              <span>GitHub</span>
            </button>
          </div>
        </div>

        <!-- 注册提示 -->
        <div class="register-prompt">
          <span>还没有账号？</span>
          <a href="#" class="register-link" @click.prevent="goToRegister">立即注册</a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ApiService from '@/services/api'

const router = useRouter()

// 表单数据
const username = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')

// 处理登录
const handleLogin = async () => {
  try {
    // 重置错误信息
    errorMessage.value = ''
    isLoading.value = true

    // 表单验证
    if (!username.value || !password.value) {
      errorMessage.value = '请填写完整的登录信息'
      return
    }

    // 使用API服务进行登录请求
    const response = await ApiService.login({
      email: username.value,
      password: password.value
    })
    // 这里假设response已经是正确的登录响应数据
    if (response && response.data && response.data.token) {
      // 登录成功，存储token，设置一个默认头像（占位）
      localStorage.setItem('token', response.data?.token?.toString() || '')
      localStorage.setItem('userAvatar', 'https://picsum.photos/id/64/40/40')
      router.push({ name: 'roleSelect' })
    } else {
      throw new Error('登录失败，无效的响应数据')
    }
  } catch (error) {
    console.error('登录失败:', error)
    // 显示错误提示
    errorMessage.value = error instanceof Error ? error.message : '登录失败，请重试'
  } finally {
    isLoading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 跳转到注册页面
const goToRegister = () => {
  router.push({ name: 'register' })
}
</script>

<style scoped>
/* 登录页面容器 */
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  position: relative;
}

/* 返回按钮 */
.back-button {
  position: absolute;
  top: 1.2rem;
  left: 1.2rem;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-button svg {
  width: 18px;
  height: 18px;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(-5px);
}

/* 登录卡片 */
.login-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  width: 100%;
  max-width: 360px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 登录头部 */
.login-header {
  text-align: center;
  margin-bottom: 1.5rem;
}

.logo {
  width: 50px;
  height: 50px;
  margin: 0 auto 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
}

.logo svg {
  width: 30px;
  height: 30px;
}

.login-title {
  font-size: 1.4rem;
  color: #1a202c;
  margin-bottom: 0.2rem;
  font-weight: 700;
}

.login-subtitle {
  font-size: 0.8rem;
  color: #718096;
}

/* 表单组 */
.form-group {
  margin-bottom: 1rem;
}

.label-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.3rem;
}

.form-label {
  display: block;
  font-size: 0.8rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.3rem;
}

.forgot-password {
  font-size: 0.75rem;
  color: #667eea;
  text-decoration: none;
  transition: color 0.3s ease;
}

.forgot-password:hover {
  color: #764ba2;
}

/* 输入容器 */
.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 0.7rem;
  color: #9CA3AF;
  width: 14px;
  height: 14px;
}

/* 表单输入框 */
.form-input {
  width: 100%;
  padding: 0.75rem 0.75rem 0.75rem 2.2rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.85rem;
  transition: all 0.3s ease;
  background: white;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
  color: #a0aec0;
}

/* 错误消息样式 */
.error-message {
  color: #e53e3e;
  font-size: 0.75rem;
  margin-bottom: 1rem;
  padding: 0.5rem;
  background: #fff5f5;
  border: 1px solid #fed7d7;
  border-radius: 6px;
}

/* 切换密码可见性按钮 */
.toggle-password {
  position: absolute;
  right: 0.75rem;
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.3rem;
}

.toggle-password svg {
  width: 14px;
  height: 14px;
}

/* 登录按钮 */
.login-button {
  width: 100%;
  padding: 0.75rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
}

.login-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 加载动画 */
.loading-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 其他登录方式 */
.other-login-options {
  margin-bottom: 1rem;
}

.divider {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e2e8f0;
}

.divider span {
  padding: 0 0.8rem;
  font-size: 0.75rem;
  color: #718096;
}

/* 社交登录按钮 */
.social-login-buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.8rem;
}

.social-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.65rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.8rem;
  font-weight: 500;
}

.social-button:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
}

.social-button svg {
  width: 14px;
  height: 14px;
}

.social-button.google {
  color: #ea4335;
}

.social-button.github {
  color: #333;
}

/* 注册提示 */
.register-prompt {
  text-align: center;
  font-size: 0.8rem;
  color: #718096;
}

.register-link {
  color: #667eea;
  font-weight: 600;
  text-decoration: none;
  transition: color 0.3s ease;
}

.register-link:hover {
  color: #764ba2;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-card {
    padding: 1.2rem;
    max-width: 320px;
  }

  .back-button {
    top: 0.8rem;
    left: 0.8rem;
    width: 32px;
    height: 32px;
  }

  .login-title {
    font-size: 1.3rem;
  }

  .social-login-buttons {
    grid-template-columns: 1fr;
  }
}
</style>
