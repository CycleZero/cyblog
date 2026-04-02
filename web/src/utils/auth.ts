const TOKEN_KEY = 'cyblog_token'
const USER_KEY = 'cyblog_user'

// 获取token
export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

// 设置token
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

// 移除token
export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

// 获取用户信息
export function getUserInfo(): any {
  const saved = localStorage.getItem(USER_KEY)
  if (saved) {
    try {
      return JSON.parse(saved)
    } catch {
      return null
    }
  }
  return null
}

// 设置用户信息
export function setUserInfo(info: any): void {
  localStorage.setItem(USER_KEY, JSON.stringify(info))
}

// 移除用户信息
export function removeUserInfo(): void {
  localStorage.removeItem(USER_KEY)
}

// 清除所有认证信息
export function clearAuth(): void {
  removeToken()
  removeUserInfo()
}
