# Cyblog 前端技术选型与设计方案

## 1. 项目概述

Cyblog 是一个基于 Go + Gin 的博客系统后端，本方案为其配套的前端设计。

### 1.1 功能需求
| 模块 | 功能 |
|------|------|
| **用户模块** | 用户注册、登录、个人信息管理 |
| **文章模块** | 文章列表、文章详情、文章搜索、文章创建/编辑/删除、文章点赞 |
| **分类模块** | 分类列表、分类管理（管理员） |
| **标签模块** | 标签列表、标签管理（管理员） |
| **评论模块** | 评论列表、评论发布、评论回复、评论编辑/删除、评论点赞 |
| **管理后台** | 文章管理、分类管理、标签管理、评论管理 |

### 1.2 技术约束
- 后端已提供 RESTful API（以 `/api` 前缀）
- API 文档：Swagger JSON/YAML 格式
- 认证方式：JWT Token

## 2. 技术选型

### 2.1 核心框架
| 技术 | 选型 | 说明 |
|------|------|------|
| **框架** | Vue 3 + TypeScript | 现代化、类型安全 |
| **构建工具** | Vite | 快速开发体验 |
| **路由** | Vue Router 4 | 单页应用路由管理 |
| **状态管理** | Pinia | Vue 3 官方推荐 |
| **HTTP 客户端** | Axios | 请求拦截、响应拦截 |
| **UI 组件库** | Element Plus | 国内流行，生态完善 |
| **Markdown 编辑器** | @bytemd/vue-next | 字节跳动 Markdown 编辑器 |
| **代码高亮** | Prism.js | 代码语法高亮 |
| **测试框架** | Vitest + Playwright | 单元测试和 E2E 测试 |
| **代码规范** | ESLint + Prettier + Oxlint | 现代化代码检查工具 |

### 2.2 开发工具
| 工具 | 说明 |
|------|------|
| **ESLint** | 代码规范检查 |
| **Prettier** | 代码格式化 |
| **Oxlint** | 快速代码检查 |
| **Husky + lint-staged** | Git 提交前检查 |
| **Vue DevTools** | Vue 浏览器调试工具 |

## 3. 项目结构设计

```
web/
├── public/                 # 静态资源
│   └── favicon.ico
├── src/
│   ├── api/              # API 请求层
│   │   ├── request.ts   # Axios 实例配置
│   │   ├── types.ts     # API 类型定义
│   │   ├── auth.ts     # 认证相关 API
│   │   ├── article.ts  # 文章相关 API
│   │   ├── category.ts # 分类相关 API
│   │   ├── tag.ts      # 标签相关 API
│   │   ├── comment.ts  # 评论相关 API
│   │   └── user.ts     # 用户相关 API
│   ├── assets/           # 资源文件
│   │   ├── base.css
│   │   ├── logo.svg
│   │   └── main.css
│   ├── components/       # 公共组件
│   │   ├── __tests__/   # 组件测试
│   │   ├── icons/       # 图标组件
│   │   ├── layout/      # 布局组件
│   │   ├── editor/      # Markdown 编辑器
│   │   ├── comment/     # 评论组件
│   │   └── common/      # 通用组件
│   ├── composables/      # 组合式函数
│   │   ├── useAuth.ts
│   │   ├── useArticle.ts
│   │   └── ...
│   ├── router/           # 路由配置
│   │   └── index.ts
│   ├── stores/           # Pinia 状态管理
│   │   ├── counter.ts   # 示例 store
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   └── ...
│   ├── utils/           # 工具函数
│   │   ├── auth.ts     # 认证工具
│   │   ├── storage.ts  # 本地存储
│   │   └── format.ts   # 格式化工具
│   ├── views/            # 页面视图
│   │   ├── HomeView.vue # 首页视图（已存在）
│   │   ├── AboutView.vue # 关于视图（已存在）
│   │   ├── home/        # 首页
│   │   ├── article/     # 文章相关页
│   │   ├── category/    # 分类相关页
│   │   ├── user/        # 用户相关页
│   │   ├── admin/       # 管理后台
│   │   └── ...
│   ├── App.vue
│   └── main.ts
├── e2e/                  # E2E 测试
│   ├── tsconfig.json
│   └── vue.spec.ts
├── .env.development
├── .env.production
├── .eslintrc.cjs
├── .prettierrc.json
├── .oxlintrc.json
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.app.json
├── tsconfig.node.json
├── tsconfig.vitest.json
├── vite.config.ts
├── vitest.config.ts
├── playwright.config.ts
└── README.md
```

## 4. 核心模块设计

### 4.1 API 请求层设计

#### src/api/request.ts
```typescript
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from '@/utils/auth'
import router from '@/router'

// 响应数据结构
interface ResponseData<T = any> {
  code: number
  msg: string
  data: T
}

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ResponseData>) => {
    const { code, msg, data } = response.data
    if (code === 200) {
      return data
    } else {
      ElMessage.error(msg || '请求失败')
      return Promise.reject(new Error(msg))
    }
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        ElMessage.error('登录已过期，请重新登录')
        removeToken()
        router.push('/login')
      } else if (status === 403) {
        ElMessage.error('无权限操作')
      } else {
        ElMessage.error(error.response.data?.msg || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请稍后重试')
    }
    return Promise.reject(error)
  }
)

export default request
```

### 4.2 状态管理设计

#### src/stores/auth.ts
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login as loginApi, register as registerApi, getUserInfo } from '@/api/auth'
import type { LoginRequest, RegisterRequest, UserInfo } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(getToken())
  const userInfo = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  // 登录
  async function login(data: LoginRequest) {
    const res = await loginApi(data)
    token.value = res.token
    setToken(res.token)
    userInfo.value = res
  }

  // 注册
  async function register(data: RegisterRequest) {
    const res = await registerApi(data)
    token.value = res.token
    setToken(res.token)
    userInfo.value = res
  }

  // 获取用户信息
  async function fetchUserInfo() {
    const res = await getUserInfo()
    userInfo.value = res
  }

  // 登出
  function logout() {
    token.value = null
    userInfo.value = null
    removeToken()
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    login,
    register,
    fetchUserInfo,
    logout,
  }
})
```

### 4.3 路由设计

#### src/router/index.ts
```typescript
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/article',
      name: 'ArticleList',
      component: () => import('@/views/article/list.vue'),
    },
    {
      path: '/article/:id',
      name: 'ArticleDetail',
      component: () => import('@/views/article/detail.vue'),
    },
    {
      path: '/article/create',
      name: 'ArticleCreate',
      component: () => import('@/views/article/edit.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/article/:id/edit',
      name: 'ArticleEdit',
      component: () => import('@/views/article/edit.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/category',
      name: 'CategoryList',
      component: () => import('@/views/category/list.vue'),
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/auth/login.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/auth/register.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/user',
      name: 'UserCenter',
      component: () => import('@/views/user/index.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('@/views/admin/index.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: 'article',
          name: 'AdminArticle',
          component: () => import('@/views/admin/article.vue'),
        },
        {
          path: 'category',
          name: 'AdminCategory',
          component: () => import('@/views/admin/category.vue'),
        },
        {
          path: 'tag',
          name: 'AdminTag',
          component: () => import('@/views/admin/tag.vue'),
        },
        {
          path: 'comment',
          name: 'AdminComment',
          component: () => import('@/views/admin/comment.vue'),
        },
      ],
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.guestOnly && authStore.isLoggedIn) {
    next({ name: 'home' })
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
```

### 4.4 评论组件设计

#### src/components/comment/CommentList.vue
```vue
<template>
  <div class="comment-list">
    <el-divider content-position="left">评论 ({{ total }})</el-divider>
    
    <!-- 发表评论 -->
    <CommentEditor
      v-if="isLoggedIn"
      :article-id="articleId"
      @created="handleCommentCreated"
    />
    <el-empty v-else description="登录后发表评论" />

    <!-- 评论列表 -->
    <div class="comments">
      <CommentItem
        v-for="comment in list"
        :key="comment.id"
        :comment="comment"
        @reply="handleReply"
      />
    </div>

    <!-- 加载更多 -->
    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="currentPage"
      :total="total"
      :page-size="pageSize"
      @current-change="fetchComments"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getComments } from '@/api/comment'
import type { CommentResponse } from '@/api/types'
import CommentItem from './CommentItem.vue'
import CommentEditor from './CommentEditor.vue'

const props = defineProps<{
  articleId: number
}>()

const authStore = useAuthStore()
const list = ref<CommentResponse[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const isLoggedIn = computed(() => authStore.isLoggedIn)

async function fetchComments() {
  const res = await getComments({
    article_id: props.articleId,
    page: currentPage.value,
    page_size: pageSize.value,
  })
  list.value = res.list
  total.value = res.total
}

function handleCommentCreated() {
  fetchComments()
}

function handleReply(parentId: number) {
  // 处理回复
}

watch(() => props.articleId, () => {
  currentPage.value = 1
  fetchComments()
})

onMounted(() => {
  fetchComments()
})
</script>
```

## 5. UI 设计要点

### 5.1 响应式布局
- PC端：侧边栏 + 主内容区
- 移动端：顶部导航 + 抽屉菜单

### 5.2 主题配置
```typescript
// src/utils/theme.ts
import type { InstallOptions } from 'element-plus'

export const elementPlusConfig: InstallOptions = {
  size: 'default',
  zIndex: 3000,
}
```

### 5.3 代码高亮
使用 Prism.js 进行代码语法高亮：
```typescript
import Prism from 'prismjs'
import 'prismjs/themes/prism-tomorrow.css'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-css'
import 'prismjs/components/prism-markup'
```

## 6. 开发与构建

### 6.1 环境变量
```env
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api

# .env.production
VITE_API_BASE_URL=https://api.example.com
```

### 6.2 构建命令
```json
{
  "scripts": {
    "dev": "vite",
    "build": "run-p type-check \"build-only {@}\" --",
    "preview": "vite preview",
    "test:unit": "vitest",
    "test:e2e": "playwright test",
    "build-only": "vite build",
    "type-check": "vue-tsc --build",
    "lint": "run-s lint:*",
    "lint:oxlint": "oxlint . --fix",
    "lint:eslint": "eslint . --fix --cache",
    "format": "prettier --write --experimental-cli src/"
  }
}
```

## 7. 部署方案

### 7.1 Nginx 配置示例
```nginx
server {
    listen 80;
    server_name yourdomain.com;

    root /var/www/cyblog/web/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

## 8. 后续优化方向

1. **性能优化**
   - 路由懒加载
   - 图片懒加载
   - 使用 CDN 加速静态资源
   - 评论列表虚拟滚动

2. **用户体验**
   - 评论实时通知
   - 离线缓存（Service Worker）
   - PWA 支持

3. **SEO 优化**
   - 服务端渲染 (Nuxt.js) 或预渲染
   - Meta 标签动态设置
   - Sitemap 生成

4. **安全加固**
   - XSS 防护
   - CSRF 防护
