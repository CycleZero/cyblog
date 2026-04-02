<template>
  <div class="min-h-screen bg-gradient-to-br from-sky-50 via-white to-cyan-50">
    <!-- Header -->
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="$router.push('/')">
          <span class="logo-text">🌊 Cyblog</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          :ellipsis="false"
          class="nav-menu"
          router
        >
          <el-menu-item index="/">✨ 首页</el-menu-item>
          <el-menu-item index="/article">📝 文章</el-menu-item>
          <el-menu-item index="/category">🏷️ 分类</el-menu-item>
        </el-menu>
        <div class="user-actions">
          <template v-if="!isLoggedIn">
            <el-button class="login-btn" @click="$router.push('/login')">登录</el-button>
            <el-button type="primary" class="register-btn" @click="$router.push('/register')">注册</el-button>
          </template>
          <template v-else>
            <el-dropdown @command="handleCommand">
              <span class="user-name">
                <el-avatar :size="36" :src="userInfo?.avatar" class="user-avatar" />
                <span class="name">{{ userInfo?.name }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">👤 个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout">🚪 退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </div>
      </div>
    </el-header>

    <!-- Main Content -->
    <el-main class="main">
      <slot />
    </el-main>

    <!-- Footer -->
    <el-footer class="footer">
      <div class="footer-content">
        <div class="footer-logo">🌊 Cyblog</div>
        <p class="footer-text">用心分享，每一篇文章都是一份礼物 🎁</p>
        <div class="footer-links">
          <a href="/" class="footer-link">🏠 首页</a>
          <a href="/article" class="footer-link">📝 文章</a>
          <a href="/category" class="footer-link">🏷️ 分类</a>
        </div>
        <div class="footer-copyright">
          © 2026 Cyblog. Made with 💖 and ☕
        </div>
      </div>
    </el-footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { UserInfo } from '@/api/types'
// import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
// const authStore = useAuthStore()

const activeMenu = computed(() => route.path)
// const isLoggedIn = computed(() => authStore.isLoggedIn)
// const userInfo = computed(() => authStore.userInfo)
const isLoggedIn = computed(() => false)
const userInfo = computed<UserInfo | null>(() => null)

function handleCommand(command: string) {
  if (command === 'profile') {
    router.push('/user')
  } else if (command === 'logout') {
    // authStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
/* Header Styles */
.header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.1);
  border-bottom: 1px solid rgba(14, 165, 233, 0.1);
  padding: 0;
  height: 72px;
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 24px;
  gap: 32px;
}

.logo {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.logo:hover {
  transform: scale(1.05);
}

.logo-text {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 50%, #14b8a6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-menu {
  flex: 1;
  border: none !important;
  background: transparent;
}

.nav-menu :deep(.el-menu-item) {
  font-weight: 500;
  color: #4b5563;
  border-radius: 12px;
  margin: 0 4px;
  transition: all 0.3s ease;
}

.nav-menu :deep(.el-menu-item:hover) {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
  color: #0ea5e9;
}

.nav-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
  color: white;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.login-btn {
  border-radius: 20px;
  border: 2px solid #e5e7eb;
  font-weight: 500;
  transition: all 0.3s ease;
}

.login-btn:hover {
  border-color: #0ea5e9;
  color: #0ea5e9;
}

.register-btn {
  border-radius: 20px;
  background: linear-gradient(135deg, #0ea5e9 0%, #06b6d4 100%);
  border: none;
  font-weight: 500;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(14, 165, 233, 0.3);
}

.user-name {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 10px;
  padding: 6px 16px;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
  border-radius: 24px;
  transition: all 0.3s ease;
}

.user-name:hover {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.2) 0%, rgba(6, 182, 212, 0.2) 100%);
}

.user-avatar {
  border: 2px solid rgba(14, 165, 233, 0.3);
}

.name {
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
}

/* Main Styles */
.main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
  width: 100%;
  min-height: calc(100vh - 200px);
}

/* Footer Styles */
.footer {
  background: linear-gradient(135deg, #1f2937 0%, #0c4a6e 50%, #1f2937 100%);
  padding: 48px 24px;
  position: relative;
  overflow: hidden;
}

.footer::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Ccircle cx='30' cy='30' r='2'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  pointer-events: none;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
  position: relative;
  z-index: 1;
}

.footer-logo {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #7dd3fc 0%, #22d3ee 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 16px;
}

.footer-text {
  color: #9ca3af;
  margin-bottom: 24px;
  font-size: 15px;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 32px;
}

.footer-link {
  color: #d1d5db;
  text-decoration: none;
  transition: color 0.3s ease;
  font-size: 14px;
}

.footer-link:hover {
  color: #7dd3fc;
}

.footer-copyright {
  color: #6b7280;
  font-size: 13px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* Dropdown Menu Styles */
:deep(.el-dropdown-menu) {
  border-radius: 16px;
  padding: 8px;
  box-shadow: 0 10px 40px rgba(14, 165, 233, 0.2);
  border: 1px solid rgba(14, 165, 233, 0.1);
}

:deep(.el-dropdown-menu__item) {
  border-radius: 10px;
  padding: 12px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-dropdown-menu__item:hover) {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
  color: #0ea5e9;
}

/* Responsive Design */
@media (max-width: 768px) {
  .header-content {
    padding: 0 16px;
    gap: 16px;
  }

  .logo-text {
    font-size: 24px;
  }

  .nav-menu {
    display: none;
  }

  .main {
    padding: 24px 16px;
  }

  .footer {
    padding: 32px 16px;
  }
}
</style>
