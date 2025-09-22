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

// 角色相关类型（根据你的后端调整）
export interface Role {
  id: string
  name: string
  category: string
  avatarUrl: string
  description: string
}

export interface ChatRequest {
  roleId: string
  text: string
}

// 完成ChatResponse接口定义
export interface ChatResponse {
  text: string
  timestamp: string
  audioUrl?: string | null
}

// 为消息列表添加音频URL支持
// interface ChatMessage {
//   id: string
//   text: string
//   isUser: boolean
//   timestamp: string
//   audioUrl?: string | null
// }
