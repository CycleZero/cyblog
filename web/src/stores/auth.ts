import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { getToken, setToken as saveToken, removeToken } from '@/utils/auth'
import { login as apiLogin, register as apiRegister } from '@/api/auth'
import type { LoginRequest, RegisterRequest, LoginResponse, UserInfo } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  function setToken(newToken: string) {
    token.value = newToken
    saveToken(newToken)
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = info
    // 保存用户信息到 localStorage
    localStorage.setItem('cyblog_user', JSON.stringify(info))
  }

  function loadUserInfo() {
    const saved = localStorage.getItem('cyblog_user')
    if (saved) {
      try {
        userInfo.value = JSON.parse(saved)
      } catch {
        userInfo.value = null
      }
    }
  }

  function logout() {
    token.value = null
    userInfo.value = null
    removeToken()
    localStorage.removeItem('cyblog_user')
  }

  async function login(data: LoginRequest) {
    const res = await apiLogin(data)
    setToken(res.token)
    setUserInfo({
      id: res.id,
      name: res.name,
      email: res.email,
      avatar: res.avatar,
      role: res.role,
      status: 1,
      createdAt: new Date().toISOString(),
    })
    return res
  }

  async function register(data: RegisterRequest) {
    const res = await apiRegister(data)
    setToken(res.token)
    setUserInfo({
      id: res.id,
      name: res.name,
      email: res.email,
      avatar: res.avatar,
      role: 'user',
      status: 1,
      createdAt: new Date().toISOString(),
    })
    return res
  }

  // 初始化时加载用户信息
  if (token.value) {
    loadUserInfo()
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    setToken,
    setUserInfo,
    loadUserInfo,
    logout,
    login,
    register,
  }
})
