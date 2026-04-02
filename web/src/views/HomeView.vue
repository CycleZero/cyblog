<template>
  <!-- Hero Section - 清新天蓝色风格 -->
  <section class="relative py-24 px-4 overflow-hidden">
    <div class="max-w-7xl mx-auto text-center relative z-10">
      <div class="animate-bounce-in">
        <span class="inline-block px-6 py-2 bg-gradient-to-r from-sky-100 to-cyan-100 text-sky-700 rounded-full text-sm font-medium mb-6">
          ✨ 欢迎来到 Cyblog
        </span>
      </div>
      <h1 class="text-5xl md:text-7xl font-bold mb-6 animate-fade-in-up">
        <span class="bg-gradient-to-r from-sky-600 via-cyan-600 to-teal-600 bg-clip-text text-transparent bg-[length:200%_auto] animate-gradient">
          分享技术，记录生活
        </span>
      </h1>
      <p class="text-xl md:text-2xl text-gray-600 mb-10 animate-fade-in-up" style="animation-delay: 0.2s;">
        🌊 在这里，发现精彩文章，与志同道合的朋友一起成长
      </p>
      <div class="flex justify-center gap-4 animate-fade-in-up" style="animation-delay: 0.4s;">
        <router-link to="/article" class="px-8 py-3 bg-gradient-to-r from-sky-500 to-cyan-500 text-white rounded-full font-medium hover:shadow-xl hover:shadow-cyan-500/30 transition-all duration-300 transform hover:scale-105">
          🚀 探索文章
        </router-link>
        <router-link to="/register" class="px-8 py-3 bg-white text-sky-600 rounded-full font-medium border-2 border-sky-200 hover:border-sky-400 transition-all duration-300">
          💫 加入我们
        </router-link>
      </div>
    </div>
    <!-- 装饰性波浪 -->
    <div class="absolute bottom-0 left-0 w-full">
      <svg viewBox="0 0 1440 120" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M0 120L60 110C120 100 240 80 360 70C480 60 600 60 720 65C840 70 960 80 1080 85C1200 90 1320 90 1380 90L1440 90V120H1380C1320 120 1200 120 1080 120C960 120 840 120 720 120C600 120 480 120 360 120C240 120 120 120 60 120H0Z" fill="white" fill-opacity="0.5"/>
      </svg>
    </div>
  </section>

  <!-- Main Content -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 relative z-10">
    <div class="flex gap-8">
      <!-- Article List -->
      <div class="flex-1">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-3xl font-bold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">
            📚 最新文章
          </h2>
          <router-link to="/article" class="text-sky-600 hover:text-cyan-600 font-medium flex items-center gap-1 transition-colors">
            查看全部
            <span>→</span>
          </router-link>
        </div>

        <div v-if="loading" class="space-y-6">
          <div v-for="i in 5" :key="i" class="bg-white rounded-2xl shadow-soft p-6 animate-pulse">
            <div class="h-8 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-lg w-3/4 mb-4"></div>
            <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-1/2 mb-3"></div>
            <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-full mb-2"></div>
            <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-2/3"></div>
          </div>
        </div>

        <template v-else>
          <div v-if="articles?.length > 0" class="space-y-6">
            <article
              v-for="(article, index) in articles"
              :key="article.id"
              class="group bg-white rounded-2xl shadow-soft hover:shadow-xl transition-all duration-300 p-6 cursor-pointer transform hover:-translate-y-1 border border-gray-100 hover:border-sky-200"
              :style="{ animationDelay: `${index * 0.1}s` }"
              @click="$router.push(`/article/${article.id}`)"
            >
              <div class="flex items-start gap-4">
                <div class="flex-1">
                  <h3 class="text-xl font-bold text-gray-900 mb-3 group-hover:text-transparent group-hover:bg-gradient-to-r group-hover:from-sky-600 group-hover:to-cyan-600 group-hover:bg-clip-text transition-all duration-300">
                    {{ article.title }}
                  </h3>
                  <div class="flex items-center flex-wrap gap-3 text-sm text-gray-500 mb-4">
                    <span
                      v-if="article.category"
                      class="px-3 py-1 bg-gradient-to-r from-sky-100 to-cyan-100 text-sky-700 rounded-full font-medium"
                    >
                      🏷️ {{ article.category.name }}
                    </span>
                    <span class="flex items-center gap-1.5">
                      <div class="w-6 h-6 rounded-full bg-gradient-to-br from-sky-400 to-cyan-400 flex items-center justify-center text-white text-xs">
                        {{ article.author.name.charAt(0) }}
                      </div>
                      {{ article.author.name }}
                    </span>
                    <span class="flex items-center gap-1.5">
                      📅 {{ formatDate(article.createdAt) }}
                    </span>
                  </div>
                  <p class="text-gray-600 line-clamp-2 mb-4">{{ article.summary }}</p>
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-4 text-sm text-gray-500">
                      <span class="flex items-center gap-1.5 bg-gray-50 px-3 py-1 rounded-full">
                        👁️ {{ article.views }}
                      </span>
                      <span class="flex items-center gap-1.5 bg-gray-50 px-3 py-1 rounded-full">
                        ❤️ {{ article.likes }}
                      </span>
                    </div>
                    <div class="flex gap-2">
                      <span
                        v-for="tag in (article.tags || []).slice(0, 3)"
                        :key="tag.id"
                        class="px-2.5 py-1 bg-gradient-to-r from-sky-50 to-cyan-50 text-sky-600 rounded-full text-xs font-medium"
                      >
                        #{{ tag.name }}
                      </span>
                    </div>
                  </div>
                </div>
                <div class="hidden md:block w-24 h-24 bg-gradient-to-br from-sky-100 to-cyan-100 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
                  <span class="text-4xl">📖</span>
                </div>
              </div>
            </article>
          </div>

          <!-- Pagination -->
          <div v-if="total > 0" class="flex justify-center items-center gap-2 mt-12">
            <button
              :disabled="currentPage <= 1"
              @click="currentPage--, fetchArticles()"
              class="px-5 py-2.5 bg-white border-2 border-gray-200 rounded-full hover:border-sky-400 hover:text-sky-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300"
            >
              ← 上一页
            </button>
            <div class="px-4 py-2 bg-gradient-to-r from-sky-100 to-cyan-100 rounded-full">
              <span class="text-sky-700 font-medium">
                {{ currentPage }} / {{ Math.ceil(total / pageSize) }}
              </span>
            </div>
            <button
              :disabled="currentPage >= Math.ceil(total / pageSize)"
              @click="currentPage++, fetchArticles()"
              class="px-5 py-2.5 bg-white border-2 border-gray-200 rounded-full hover:border-sky-400 hover:text-sky-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300"
            >
              下一页 →
            </button>
          </div>

          <div v-else class="text-center py-20">
            <div class="w-24 h-24 mx-auto bg-gradient-to-br from-sky-100 to-cyan-100 rounded-full flex items-center justify-center mb-6">
              <span class="text-5xl">📝</span>
            </div>
            <p class="text-gray-500 text-xl font-medium">暂无文章</p>
            <p class="text-gray-400 mt-2">快来发布第一篇文章吧~</p>
          </div>
        </template>
      </div>

      <!-- Sidebar -->
      <aside class="w-80 hidden lg:block space-y-6">
        <!-- Tags -->
        <div class="bg-white rounded-2xl shadow-soft p-6 border border-gray-100 hover:border-sky-200 transition-colors">
          <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-8 h-8 bg-gradient-to-br from-sky-500 to-cyan-500 rounded-lg flex items-center justify-center text-white">
              🏷️
            </span>
            热门标签
          </h3>
          <div v-if="tagsLoading" class="space-y-2">
            <div v-for="i in 5" :key="i" class="h-8 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-full animate-pulse"></div>
          </div>
          <div v-else-if="tags?.length > 0" class="flex flex-wrap gap-2">
            <button
              v-for="tag in tags"
              :key="tag.id"
              @click="goToArticleByTag(tag.id)"
              class="px-4 py-2 bg-gradient-to-r from-sky-50 to-cyan-50 text-sky-700 rounded-full text-sm font-medium hover:from-sky-100 hover:to-cyan-100 hover:shadow-md transition-all duration-300 transform hover:scale-105"
            >
              #{{ tag.name }} ({{ tag.count }})
            </button>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-3xl">🏷️</span>
            <p class="mt-2">暂无标签</p>
          </div>
        </div>

        <!-- Categories -->
        <div class="bg-white rounded-2xl shadow-soft p-6 border border-gray-100 hover:border-teal-200 transition-colors">
          <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-8 h-8 bg-gradient-to-br from-cyan-500 to-teal-500 rounded-lg flex items-center justify-center text-white">
              📂
            </span>
            分类
          </h3>
          <div v-if="categoriesLoading" class="space-y-2">
            <div
              v-for="i in 5"
              :key="i"
              class="h-10 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-lg animate-pulse"
            ></div>
          </div>
          <div v-else-if="categories.length > 0" class="space-y-2">
            <button
              v-for="(category, index) in categories"
              :key="category.id"
              @click="goToArticleByCategory(category.id)"
              class="w-full text-left px-4 py-3 bg-gradient-to-r from-gray-50 to-white text-gray-700 rounded-xl hover:from-cyan-50 hover:to-teal-50 hover:text-cyan-700 transition-all duration-300 flex items-center gap-3 group"
            >
              <span class="w-2 h-2 bg-gradient-to-r from-cyan-400 to-teal-400 rounded-full group-hover:w-3 transition-all"></span>
              {{ category.name }}
            </button>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-3xl">📂</span>
            <p class="mt-2">暂无分类</p>
          </div>
        </div>

        <!-- 快捷链接 -->
        <div class="bg-gradient-to-br from-sky-500 to-cyan-500 rounded-2xl shadow-lg p-6 text-white">
          <h3 class="text-lg font-bold mb-4 flex items-center gap-2">
            ✨ 快捷链接
          </h3>
          <div class="space-y-3">
            <router-link to="/article" class="flex items-center gap-3 p-3 bg-white/20 rounded-xl hover:bg-white/30 transition-colors">
              <span class="text-xl">📝</span>
              <span>浏览全部文章</span>
            </router-link>
            <router-link to="/register" class="flex items-center gap-3 p-3 bg-white/20 rounded-xl hover:bg-white/30 transition-colors">
              <span class="text-xl">💫</span>
              <span>加入我们</span>
            </router-link>
          </div>
        </div>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles } from '@/api/article'
import { getTags } from '@/api/tag'
import { getCategories } from '@/api/category'
import type { Article, Tag, Category } from '@/api/types'

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
      pageSize: pageSize.value,
      status: 2,
      sortBy: 'createdAt',
      sortOrder: 'desc',
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
    const res = await getTags({ page: 1, pageSize: 20 })
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
    const res = await getCategories({ page: 1, pageSize: 20 })
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

<style scoped>
/* 自定义动画 */
@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

@keyframes gradient {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

@keyframes fadeInUp {
  0% { opacity: 0; transform: translateY(30px); }
  100% { opacity: 1; transform: translateY(0); }
}

@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.3); }
  50% { transform: scale(1.05); }
  70% { transform: scale(0.9); }
  100% { opacity: 1; transform: scale(1); }
}

.animate-float {
  animation: float 6s ease-in-out infinite;
}

.animate-gradient {
  background-size: 200% auto;
  animation: gradient 8s ease infinite;
}

.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
}

.animate-bounce-in {
  animation: bounceIn 0.8s ease-out forwards;
}
</style>
