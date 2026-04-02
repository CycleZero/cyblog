<template>
  <header class="bg-white/80 backdrop-blur-xl shadow-lg sticky top-0 z-50 border-b border-sky-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center group cursor-pointer" @click="$router.push('/')">
          <div class="relative">
            <span class="text-3xl font-bold bg-gradient-to-r from-sky-600 via-cyan-600 to-teal-600 bg-clip-text text-transparent">
              🌊 Cyblog
            </span>
            <div class="absolute -bottom-1 left-0 w-0 h-0.5 bg-gradient-to-r from-sky-500 to-cyan-500 group-hover:w-full transition-all duration-500"></div>
          </div>
        </div>
        <nav class="hidden md:flex items-center space-x-1">
          <router-link
            to="/"
            :class="[
              'px-4 py-2 rounded-full font-medium transition-all duration-300',
              currentPath === '/'
                ? 'text-white bg-gradient-to-r from-sky-500 to-cyan-500'
                : 'text-gray-700 hover:text-white hover:bg-gradient-to-r hover:from-sky-500 hover:to-cyan-500'
            ]"
          >
            ✨ 首页
          </router-link>
          <router-link
            to="/article"
            :class="[
              'px-4 py-2 rounded-full font-medium transition-all duration-300',
              currentPath === '/article' || currentPath.startsWith('/article/')
                ? 'text-white bg-gradient-to-r from-sky-500 to-cyan-500'
                : 'text-gray-700 hover:text-white hover:bg-gradient-to-r hover:from-sky-500 hover:to-cyan-500'
            ]"
          >
            📝 文章
          </router-link>
          <router-link
            to="/category"
            :class="[
              'px-4 py-2 rounded-full font-medium transition-all duration-300',
              currentPath === '/category'
                ? 'text-white bg-gradient-to-r from-cyan-500 to-teal-500'
                : 'text-gray-700 hover:text-white hover:bg-gradient-to-r hover:from-cyan-500 hover:to-teal-500'
            ]"
          >
            🏷️ 分类
          </router-link>
        </nav>
        <div class="flex items-center space-x-3">
          <template v-if="!authStore.isLoggedIn">
            <router-link to="/login" class="px-5 py-2.5 text-gray-700 font-medium rounded-full border-2 border-gray-200 hover:border-sky-400 hover:text-sky-600 transition-all duration-300">
              登录
            </router-link>
            <router-link to="/register" class="px-5 py-2.5 bg-gradient-to-r from-sky-500 to-cyan-500 text-white rounded-full hover:shadow-lg hover:shadow-cyan-500/30 font-medium transition-all duration-300 transform hover:scale-105">
              注册
            </router-link>
          </template>
          <template v-else>
            <div class="relative" ref="userMenuRef">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center space-x-3 bg-gradient-to-r from-sky-50 to-cyan-50 px-4 py-2 rounded-full border border-sky-100 hover:border-sky-300 transition-all duration-300"
              >
                <div class="w-8 h-8 rounded-full bg-gradient-to-br from-sky-400 to-cyan-400 flex items-center justify-center text-white font-bold">
                  {{ authStore.userInfo?.name?.charAt(0) || 'U' }}
                </div>
                <span class="text-gray-700 font-medium hidden sm:block">{{ authStore.userInfo?.name }}</span>
                <svg class="w-4 h-4 text-gray-500 transition-transform duration-300" :class="{ 'rotate-180': showUserMenu }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </button>

              <!-- Dropdown Menu -->
              <Transition
                enter-active-class="transition duration-200 ease-out"
                enter-from-class="opacity-0 scale-95"
                enter-to-class="opacity-100 scale-100"
                leave-active-class="transition duration-150 ease-in"
                leave-from-class="opacity-100 scale-100"
                leave-to-class="opacity-0 scale-95"
              >
                <div
                  v-if="showUserMenu"
                  class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-2 overflow-hidden"
                >
                  <div class="px-4 py-3 border-b border-gray-100">
                    <p class="text-sm font-medium text-gray-900">{{ authStore.userInfo?.name }}</p>
                    <p class="text-xs text-gray-500 truncate">{{ authStore.userInfo?.email }}</p>
                  </div>
                  <!-- 管理界面入口（仅管理员可见） -->
                  <button
                    v-if="authStore.userInfo?.role === 'admin'"
                    @click="goToAdmin"
                    class="w-full text-left px-4 py-2.5 text-sm text-sky-600 hover:bg-sky-50 transition-colors flex items-center gap-2"
                  >
                    <span>⚙️</span>
                    管理界面
                  </button>
                  <button
                    @click="handleLogout"
                    class="w-full text-left px-4 py-2.5 text-sm text-red-600 hover:bg-red-50 transition-colors flex items-center gap-2"
                  >
                    <span>🚪</span>
                    退出登录
                  </button>
                </div>
              </Transition>
            </div>
          </template>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const showUserMenu = ref(false)
const userMenuRef = ref<HTMLElement | null>(null)

const currentPath = computed(() => route.path)

function handleLogout() {
  authStore.logout()
  showUserMenu.value = false
  router.push('/')
}

function goToAdmin(): void {
  showUserMenu.value = false
  router.push('/admin')
}

// 点击外部关闭菜单
function handleClickOutside(event: MouseEvent) {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target as Node)) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* 响应式导航菜单 */
@media (max-width: 768px) {
  nav {
    display: none;
  }
}
</style>
