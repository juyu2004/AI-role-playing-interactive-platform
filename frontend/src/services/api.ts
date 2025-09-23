import type {
  Role,
  ChatRequest,
  ChatResponse,
  BaseResponse,
  LoginRequest,
  LoginResponse,
  ASRResponse,
  TTSRequest,
  TTSResponse,
  Conversation,
  Message
} from '@/types/api'
import request from '@/utils/request'

class ApiService {
  /**
   * 健康检查
   * 用于验证后端服务是否正常运行
   */
  public static async healthCheck(): Promise<string> {
    try {
      // 直接使用 axios 实例发送请求，避免响应拦截器的处理
      const response = await request.getInstance().get('/api/health', {
        responseType: 'text'
      })
      return response.data
    } catch (error) {
      console.error('Health check failed:', error)
      throw error
    }
  }

  /**
   * 获取角色列表
   * 返回可用的AI角色信息列表
   */
  public static async getRoles(): Promise<BaseResponse<Role[]>> {
    return request.get<Role[]>('/api/roles')
  }

  /**
   * 获取角色详情
   * 根据角色ID获取特定角色的详细信息
   * @param roleId 角色唯一标识
   */
  public static async getRoleDetail(roleId: string): Promise<BaseResponse<Role>> {
    return request.get<Role>(`/api/roles/${roleId}`)
  }

  /**
   * 发送聊天消息
   * 向指定角色发送文本消息并获取回复
   * @param params 聊天请求参数，包含角色ID和消息文本
   */
  public static async sendChat(params: ChatRequest): Promise<BaseResponse<ChatResponse>> {
    return request.post<ChatResponse>('/api/chat', params)
  }

  /**
   * 流式聊天（占位实现）
   * 用于实现实时的对话流式响应
   * @param params 聊天请求参数
   * @param onMessage 收到消息片段时的回调函数
   */
  public static async streamChat(
    params: ChatRequest,
    onMessage: (text: string) => void
  ): Promise<void> {
    // 占位实现，实际项目中需要实现SSE连接
    // 这里仅做模拟，实际应使用EventSource或WebSocket
    return new Promise((resolve) => {
      setTimeout(() => {
        // 模拟流式响应
        onMessage('这是流式响应的一部分...')
        setTimeout(() => {
          onMessage('这是流式响应的另一部分。')
          resolve()
        }, 500)
      }, 1000)
    })
  }

  /**
   * 语音识别
   * 将音频文件转换为文本
   * @param formData 包含音频文件的表单数据
   */
  public static async recognizeSpeech(formData: FormData): Promise<BaseResponse<ASRResponse>> {
    return request.post<ASRResponse>('/api/asr', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }

  /**
   * 语音合成
   * 将文本转换为语音
   * @param params 文本转语音请求参数
   */
  public static async synthesizeSpeech(params: TTSRequest): Promise<BaseResponse<TTSResponse>> {
    return request.post<TTSResponse>('/api/tts', params)
  }

  /**
   * 用户登录
   * 用户身份验证并获取访问令牌
   * @param params 登录请求参数，包含邮箱和密码
   */
  public static async login(params: LoginRequest): Promise<BaseResponse<LoginResponse>> {
    return request.post<LoginResponse>('/api/auth/login', params)
  }

  /**
   * 列出会话
   * 获取用户的所有对话会话列表
   */
  public static async getConversations(): Promise<BaseResponse<Conversation[]>> {
    return request.get<Conversation[]>('/api/conversations')
  }

  /**
   * 列出消息
   * 获取指定会话的所有消息
   * @param conversationId 会话唯一标识
   */
  public static async getMessages(conversationId: string): Promise<BaseResponse<Message[]>> {
    return request.get<Message[]>(`/api/messages?conversationId=${conversationId}`)
  }

  /**
   * 上传音频文件
   * 用于上传语音消息或其他音频文件
   * @param formData 包含音频文件的表单数据
   */
  public static async uploadAudio(formData: FormData): Promise<BaseResponse<{ url: string }>> {
    return request.post<{ url: string }>('/api/upload/audio', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}

export default ApiService
