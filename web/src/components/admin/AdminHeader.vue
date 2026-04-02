<template>
  <div class="header-container">
    <div class="left">
      <el-button
        text
        @click="handleToggleSidebar"
      >
        <el-icon :size="20">
          <Fold v-if="!isSidebarCollapsed" />
          <Expand v-else />
        </el-icon>
      </el-button>
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/admin' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item v-if="currentRouteName">
          {{ currentRouteName }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="right">
      <el-dropdown @command="handleCommand">
        <span class="user-info">
          <el-avatar :size="32" :src="userAvatar">
            {{ userName?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="user-name">{{ userName }}</span>
          <el-icon><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人中心</el-dropdown-item>
            <el-dropdown-item command="home">返回前台</el-dropdown-item>
            <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessageBox } from 'element-plus'
import { Fold, Expand, ArrowDown } from '@element-plus/icons-vue'

defineEmits<{
  toggleSidebar: []
}>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isSidebarCollapsed = computed(() => false)

const currentRouteName = computed(() => {
  const routeNames: Record<string, string> = {
    'admin-dashboard': '仪表盘',
    'admin-article': '文章管理',
    'admin-article-create': '创建文章',
    'admin-article-edit': '编辑文章',
    'admin-category': '分类管理',
    'admin-tag': '标签管理',
    'admin-comment': '评论管理',
    'admin-user': '用户管理',
    'admin-setting': '系统设置',
  }
  return routeNames[route.name as string] || ''
})

const userName = computed(() => authStore.userInfo?.name || '用户')
const userAvatar = computed(() => authStore.userInfo?.avatar || '')

function handleToggleSidebar(): void {
  // emit('toggleSidebar')
}

async function handleCommand(command: string): Promise<void> {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'home':
      router.push('/')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          type: 'warning',
        })
        authStore.logout()
        router.push('/login')
      } catch {
        // 用户取消
      }
      break
  }
}
</script>

<style scoped>
.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: var(--el-fill-color-light);
}

.user-name {
  color: var(--el-text-color-primary);
  font-size: 14px;
}
</style>
