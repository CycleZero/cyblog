<template>
  <div class="sidebar-container">
    <div class="logo-container">
      <router-link to="/" class="logo-link">
        <span class="logo-icon">🌊</span>
        <span v-show="!isCollapsed" class="logo-text">Cyblog</span>
      </router-link>
    </div>
    <el-menu
      :default-active="activeMenu"
      :collapse="isCollapsed"
      :collapse-transition="false"
      class="sidebar-menu"
      router
    >
      <el-menu-item index="/admin">
        <el-icon><House /></el-icon>
        <template #title>仪表盘</template>
      </el-menu-item>
      <el-menu-item index="/admin/article">
        <el-icon><Document /></el-icon>
        <template #title>文章管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/category">
        <el-icon><FolderOpened /></el-icon>
        <template #title>分类管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/tag">
        <el-icon><PriceTag /></el-icon>
        <template #title>标签管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/comment">
        <el-icon><ChatDotRound /></el-icon>
        <template #title>评论管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/user">
        <el-icon><User /></el-icon>
        <template #title>用户管理</template>
      </el-menu-item>
      <el-menu-item index="/admin/setting">
        <el-icon><Setting /></el-icon>
        <template #title>系统设置</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  House,
  Document,
  FolderOpened,
  PriceTag,
  ChatDotRound,
  User,
  Setting,
} from '@element-plus/icons-vue'

const route = useRoute()

const isCollapsed = computed(() => {
  return false
})

const activeMenu = computed(() => {
  const path = route.path
  if (path.startsWith('/admin/article')) {
    return '/admin/article'
  }
  return path
})
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo-container {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--el-border-color-light);
}

.logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
}

.logo-icon {
  font-size: 24px;
  margin-right: 8px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}
</style>
