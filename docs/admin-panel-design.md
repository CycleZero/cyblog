# Cyblog 后台管理系统实现规范

## 一、代码规范要求

### 1.1 TypeScript 规范
- 严格模式：`"strict": true`
- 所有组件使用 `<script setup lang="ts">` 语法
- 接口和类型定义必须完整，无 `any` 类型
- 使用 `type` 而非 `interface` 定义简单类型
- 箭头函数优先，必要时使用普通函数
- 导出函数必须有 JSDoc 注释

### 1.2 ESLint 规范
- 遵循 `@typescript-eslint/recommended` 规则
- 遵循 `vue/vue3-recommended` 规则
- 组件名使用 PascalCase
- 事件处理函数使用 `handle` 前缀
- 布尔值变量使用 `is`、`has`、`can` 前缀

### 1.3 Element Plus 使用规范
- **优先使用 Element Plus 原生样式**，不过度自定义
- 布局使用 `el-container` / `el-aside` / `el-header` / `el-main`
- 表格使用 `el-table` / `el-table-column`
- 表单使用 `el-form` / `el-form-item`
- 对话框使用 `el-dialog`
- 分页使用 `el-pagination`
- 消息提示使用 `ElMessage`
- 下拉菜单使用 `el-dropdown`
- 标签使用 `el-tag`
- 卡片使用 `el-card`
- 按钮使用 `el-button`

### 1.4 CSS 规范
- 优先使用 Element Plus 组件的默认样式
- 如需自定义，使用 `<style scoped>` 局部样式
- 避免全局 CSS 污染
- 使用 CSS 变量统一管理颜色

---

## 二、文件结构

```昂
web/src/
├── api/
│   ├── admin.ts          # 后台管理 API（待实现）
│   └── ...（现有）
├── stores/
│   └── admin.ts          # 后台状态管理（待实现）
├── views/
│   └── admin/            # 后台页面（待实现）
│       ├── DashboardView.vue
│       ├── ArticleListView.vue
│       ├── ArticleFormView.vue
│       ├── CategoryView.vue
│       ├── TagView.vue
│       ├── CommentView.vue
│       ├── UserView.vue
│       └── SettingView.vue
├── components/
│   └── admin/            # 后台组件（待实现）
│       ├── AdminLayout.vue
│       ├── AdminSidebar.vue
│       └── AdminHeader.vue
└── router/
    └── index.ts          # 更新路由
```

---

## 三、API 类型定义

### 3.1 已存在 API（直接使用）

```typescript
// web/src/api/types.ts 已定义
import type {
  User,
  UserInfo,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  Article,
  ArticleListParams,
  ArticleListResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
  Category,
  CategoryListParams,
  CategoryListResponse,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  Tag,
  TagListParams,
  TagListResponse,
  CreateTagRequest,
  UpdateTagRequest,
  Comment,
  CommentListParams,
  CommentListResponse,
  CreateCommentRequest,
  UpdateCommentRequest,
} from './types'
```

### 3.2 后台 API 类型（待添加）

```typescript
// web/src/api/admin.ts

import request from './request'

// ============== 仪表盘 ==============

/**
 * 获取仪表盘统计数据
 */
export function getDashboard(): Promise<DashboardData> {
  return request.get('/admin/dashboard')
}

// ============== 文章管理 ==============

/**
 * 管理端获取文章列表（包含所有状态）
 */
export function getAdminArticles(params: AdminArticleListParams): Promise<ArticleListResponse> {
  return request.get('/admin/articles', { params })
}

/**
 * 置顶/取消置顶文章
 */
export function setArticleTop(id: number, isTop: boolean): Promise<void> {
  return request.post(`/admin/articles/${id}/top`, { is_top: isTop })
}

/**
 * 批量删除文章
 */
export function batchDeleteArticles(ids: number[]): Promise<void> {
  return request.post('/admin/articles/batch-delete', { ids })
}

/**
 * 批量更新文章状态
 */
export function batchUpdateArticleStatus(ids: number[], status: number): Promise<void> {
  return request.post('/admin/articles/batch-status', { ids, status })
}

// ============== 评论管理 ==============

/**
 * 管理端获取评论列表
 */
export function getAdminComments(params: AdminCommentListParams): Promise<CommentListResponse> {
  return request.get('/admin/comments', { params })
}

/**
 * 审核评论
 */
export function auditComment(id: number, status: number): Promise<void> {
  return request.put(`/admin/comments/${id}/audit`, { status })
}

// ============== 用户管理 ==============

/**
 * 获取用户列表
 */
export function getAdminUsers(params: AdminUserListParams): Promise<AdminUserListResponse> {
  return request.get('/admin/users', { params })
}

/**
 * 更新用户角色/状态
 */
export function updateAdminUser(id: number, data: UpdateAdminUserRequest): Promise<void> {
  return request.put(`/admin/users/${id}`, data)
}

/**
 * 删除用户
 */
export function deleteAdminUser(id: number): Promise<void> {
  return request.delete(`/admin/users/${id}`)
}

// ============== 系统设置 ==============

/**
 * 获取系统设置
 */
export function getSettings(): Promise<SystemSettings> {
  return request.get('/admin/settings')
}

/**
 * 更新系统设置
 */
export function updateSettings(data: SystemSettings): Promise<void> {
  return request.put('/admin/settings', data)
}
```

---

## 四、组件实现规范

### 4.1 布局组件 (AdminLayout)

```vue
<template>
  <el-container class="admin-layout">
    <el-aside width="220px">
      <AdminSidebar />
    </el-aside>
    <el-container>
      <el-header height="60px">
        <AdminHeader />
      </el-header>
      <el-main>
        <slot />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import AdminSidebar from './AdminSidebar.vue'
import AdminHeader from './AdminHeader.vue'
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}
</style>
```

### 4.2 表格页面规范

```vue
<template>
  <div class="article-list">
    <!-- 搜索筛选区域 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="关键词">
          <el-input v-model="queryParams.keyword" placeholder="请输入关键词" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="全部" :value="undefined" />
            <el-option label="草稿" :value="1" />
            <el-option label="已发布" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card class="table-card">
      <template #header>
        <div class="table-header">
          <span>文章列表</span>
          <el-button type="primary" @click="handleCreate">新增</el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { Article } from '@/api/types'
import { getAdminArticles, deleteArticle } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

// 类型定义
interface QueryParams {
  keyword?: string
  status?: number
  page: number
  pageSize: number
}

interface Pagination {
  page: number
  pageSize: number
  total: number
}

// 状态
const loading = ref(false)
const tableData = ref<Article[]>([])
const selectedRows = ref<Article[]>([])
const queryParams = reactive<QueryParams>({
  page: 1,
  pageSize: 10,
})
const pagination = reactive<Pagination>({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 方法
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const res = await getAdminArticles(queryParams)
    tableData.value = res.list
    pagination.total = res.total
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

function handleSearch(): void {
  pagination.page = 1
  fetchData()
}

function handleReset(): void {
  queryParams.keyword = undefined
  queryParams.status = undefined
  handleSearch()
}

function handleSelectionChange(rows: Article[]): void {
  selectedRows.value = rows
}

function handleCreate(): void {
  // 跳转创建页面
}

function handleEdit(row: Article): void {
  // 跳转编辑页面
}

async function handleDelete(row: Article): Promise<void> {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning',
    })
    await deleteArticle(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

function handleSizeChange(size: number): void {
  pagination.pageSize = size
  fetchData()
}

function handlePageChange(page: number): void {
  pagination.page = page
  fetchData()
}

function getStatusType(status: number): string {
  const types: Record<number, string> = {
    1: 'info',
    2: 'success',
    3: 'warning',
  }
  return types[status] || 'info'
}

function getStatusText(status: number): string {
  const texts: Record<number, string> = {
    1: '草稿',
    2: '已发布',
    3: '私密',
  }
  return texts[status] || '未知'
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.article-list {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
```

---

## 五、确认事项

请确认以下内容，我将严格遵循执行：

1. **TypeScript 严格模式**：所有类型定义必须完整，无 any
2. **Element Plus 原生样式**：不过度自定义，优先使用组件默认样式
3. **ESLint 规范**：遵循 @typescript-eslint 和 vue3-recommended 规则
4. **待添加 API**：是否需要我先实现前端页面骨架，使用模拟数据？

确认后我立即开始实现。
