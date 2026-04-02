import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'
import ContentLayout from '@/components/layout/ContentLayout.vue'
import AdminLayout from '@/components/admin/AdminLayout.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // 前台页面 - 使用 ContentLayout
    {
      path: '/',
      component: ContentLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../views/HomeView.vue'),
        },
        {
          path: 'article',
          name: 'article-list',
          component: () => import('../views/ArticleListView.vue'),
        },
        {
          path: 'article/:id',
          name: 'article-detail',
          component: () => import('../views/ArticleDetailView.vue'),
        },
        {
          path: 'about',
          name: 'about',
          component: () => import('../views/AboutView.vue'),
        },
      ],
    },
    // 认证页面 - 不使用任何布局（独立页面）
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { guestOnly: true },
    },
    // 后台管理路由 - 使用 AdminLayout（独立布局，无博客页头页脚）
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/DashboardView.vue'),
        },
        {
          path: 'article',
          name: 'admin-article',
          component: () => import('../views/admin/ArticleListView.vue'),
        },
        {
          path: 'article/create',
          name: 'admin-article-create',
          component: () => import('../views/admin/ArticleFormView.vue'),
        },
        {
          path: 'article/:id/edit',
          name: 'admin-article-edit',
          component: () => import('../views/admin/ArticleFormView.vue'),
        },
        {
          path: 'category',
          name: 'admin-category',
          component: () => import('../views/admin/CategoryView.vue'),
        },
        {
          path: 'tag',
          name: 'admin-tag',
          component: () => import('../views/admin/TagView.vue'),
        },
        {
          path: 'comment',
          name: 'admin-comment',
          component: () => import('../views/admin/CommentView.vue'),
        },
        {
          path: 'user',
          name: 'admin-user',
          component: () => import('../views/admin/UserView.vue'),
        },
        {
          path: 'setting',
          name: 'admin-setting',
          component: () => import('../views/admin/SettingView.vue'),
        },
      ],
    },
  ],
})

// 路由守卫 - 使用新的返回方式
router.beforeEach((to) => {
  const authStore = useAuthStore()

  // 如果目标路由只允许未登录用户访问，且用户已登录，则跳转到首页
  if (to.meta.guestOnly && authStore.isLoggedIn) {
    return '/'
  }

  // 如果目标路由需要登录
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return '/login'
  }

  // 检查是否是管理员
  if (to.path.startsWith('/admin') && authStore.userInfo?.role !== 'admin') {
    ElMessage.error('没有权限访问后台管理系统')
    return '/'
  }

  // 继续导航
  return true
})

export default router
