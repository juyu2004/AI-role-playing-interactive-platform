<template>
  <div class="register-container">
    <!-- 返回按钮 -->
    <button class="back-button" @click="goBack">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M19 12H5M5 12L12 5M5 12L12 19" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <!-- 注册表单卡片 -->
    <div class="register-card">
      <!-- Logo 和标题 -->
      <div class="register-header">
        <div class="logo">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="#667eea"/>
            <path d="M2 17L12 22L22 17" stroke="#667eea" stroke-width="2" stroke-linecap="round"/>
            <path d="M2 12L12 17L22 12" stroke="#667eea" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <h2 class="register-title">创建账号</h2>
        <p class="register-subtitle">加入我们，开始与AI角色对话的奇妙旅程</p>
      </div>

      <!-- 注册表单 -->
      <form class="register-form" @submit.prevent="handleRegister">
        <!-- 用户名输入框 -->
        <div class="form-group">
          <label for="username" class="form-label">用户名</label>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="#9CA3AF" stroke-width="2"/>
              <path d="M12 16V12M12 8H12.01" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <input
              type="text"
              id="username"
              v-model="username"
              placeholder="请设置用户名"
              class="form-input"
              required
              minlength="3"
              maxlength="20"
            />
          </div>
          <p v-if="usernameError" class="error-message">{{ usernameError }}</p>
        </div>

        <!-- 邮箱输入框 -->
        <div class="form-group">
          <label for="email" class="form-label">邮箱</label>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 4H20C20.5304 4 21.0391 4.21071 21.4142 4.58579C21.7893 4.96086 22 5.46957 22 6V18C22 18.5304 21.7893 19.0391 21.4142 19.4142C21.0391 19.7893 20.5304 20 20 20H4C3.46957 20 2.96086 19.7893 2.58579 19.4142C2.21071 19.0391 2 18.5304 2 18V6C2 5.46957 2.21071 4.96086 2.58579 4.58579C2.96086 4.21071 3.46957 4 4 4Z" stroke="#9CA3AF" stroke-width="2"/>
              <path d="M22 6L12 13L2 6" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <input
              type="email"
              id="email"
              v-model="email"
              placeholder="请输入邮箱"
              class="form-input"
              required
            />
          </div>
          <p v-if="emailError" class="error-message">{{ emailError }}</p>
        </div>

        <!-- 密码输入框 -->
        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 15V17M12 7V9M5 12H19M6.17157 6.17157L4.75736 7.58579M19.2426 16.4142L17.8284 17.8284M19.2426 6.17157L17.8284 7.58579M6.17157 17.8284L4.75736 16.4142M17 12C17 14.7614 14.7614 17 12 17C9.23858 17 7 14.7614 7 12C7 9.23858 9.23858 7 12 7C14.7614 7 17 9.23858 17 12Z" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <input
              :type="showPassword ? 'text' : 'password'"
              id="password"
              v-model="password"
              placeholder="请设置密码"
              class="form-input"
              required
              minlength="8"
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
          <p v-if="passwordError" class="error-message">{{ passwordError }}</p>
          <p class="password-hint">密码至少8位，包含字母和数字</p>
        </div>

        <!-- 确认密码输入框 -->
        <div class="form-group">
          <label for="confirmPassword" class="form-label">确认密码</label>
          <div class="input-container">
            <svg class="input-icon" width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 15V17M12 7V9M5 12H19M6.17157 6.17157L4.75736 7.58579M19.2426 16.4142L17.8284 17.8284M19.2426 6.17157L17.8284 7.58579M6.17157 17.8284L4.75736 16.4142M17 12C17 14.7614 14.7614 17 12 17C9.23858 17 7 14.7614 7 12C7 9.23858 9.23858 7 12 7C14.7614 7 17 9.23858 17 12Z" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <input
              :type="showPassword ? 'text' : 'password'"
              id="confirmPassword"
              v-model="confirmPassword"
              placeholder="请再次输入密码"
              class="form-input"
              required
            />
          </div>
          <p v-if="confirmPasswordError" class="error-message">{{ confirmPasswordError }}</p>
        </div>

        <!-- 用户协议 -->
        <div class="form-group">
          <label class="agreement-container">
            <input type="checkbox" v-model="agreedToTerms" required/>
            <span class="checkmark"></span>
            <span class="agreement-text">
              我已阅读并同意 <a href="#" class="agreement-link">服务条款</a> 和 <a href="#" class="agreement-link">隐私政策</a>
            </span>
          </label>
        </div>

        <!-- 注册按钮 -->
        <button type="submit" class="register-button" :disabled="isLoading || !isFormValid">
          <span v-if="!isLoading">立即注册</span>
          <span v-else class="loading-spinner"></span>
        </button>

        <!-- 登录提示 -->
        <div class="login-prompt">
          <span>已有账号？</span>
          <a href="#" class="login-link" @click.prevent="goToLogin">立即登录</a>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 表单数据
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const agreedToTerms = ref(false)
const isLoading = ref(false)

// 错误信息
const usernameError = ref('')
const emailError = ref('')
const passwordError = ref('')
const confirmPasswordError = ref('')

// 表单验证
const validateUsername = () => {
  if (!username.value.trim()) {
    usernameError.value = '请输入用户名'
    return false
  }
  if (username.value.length < 3 || username.value.length > 20) {
    usernameError.value = '用户名长度应在3-20个字符之间'
    return false
  }
  usernameError.value = ''
  return true
}

const validateEmail = () => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!email.value.trim()) {
    emailError.value = '请输入邮箱'
    return false
  }
  if (!emailRegex.test(email.value)) {
    emailError.value = '请输入有效的邮箱地址'
    return false
  }
  emailError.value = ''
  return true
}

const validatePassword = () => {
  const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/
  if (!password.value) {
    passwordError.value = '请设置密码'
    return false
  }
  if (password.value.length < 8) {
    passwordError.value = '密码长度至少为8位'
    return false
  }
  if (!passwordRegex.test(password.value)) {
    passwordError.value = '密码必须包含字母和数字'
    return false
  }
  passwordError.value = ''
  return true
}

const validateConfirmPassword = () => {
  if (!confirmPassword.value) {
    confirmPasswordError.value = '请确认密码'
    return false
  }
  if (confirmPassword.value !== password.value) {
    confirmPasswordError.value = '两次输入的密码不一致'
    return false
  }
  confirmPasswordError.value = ''
  return true
}

// 监听输入变化进行验证
watch(username, validateUsername)
watch(email, validateEmail)
watch(password, () => {
  validatePassword()
  validateConfirmPassword()
})
watch(confirmPassword, validateConfirmPassword)

// 表单是否有效
const isFormValid = computed(() => {
  return (
    validateUsername() &&
    validateEmail() &&
    validatePassword() &&
    validateConfirmPassword() &&
    agreedToTerms.value
  )
})

// 处理注册
const handleRegister = async () => {
  if (!isFormValid.value) {
    // 触发所有验证
    validateUsername()
    validateEmail()
    validatePassword()
    validateConfirmPassword()
    return
  }

  try {
    isLoading.value = true

    // 模拟注册请求
    await new Promise(resolve => setTimeout(resolve, 2000))

    // 注册成功后，存储用户信息并跳转
    localStorage.setItem('userToken', 'mock-token')
    localStorage.setItem('userAvatar', 'https://picsum.photos/id/64/40/40')

    // 跳转到角色选择页面
    router.push({ name: 'roleSelect' })
  } catch (error) {
    console.error('注册失败:', error)
    // 实际项目中这里应该显示错误提示
  } finally {
    isLoading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 跳转到登录页面
const goToLogin = () => {
  router.push({ name: 'login' })
}
</script>

<style scoped>
/* 注册页面容器 */
.register-container {
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

/* 注册卡片 */
.register-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  width: 100%;
  max-width: 360px;
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
  animation: fadeIn 0.5s ease-in-out;
}

/* 注册头部 */
.register-header {
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

.register-title {
  font-size: 1.4rem;
  color: #1a202c;
  margin-bottom: 0.2rem;
  font-weight: 700;
}

.register-subtitle {
  font-size: 0.8rem;
  color: #718096;
}

/* 表单组 */
.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  font-size: 0.8rem;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.3rem;
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

/* 错误消息 */
.error-message {
  color: #e53e3e;
  font-size: 0.7rem;
  margin-top: 0.3rem;
}

/* 密码提示 */
.password-hint {
  color: #718096;
  font-size: 0.7rem;
  margin-top: 0.3rem;
}

/* 用户协议 */
.agreement-container {
  display: flex;
  align-items: flex-start;
  cursor: pointer;
  user-select: none;
  font-size: 0.8rem;
  color: #4a5568;
}

/* 注册按钮 */
.register-button {
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

/* 加载动画 */
.loading-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

/* 登录提示 */
.login-prompt {
  text-align: center;
  font-size: 0.8rem;
  color: #718096;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-card {
    padding: 1.2rem;
    max-width: 320px;
  }

  .back-button {
    top: 0.8rem;
    left: 0.8rem;
    width: 32px;
    height: 32px;
  }

  .register-title {
    font-size: 1.3rem;
  }
}
</style>
