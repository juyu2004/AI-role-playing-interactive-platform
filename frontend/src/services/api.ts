import type { Role, ChatRequest, ChatResponse, BaseResponse } from '@/types/api'
import request from '@/utils/request'

class ApiService {
  // 获取角色列表
  public static async getRoles(): Promise<BaseResponse<Role[]>> {
    return request.get<Role[]>('/api/roles')
  }

  // 获取角色详情
  public static async getRoleDetail(roleId: string): Promise<BaseResponse<Role>> {
    return request.get<Role>(`/api/roles/${roleId}`)
  }

  // 发送聊天消息
  public static async sendChat(params: ChatRequest): Promise<BaseResponse<ChatResponse>> {
    return request.post<ChatResponse>('/api/chat', params)
  }

  // 上传音频文件
  public static async uploadAudio(formData: FormData): Promise<BaseResponse<{ url: string }>> {
    return request.post<{ url: string }>('/api/upload/audio', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

export default ApiService
