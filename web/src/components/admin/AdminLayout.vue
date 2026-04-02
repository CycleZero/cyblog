<template>
  <el-container class="admin-layout">
    <el-aside :width="isCollapsed ? '64px' : '220px'" class="aside">
      <AdminSidebar />
    </el-aside>
    <el-container>
      <el-header height="60px" class="header">
        <AdminHeader @toggle-sidebar="handleToggleSidebar" />
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AdminSidebar from './AdminSidebar.vue'
import AdminHeader from './AdminHeader.vue'

const isCollapsed = ref(false)

function handleToggleSidebar(): void {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.aside {
  background-color: var(--el-bg-color-page);
  transition: width 0.3s;
  overflow: hidden;
}

.header {
  background-color: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.main {
  background-color: var(--el-bg-color-page);
  padding: 20px;
  overflow-y: auto;
}
</style>
