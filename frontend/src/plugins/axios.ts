import type { App } from 'vue'
import request from '@/utils/request'
import ApiService from '@/services/api'

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $http: typeof request
    $api: typeof ApiService
  }
}

export default {
  install(app: App): void {
    // 挂载请求实例
    app.config.globalProperties.$http = request

    // 挂载 API 服务
    app.config.globalProperties.$api = ApiService

    // 提供注入方式
    app.provide('$http', request)
    app.provide('$api', ApiService)
  }
}
