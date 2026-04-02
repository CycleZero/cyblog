<template>
  <el-container class="main-layout">
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="$router.push('/')">
          <h1>Cyblog</h1>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          :ellipsis="false"
          class="nav-menu"
          router
        >
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/article">文章</el-menu-item>
          <el-menu-item index="/category">分类</el-menu-item>
        </el-menu>
        <div class="user-actions">
          <template v-if="!isLoggedIn">
            <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </template>
          <template v-else>
            <el-dropdown @command="handleCommand">
              <span class="user-name">
                <el-avatar :size="32" :src="userInfo?.avatar" />
                <span class="name">{{ userInfo?.name }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </div>
      </div>
    </el-header>
    <el-main class="main">
      <slot />
    </el-main>
    <el-footer class="footer">
      <p>© 2026 Cyblog. All rights reserved.</p>
    </el-footer>
  </el-container>
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
.main-layout {
  min-height: 100vh;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
  height: 64px;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.logo {
  cursor: pointer;
  margin-right: 40px;
}

.logo h1 {
  margin: 0;
  font-size: 24px;
  color: #409eff;
}

.nav-menu {
  flex: 1;
  border: none;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}

.name {
  font-size: 14px;
}

.main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 20px;
  width: 100%;
}

.footer {
  text-align: center;
  background-color: #f5f7fa;
  color: #909399;
  padding: 20px;
}

.footer p {
  margin: 0;
}
</style>
