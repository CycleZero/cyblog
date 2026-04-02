<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm sticky top-0 z-50">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <a href="/" class="text-2xl font-bold text-blue-600">Cyblog</a>
          </div>
          <nav class="hidden md:flex space-x-8">
            <a href="/" class="text-gray-900 font-medium hover:text-blue-600 transition-colors"
              >首页</a
            >
            <a href="/article" class="text-gray-600 hover:text-blue-600 transition-colors">文章</a>
            <a href="/category" class="text-gray-600 hover:text-blue-600 transition-colors">分类</a>
          </nav>
          <div class="flex items-center space-x-4">
            <template v-if="!isLoggedIn">
              <a href="/login" class="px-4 py-2 text-gray-700 hover:text-gray-900 font-medium"
                >登录</a
              >
              <a
                href="/register"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors"
                >注册</a
              >
            </template>
            <template v-else>
              <div class="flex items-center space-x-3">
                <img :src="userInfo?.avatar" class="w-8 h-8 rounded-full" />
                <span class="text-gray-700">{{ userInfo?.name }}</span>
              </div>
            </template>
          </div>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="bg-gradient-to-br from-blue-600 via-purple-600 to-indigo-700 text-white py-20">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h1 class="text-5xl md:text-6xl font-bold mb-6">欢迎来到 Cyblog</h1>
        <p class="text-xl md:text-2xl opacity-90">分享技术，记录生活</p>
      </div>
    </section>

    <!-- Main Content -->
    <main class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="flex gap-8">
        <!-- Article List -->
        <div class="flex-1">
          <h2 class="text-2xl font-bold text-gray-900 mb-8">最新文章</h2>

          <div v-if="loading" class="space-y-6">
            <div v-for="i in 5" :key="i" class="bg-white rounded-xl shadow-sm p-6">
              <div class="h-7 bg-gray-200 rounded w-3/4 mb-4 animate-pulse"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2 mb-3 animate-pulse"></div>
              <div class="h-4 bg-gray-200 rounded w-full mb-2 animate-pulse"></div>
              <div class="h-4 bg-gray-200 rounded w-2/3 animate-pulse"></div>
            </div>
          </div>

          <template v-else>
            <div v-if="articles?.length > 0" class="space-y-6">
              <article
                v-for="article in articles"
                :key="article.id"
                class="bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow p-6 cursor-pointer"
                @click="$router.push(`/article/${article.id}`)"
              >
                <h3
                  class="text-xl font-bold text-gray-900 mb-3 hover:text-blue-600 transition-colors"
                >
                  {{ article.title }}
                </h3>
                <div class="flex items-center gap-4 text-sm text-gray-500 mb-4">
                  <span
                    v-if="article.category"
                    class="bg-blue-100 text-blue-700 px-3 py-1 rounded-full font-medium"
                  >
                    {{ article.category.name }}
                  </span>
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                      ></path>
                    </svg>
                    {{ article.author.name }}
                  </span>
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                      ></path>
                    </svg>
                    {{ formatDate(article.created_at) }}
                  </span>
                </div>
                <p class="text-gray-600 line-clamp-3 mb-4">{{ article.summary }}</p>
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-4 text-sm text-gray-500">
                    <span class="flex items-center gap-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                        ></path>
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                        ></path>
                      </svg>
                      {{ article.views }}
                    </span>
                    <span class="flex items-center gap-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                        ></path>
                      </svg>
                      {{ article.likes }}
                    </span>
                  </div>
                  <div class="flex gap-2">
                    <span
                      v-for="tag in article.tags || []"
                      :key="tag.id"
                      class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded"
                    >
                      {{ tag.name }}
                    </span>
                  </div>
                </div>
              </article>
            </div>

            <!-- Pagination -->
            <div v-if="total > 0" class="flex justify-center items-center gap-2 mt-12">
              <button
                :disabled="currentPage <= 1"
                @click="
                  currentPage-- ,
                  fetchArticles()
                "
                class="px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                上一页
              </button>
              <span class="text-gray-600">
                第 {{ currentPage }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页
              </span>
              <button
                :disabled="currentPage >= Math.ceil(total / pageSize)"
                @click="
                  currentPage++,
                  fetchArticles()
                "
                class="px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                下一页
              </button>
            </div>

            <div v-else class="text-center py-16">
              <svg
                class="w-16 h-16 mx-auto text-gray-300 mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                ></path>
              </svg>
              <p class="text-gray-500 text-lg">暂无文章</p>
            </div>
          </template>
        </div>

        <!-- Sidebar -->
        <aside class="w-80 hidden lg:block">
          <!-- Tags -->
          <div class="bg-white rounded-xl shadow-sm p-6 mb-6">
            <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
              <svg
                class="w-5 h-5 text-blue-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
                ></path>
              </svg>
              热门标签
            </h3>
            <div v-if="tagsLoading" class="space-y-2">
              <div v-for="i in 5" :key="i" class="h-6 bg-gray-200 rounded w-24 animate-pulse"></div>
            </div>
            <div v-else-if="tags?.length > 0" class="flex flex-wrap gap-2">
              <button
                v-for="tag in tags"
                :key="tag.id"
                @click="goToArticleByTag(tag.id)"
                class="px-3 py-1.5 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-blue-100 hover:text-blue-700 transition-colors"
              >
                {{ tag.name }} ({{ tag.count }})
              </button>
            </div>
            <div v-else class="text-center py-6 text-gray-500">暂无标签</div>
          </div>

          <!-- Categories -->
          <div class="bg-white rounded-xl shadow-sm p-6">
            <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
              <svg
                class="w-5 h-5 text-blue-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4 6h16M4 12h16M4 18h16"
                ></path>
              </svg>
              分类
            </h3>
            <div v-if="categoriesLoading" class="space-y-2">
              <div
                v-for="i in 5"
                :key="i"
                class="h-6 bg-gray-200 rounded w-full animate-pulse"
              ></div>
            </div>
            <div v-else-if="categories.length > 0" class="space-y-1">
              <button
                v-for="category in categories"
                :key="category.id"
                @click="goToArticleByCategory(category.id)"
                class="w-full text-left px-3 py-2 text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
              >
                {{ category.name }}
              </button>
            </div>
            <div v-else class="text-center py-6 text-gray-500">暂无分类</div>
          </div>
        </aside>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-gray-900 text-gray-400 py-12 mt-16">
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <p>© 2026 Cyblog. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles } from '@/api/article'
import { getTags } from '@/api/tag'
import { getCategories } from '@/api/category'
import type { Article, Tag, Category, User } from '@/api/types'

const router = useRouter()
const articles = ref<Article[]>([])
const tags = ref<Tag[]>([])
const categories = ref<Category[]>([])
const loading = ref(true)
const tagsLoading = ref(true)
const categoriesLoading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const isLoggedIn = ref(false)
const userInfo = ref<User|null>(null)

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

async function fetchArticles() {
  try {
    loading.value = true
    const res = await getArticles({
      page: currentPage.value,
      page_size: pageSize.value,
      status: 2,
      sort_by: 'created_at',
      sort_order: 'desc',
    })
    articles.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取文章失败', error)
  } finally {
    loading.value = false
  }
}

async function fetchTags() {
  try {
    tagsLoading.value = true
    const res = await getTags({ page: 1, page_size: 20 })
    tags.value = res.list
  } catch (error) {
    console.error('获取标签失败', error)
  } finally {
    tagsLoading.value = false
  }
}

async function fetchCategories() {
  try {
    categoriesLoading.value = true
    const res = await getCategories({ page: 1, page_size: 20 })
    categories.value = res.list
  } catch (error) {
    console.error('获取分类失败', error)
  } finally {
    categoriesLoading.value = false
  }
}

function goToArticleByTag(tagId: number) {
  router.push(`/article?tag_id=${tagId}`)
}

function goToArticleByCategory(categoryId: number) {
  router.push(`/article?category_id=${categoryId}`)
}

onMounted(() => {
  fetchArticles()
  fetchTags()
  fetchCategories()
})
</script>
