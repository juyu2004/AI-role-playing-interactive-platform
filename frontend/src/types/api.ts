import type { AxiosRequestConfig } from "axios"

// API 响应基础类型
export interface BaseResponse<T = unknown> {
  code: number
  data: T
  message: string
  success: boolean
}

// 错误响应类型
export interface ErrorResponse {
  code: number
  message: string
  timestamp?: string
  path?: string
}

// 请求配置扩展
export interface RequestConfig extends AxiosRequestConfig {
  showError?: boolean
}

// 角色类型定义
export interface Role {
  id: string
  name: string
  category: string
  avatarUrl: string
  description: string
}

// 聊天请求类型
export interface ChatRequest {
  roleId: string
  text: string
}

// 聊天响应类型
export interface ChatResponse {
  text: string
  audioUrl?: string | null
}

// 登录请求类型
export interface LoginRequest {
  email: string
  password: string
}

// 登录响应类型
export interface LoginResponse {
  token: string
}

// 语音识别响应类型
export interface ASRResponse {
  text: string
}

// 语音合成请求类型
export interface TTSRequest {
  text: string
}

// 语音合成响应类型
export interface TTSResponse {
  audioUrl: string
}

// 消息类型定义（用于消息列表）
export interface Message {
  id: string
  text: string
  isUser: boolean
  timestamp: string
  audioUrl?: string | null
}

// 会话类型定义
export interface Conversation {
  id: string
  title: string
  roleId: string
  createTime: string
  updateTime: string
  lastMessage?: string
  lastMessageTime?: string
}
