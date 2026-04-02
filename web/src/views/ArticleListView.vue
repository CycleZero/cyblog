<template>
  <!-- Main Content -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 relative z-10">
    <div class="flex gap-8">
      <!-- Article List -->
      <div class="flex-1">
        <!-- 搜索和筛选 -->
        <div class="bg-white rounded-2xl shadow-soft p-6 mb-8 border border-gray-100">
          <div class="flex flex-col md:flex-row gap-4">
            <div class="flex-1 relative">
              <input
                v-model="searchKeyword"
                type="text"
                placeholder="搜索文章..."
                class="w-full px-4 py-3 pl-12 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent transition-all"
                @keyup.enter="handleSearch"
              />
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">🔍</span>
            </div>
            <select
              v-model="selectedCategory"
              class="px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent cursor-pointer"
              @change="handleCategoryChange"
            >
              <option value="">全部分类</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
            <select
              v-model="selectedTag"
              class="px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent cursor-pointer"
              @change="handleTagChange"
            >
              <option value="">全部标签</option>
              <option v-for="tag in tags" :key="tag.id" :value="tag.id">
                {{ tag.name }}
              </option>
            </select>
          </div>
        </div>

        <div class="flex items-center justify-between mb-6">
          <h2 class="text-2xl font-bold text-gray-800">
            📚 文章列表
            <span class="text-sm font-normal text-gray-500 ml-2">共 {{ total }} 篇</span>
          </h2>
          <div class="flex items-center gap-2">
            <button
              @click="sortOrder = 'desc'"
              :class="[
                'px-4 py-2 rounded-lg font-medium transition-all',
                sortOrder === 'desc' ? 'bg-sky-500 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              ]"
            >
              最新
            </button>
            <button
              @click="sortOrder = 'asc'"
              :class="[
                'px-4 py-2 rounded-lg font-medium transition-all',
                sortOrder === 'asc' ? 'bg-sky-500 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              ]"
            >
              最旧
            </button>
          </div>
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
                  <div class="flex items-center gap-2 mb-3">
                    <span
                      v-if="article.is_top"
                      class="px-2 py-0.5 bg-gradient-to-r from-red-500 to-orange-500 text-white text-xs rounded-full font-medium"
                    >
                      置顶
                    </span>
                    <span
                      v-if="article.is_original"
                      class="px-2 py-0.5 bg-gradient-to-r from-sky-500 to-cyan-500 text-white text-xs rounded-full font-medium"
                    >
                      原创
                    </span>
                  </div>
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
                      📅 {{ formatDate(article.created_at) }}
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
                <div class="hidden md:block w-32 h-32 bg-gradient-to-br from-sky-100 to-cyan-100 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300 overflow-hidden">
                  <img v-if="article.cover_image" :src="article.cover_image" class="w-full h-full object-cover" alt="cover" />
                  <span v-else class="text-4xl">📖</span>
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
            <div class="flex items-center gap-2">
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="currentPage = page, fetchArticles()"
                :class="[
                  'w-10 h-10 rounded-full transition-all duration-300',
                  currentPage === page
                    ? 'bg-gradient-to-r from-sky-500 to-cyan-500 text-white shadow-lg'
                    : 'bg-white border-2 border-gray-200 text-gray-600 hover:border-sky-400'
                ]"
              >
                {{ page }}
              </button>
            </div>
            <button
              :disabled="currentPage >= totalPages"
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
            <p class="text-gray-400 mt-2">换个关键词试试吧~</p>
          </div>
        </template>
      </div>

      <!-- Sidebar -->
      <aside class="w-80 hidden lg:block space-y-6">
        <!-- 热门标签 -->
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
              @click="selectTag(tag.id)"
              :class="[
                'px-4 py-2 rounded-full text-sm font-medium transition-all duration-300 transform hover:scale-105',
                selectedTag === tag.id
                  ? 'bg-gradient-to-r from-sky-500 to-cyan-500 text-white shadow-md'
                  : 'bg-gradient-to-r from-sky-50 to-cyan-50 text-sky-700 hover:from-sky-100 hover:to-cyan-100'
              ]"
            >
              #{{ tag.name }} ({{ tag.count }})
            </button>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-3xl">🏷️</span>
            <p class="mt-2">暂无标签</p>
          </div>
        </div>

        <!-- 分类 -->
        <div class="bg-white rounded-2xl shadow-soft p-6 border border-gray-100 hover:border-teal-200 transition-colors">
          <h3 class="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
            <span class="w-8 h-8 bg-gradient-to-br from-cyan-500 to-teal-500 rounded-lg flex items-center justify-center text-white">
              📂
            </span>
            分类
          </h3>
          <div v-if="categoriesLoading" class="space-y-2">
            <div v-for="i in 5" :key="i" class="h-10 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-lg animate-pulse"></div>
          </div>
          <div v-else-if="categories.length > 0" class="space-y-2">
            <button
              v-for="category in categories"
              :key="category.id"
              @click="selectCategory(category.id)"
              :class="[
                'w-full text-left px-4 py-3 rounded-xl transition-all duration-300 flex items-center gap-3 group',
                selectedCategory === category.id
                  ? 'bg-gradient-to-r from-cyan-50 to-teal-50 text-cyan-700'
                  : 'bg-gradient-to-r from-gray-50 to-white text-gray-700 hover:from-cyan-50 hover:to-teal-50'
              ]"
            >
              <span :class="[
                'w-2 h-2 rounded-full transition-all',
                selectedCategory === category.id ? 'bg-gradient-to-r from-cyan-400 to-teal-400 w-3' : 'bg-gradient-to-r from-cyan-400 to-teal-400'
              ]"></span>
              {{ category.name }}
            </button>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <span class="text-3xl">📂</span>
            <p class="mt-2">暂无分类</p>
          </div>
        </div>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles } from '@/api/article'
import { getTags } from '@/api/tag'
import { getCategories } from '@/api/category'
import { formatDate } from '@/utils/date'
import type { Article, Tag, Category } from '@/api/types'

const route = useRoute()

const articles = ref<Article[]>([])
const tags = ref<Tag[]>([])
const categories = ref<Category[]>([])
const loading = ref(true)
const tagsLoading = ref(true)
const categoriesLoading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const selectedCategory = ref<number | ''>('')
const selectedTag = ref<number | ''>('')
const sortOrder = ref('desc')

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value
  let start = Math.max(1, current - 2)
  let end = Math.min(total, current + 2)

  if (end - start < 4) {
    if (start === 1) {
      end = Math.min(total, start + 4)
    } else {
      start = Math.max(1, end - 4)
    }
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

async function fetchArticles() {
  try {
    loading.value = true
    const res = await getArticles({
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value || undefined,
      category_id: selectedCategory.value || undefined,
      tag_id: selectedTag.value || undefined,
      status: 2,
      sort_by: 'created_at',
      sort_order: sortOrder.value,
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
    const res = await getTags({ page: 1, pageSize: 50 })
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
    const res = await getCategories({ page: 1, pageSize: 50 })
    categories.value = res.list
  } catch (error) {
    console.error('获取分类失败', error)
  } finally {
    categoriesLoading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchArticles()
}

function handleCategoryChange() {
  currentPage.value = 1
  fetchArticles()
}

function handleTagChange() {
  currentPage.value = 1
  fetchArticles()
}

function selectCategory(categoryId: number) {
  if (selectedCategory.value === categoryId) {
    selectedCategory.value = ''
  } else {
    selectedCategory.value = categoryId
  }
  currentPage.value = 1
  fetchArticles()
}

function selectTag(tagId: number) {
  if (selectedTag.value === tagId) {
    selectedTag.value = ''
  } else {
    selectedTag.value = tagId
  }
  currentPage.value = 1
  fetchArticles()
}

watch(
  () => route.query,
  (query) => {
    if (query.category_id) {
      selectedCategory.value = Number(query.category_id)
    }
    if (query.tag_id) {
      selectedTag.value = Number(query.tag_id)
    }
    fetchArticles()
  },
  { immediate: true }
)

onMounted(() => {
  if (route.query.category_id) {
    selectedCategory.value = Number(route.query.category_id)
  }
  if (route.query.tag_id) {
    selectedTag.value = Number(route.query.tag_id)
  }

  fetchArticles()
  fetchTags()
  fetchCategories()
})
</script>
