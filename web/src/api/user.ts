import request from './request'
import type {
  UserInfo,
  GetUserRequest,
  GetUserResponse,
  UpdateUserRequest,
} from './types'

// 获取当前登录用户信息
export function getCurrentUserInfo(): Promise<UserInfo> {
  return request.get('/user/info')
}

// 获取用户信息
export function getUserInfo(data: GetUserRequest): Promise<GetUserResponse> {
  return request.get(`/user/${data.id}`)
}

// 更新用户信息
export function updateUserInfo(data: UpdateUserRequest): Promise<UserInfo> {
  return request.put('/user/update', data)
}
