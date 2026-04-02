import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { getToken, removeToken } from '@/utils/auth'
import router from '@/router'

// 后端统一响应接口（泛型默认值用 unknown，替代 any，符合ESLint）
interface ResponseData<T = unknown> {
  code: number
  msg: string
  data: T
}

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
})

// 请求拦截器（无修改，仅保留原有逻辑）
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// 响应拦截器（核心修复：泛型类型 + 严格匹配 Axios 类型）
request.interceptors.response.use(
  <T>(response: AxiosResponse<ResponseData<T>>) => {
    const { code, msg, data } = response.data
    // 业务状态码成功：直接返回 data（通过泛型保证类型安全）
    if (code === 200) {
      return data as unknown as AxiosResponse['data']
    }
    // 业务状态码失败：抛出错误
    console.error(msg || '请求失败')
    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error) => {
    // 原有错误处理逻辑（无修改）
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        console.error('登录已过期，请重新登录')
        removeToken()
        router.push('/login')
      } else if (status === 403) {
        console.error('无权限操作')
      } else {
        console.error(error.response.data?.msg || '请求失败')
      }
    } else {
      console.error('网络错误，请稍后重试')
    }
    return Promise.reject(error)
  },
)

export default request
