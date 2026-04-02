import request from './request'
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from './types'

// 登录
export function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post('/auth/login', data)
}

// 注册
export function register(data: RegisterRequest): Promise<RegisterResponse> {
  return request.post('/auth/register', data)
}
