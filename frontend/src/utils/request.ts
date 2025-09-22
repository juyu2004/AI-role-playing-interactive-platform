import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig
} from 'axios'
import type { BaseResponse, ErrorResponse, RequestConfig } from '@/types/api'

// 创建自定义 axios 实例
class Request {
  private instance: AxiosInstance

  constructor(config: AxiosRequestConfig) {
    this.instance = axios.create(config)
    this.setupInterceptors()
  }

  private setupInterceptors(): void {
    // 请求拦截器
    this.instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // 添加认证令牌
        const token = localStorage.getItem('token')
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error: Error) => {
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        const { data } = response
        // 如果后端返回的数据已经是标准格式，直接返回 data
        return data
      },
      (error: unknown) => {
        this.handleError(error)
        return Promise.reject(error)
      }
    )
  }

  private handleError(error: unknown): void {
    if (axios.isAxiosError(error)) {
      const response = error.response?.data as ErrorResponse | undefined
      const status = error.response?.status
      const message = response?.message || error.message

      console.error(`HTTP Error ${status}:`, message)

      switch (status) {
        case 401:
          console.error('未授权，请重新登录')
          // 可以跳转到登录页
          break
        case 403:
          console.error('拒绝访问')
          break
        case 404:
          console.error('请求资源不存在')
          break
        case 500:
          console.error('服务器内部错误')
          break
        default:
          console.error('网络错误', message)
      }
    } else if (error instanceof Error) {
      console.error('请求错误:', error.message)
    } else {
      console.error('未知错误:', error)
    }
  }

  // 泛型请求方法
  public request<T = unknown>(config: RequestConfig): Promise<BaseResponse<T>> {
    return this.instance.request(config)
  }

  public get<T = unknown>(url: string, config?: RequestConfig): Promise<BaseResponse<T>> {
    return this.instance.get(url, config)
  }

  public post<T = unknown>(url: string, data?: unknown, config?: RequestConfig): Promise<BaseResponse<T>> {
    return this.instance.post(url, data, config)
  }

  public put<T = unknown>(url: string, data?: unknown, config?: RequestConfig): Promise<BaseResponse<T>> {
    return this.instance.put(url, data, config)
  }

  public delete<T = unknown>(url: string, config?: RequestConfig): Promise<BaseResponse<T>> {
    return this.instance.delete(url, config)
  }

  // 获取原始的 axios 实例（如果需要）
  public getInstance(): AxiosInstance {
    return this.instance
  }
}

// 创建请求实例
const request = new Request({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

export default request
