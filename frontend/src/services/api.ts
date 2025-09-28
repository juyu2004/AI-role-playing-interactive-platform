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
  Message,
  LoginApiResponse,
  User,
  UpdateAvatarRequest,
  UpdateBioRequest,
  FriendsResponse,
  AddFriendRequest,
  ReassignUserDataRequest
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

  // /**
  //  * 流式聊天（占位实现）
  //  * 用于实现实时的对话流式响应
  //  * @param params 聊天请求参数
  //  * @param onMessage 收到消息片段时的回调函数
  //  */
  // public static async streamChat(
  //   params: ChatRequest,
  //   onMessage: (text: string) => void
  // ): Promise<void> {
  //   // 占位实现，实际项目中需要实现SSE连接
  //   // 这里仅做模拟，实际应使用EventSource或WebSocket
  //   return new Promise((resolve) => {
  //     setTimeout(() => {
  //       // 模拟流式响应
  //       onMessage('这是流式响应的一部分...')
  //       setTimeout(() => {
  //         onMessage('这是流式响应的另一部分。')
  //         resolve()
  //       }, 500)
  //     }, 1000)
  //   })
  // }

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
  public static async login(params: LoginRequest): Promise<LoginApiResponse> {
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

  /**
   * 用户注册
   * 创建新用户账户并获取访问令牌
   * @param params 注册请求参数，包含邮箱和密码
   */
  public static async register(params: LoginRequest): Promise<LoginApiResponse> {
    return request.post<LoginResponse>('/api/auth/register', params)
  }

  /**
   * 获取当前用户信息
   * 需要JWT认证
   */
  public static async getCurrentUser(): Promise<BaseResponse<User>> {
    return request.get<User>('/api/user/me')
  }

  /**
   * 更新用户头像
   * 需要JWT认证
   * @param params 更新头像请求参数
   */
  public static async updateAvatar(params: UpdateAvatarRequest): Promise<BaseResponse<{ status: string }>> {
    return request.post<{ status: string }>('/api/user/avatar', params)
  }

  /**
   * 更新用户简介
   * 需要JWT认证
   * @param params 更新简介请求参数
   */
  public static async updateBio(params: UpdateBioRequest): Promise<BaseResponse<{ status: string }>> {
    return request.post<{ status: string }>('/api/user/bio', params)
  }

  /**
   * 获取好友列表
   * 需要JWT认证
   */
  public static async getFriends(): Promise<BaseResponse<FriendsResponse>> {
    return request.get<FriendsResponse>('/api/user/friends')
  }

  /**
   * 添加好友
   * 需要JWT认证
   * @param params 添加好友请求参数
   */
  public static async addFriend(params: AddFriendRequest): Promise<BaseResponse<{ status: string }>> {
    return request.post<{ status: string }>('/api/user/friends/add', params)
  }

  /**
   * 删除好友
   * 需要JWT认证
   * @param friendUserId 好友用户ID
   */
  public static async removeFriend(friendUserId: string): Promise<BaseResponse<{ status: string }>> {
    return request.delete<{ status: string }>(`/api/user/friends/${friendUserId}`)
  }

  /**
   * 重新分配用户数据
   * 需要JWT认证
   * @param params 重新分配请求参数
   */
  public static async reassignUserData(params: ReassignUserDataRequest): Promise<BaseResponse<{ status: string }>> {
    return request.post<{ status: string }>('/api/conversations/reassign', params)
  }

  /**
   * 实现真正的流式聊天
   * 使用Server-Sent Events (SSE) 实现流式响应
   * @param params 聊天请求参数
   * @param onMessage 收到消息片段时的回调函数
   * @param onComplete 完成时的回调函数
   * @param onError 错误时的回调函数
   */
  public static async streamChat(
    params: ChatRequest,
    onMessage: (text: string) => void,
    onComplete?: () => void,
    onError?: (error: Error) => void
  ): Promise<void> {
    return new Promise((resolve, reject) => {
      // 直接使用fetch API来实现SSE，因为axios默认不支持SSE
      fetch(`${import.meta.env.VITE_API_BASE || ''}/api/chat/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}` // 确保携带认证token
        },
        body: JSON.stringify(params)
      })
        .then(response => {
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
          }

          const reader = response.body?.getReader()
          const decoder = new TextDecoder()
          let buffer = ''

          const readStream = () => {
            if (!reader) {
              onComplete?.()
              resolve()
              return
            }

            reader.read()
              .then(({ done, value }) => {
                if (done) {
                  onComplete?.()
                  resolve()
                  return
                }

                // 解码数据并处理SSE格式
                const chunk = decoder.decode(value, { stream: true })
                buffer += chunk

                // 处理SSE消息格式: data: message\n\n
                const lines = buffer.split('\n')
                for (let i = 0; i < lines.length - 1; i++) {
                  const line = lines[i].trim()
                  if (line.startsWith('data:')) {
                    const message = line.substring(5).trim()
                    if (message) {
                      try {
                        // 尝试解析JSON，因为后端可能返回JSON格式的数据
                        const parsed = JSON.parse(message)
                        onMessage(parsed.text || message)
                      } catch {
                        // 如果不是JSON，直接使用文本
                        onMessage(message)
                      }
                    }
                  }
                }

                // 保留未处理的部分
                buffer = lines[lines.length - 1]

                // 继续读取
                readStream()
              })
              .catch(error => {
                console.error('Stream error:', error)
                onError?.(error)
                reject(error)
              })
          }

          readStream()
        })
        .catch(error => {
          console.error('Stream connection error:', error)
          onError?.(error)
          reject(error)
        })
    })
  }
}

export default ApiService
