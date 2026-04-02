<template>
  <!-- Loading State -->
  <div v-if="loading" class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 relative z-10">
    <div class="bg-white rounded-2xl shadow-soft p-8 animate-pulse">
      <div class="h-12 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-lg w-3/4 mb-6"></div>
      <div class="flex items-center gap-4 mb-8">
        <div class="h-10 w-10 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded-full"></div>
        <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-32"></div>
      </div>
      <div class="space-y-4">
        <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-full"></div>
        <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-full"></div>
        <div class="h-4 bg-gradient-to-r from-gray-200 via-gray-100 to-gray-200 rounded w-2/3"></div>
      </div>
    </div>
  </div>

  <!-- Article Content -->
  <main v-else-if="article" class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 relative z-10">
    <!-- Article Header -->
    <article class="bg-white rounded-2xl shadow-soft overflow-hidden mb-8">
      <!-- Cover Image -->
      <div v-if="article.cover_image" class="relative h-64 md:h-96 overflow-hidden">
        <img :src="article.cover_image" :alt="article.title" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
      </div>

      <div class="p-8">
        <!-- Tags and Badges -->
        <div class="flex flex-wrap items-center gap-2 mb-4">
          <span
            v-if="article.is_top"
            class="px-3 py-1 bg-gradient-to-r from-red-500 to-orange-500 text-white text-sm rounded-full font-medium"
          >
            🔝 置顶
          </span>
          <span
            v-if="article.is_original"
            class="px-3 py-1 bg-gradient-to-r from-sky-500 to-cyan-500 text-white text-sm rounded-full font-medium"
          >
            ✍️ 原创
          </span>
          <span
            v-if="article.category"
            class="px-3 py-1 bg-gradient-to-r from-sky-100 to-cyan-100 text-sky-700 rounded-full text-sm font-medium"
          >
            🏷️ {{ article.category.name }}
          </span>
        </div>

        <!-- Title -->
        <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
          {{ article.title }}
        </h1>

        <!-- Author Info -->
        <div class="flex items-center justify-between flex-wrap gap-4 pb-6 border-b border-gray-100">
          <div class="flex items-center gap-3">
            <div class="w-12 h-12 rounded-full bg-gradient-to-br from-sky-400 to-cyan-400 flex items-center justify-center text-white text-lg font-bold">
              {{ article.author.name.charAt(0) }}
            </div>
            <div>
              <div class="font-medium text-gray-900">{{ article.author.name }}</div>
              <div class="text-sm text-gray-500">{{ formatDate(article.created_at) }}</div>
            </div>
          </div>
          <div class="flex items-center gap-4 text-sm text-gray-500">
            <span class="flex items-center gap-1.5 bg-gray-50 px-3 py-1.5 rounded-full">
              👁️ {{ article.views }} 阅读
            </span>
            <span class="flex items-center gap-1.5 bg-gray-50 px-3 py-1.5 rounded-full">
              ❤️ {{ article.likes }} 点赞
            </span>
          </div>
        </div>

        <!-- Article Tags -->
        <div v-if="article.tags && article.tags.length > 0" class="flex flex-wrap gap-2 mt-6">
          <span
            v-for="tag in article.tags"
            :key="tag.id"
            class="px-3 py-1.5 bg-gradient-to-r from-sky-50 to-cyan-50 text-sky-600 rounded-full text-sm font-medium hover:from-sky-100 hover:to-cyan-100 transition-colors cursor-pointer"
            @click="$router.push(`/article?tag_id=${tag.id}`)"
          >
            #{{ tag.name }}
          </span>
        </div>
      </div>
    </article>

    <!-- Article Body -->
    <article class="bg-white rounded-2xl shadow-soft p-8 mb-8">
      <div class="prose prose-lg max-w-none">
        <div v-html="renderedContent" class="article-content"></div>
      </div>
    </article>

    <!-- Like Section -->
    <div class="bg-white rounded-2xl shadow-soft p-6 mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-bold text-gray-900 mb-1">喜欢这篇文章吗？</h3>
          <p class="text-gray-500 text-sm">如果觉得有帮助，就点个赞吧！</p>
        </div>
        <button
          @click="handleLike"
          :class="[
            'px-6 py-3 rounded-full font-medium transition-all duration-300 transform hover:scale-105',
            isLiked
              ? 'bg-gradient-to-r from-red-500 to-pink-500 text-white shadow-lg'
              : 'bg-gradient-to-r from-sky-500 to-cyan-500 text-white hover:shadow-lg hover:shadow-cyan-500/30'
          ]"
        >
          {{ isLiked ? '❤️ 已点赞' : '🤍 点赞' }}
        </button>
      </div>
    </div>

    <!-- Author Card -->
    <div class="bg-gradient-to-br from-sky-500 to-cyan-500 rounded-2xl shadow-lg p-6 text-white mb-8">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-white/20 flex items-center justify-center text-2xl font-bold">
          {{ article.author.name.charAt(0) }}
        </div>
        <div class="flex-1">
          <div class="text-xl font-bold mb-1">{{ article.author.name }}</div>
          <div class="text-white/80 text-sm">作者</div>
        </div>
        <button class="px-4 py-2 bg-white/20 rounded-full hover:bg-white/30 transition-colors text-sm font-medium">
          查看主页
        </button>
      </div>
    </div>

    <!-- Related Articles -->
    <div v-if="relatedArticles.length > 0" class="mb-8">
      <h3 class="text-2xl font-bold text-gray-900 mb-6">📚 相关文章</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <article
          v-for="related in relatedArticles"
          :key="related.id"
          class="group bg-white rounded-2xl shadow-soft hover:shadow-xl transition-all duration-300 p-6 cursor-pointer transform hover:-translate-y-1 border border-gray-100 hover:border-sky-200"
          @click="$router.push(`/article/${related.id}`)"
        >
          <h4 class="text-lg font-bold text-gray-900 mb-2 group-hover:text-transparent group-hover:bg-gradient-to-r group-hover:from-sky-600 group-hover:to-cyan-600 group-hover:bg-clip-text transition-all duration-300">
            {{ related.title }}
          </h4>
          <p class="text-gray-600 text-sm line-clamp-2 mb-3">{{ related.summary }}</p>
          <div class="flex items-center gap-4 text-xs text-gray-500">
            <span>👁️ {{ related.views }}</span>
            <span>❤️ {{ related.likes }}</span>
            <span>📅 {{ formatDate(related.created_at) }}</span>
          </div>
        </article>
      </div>
    </div>

    <!-- Back Button -->
    <div class="text-center">
      <router-link
        to="/article"
        class="inline-block px-6 py-3 bg-white border-2 border-gray-200 rounded-full text-gray-700 font-medium hover:border-sky-400 hover:text-sky-600 transition-all duration-300"
      >
        ← 返回文章列表
      </router-link>
    </div>
  </main>

  <!-- Not Found -->
  <div v-else class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-20 relative z-10 text-center">
    <div class="w-24 h-24 mx-auto bg-gradient-to-br from-sky-100 to-cyan-100 rounded-full flex items-center justify-center mb-6">
      <span class="text-5xl">😢</span>
    </div>
    <h2 class="text-2xl font-bold text-gray-900 mb-2">文章不存在</h2>
    <p class="text-gray-500 mb-6">该文章可能已被删除或不存在</p>
    <router-link
      to="/article"
      class="inline-block px-6 py-3 bg-gradient-to-r from-sky-500 to-cyan-500 text-white rounded-full font-medium hover:shadow-lg hover:shadow-cyan-500/30 transition-all duration-300"
    >
      ← 返回文章列表
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import { getArticle, likeArticle, unlikeArticle } from '@/api/article'
import type { Article } from '@/api/types'

const route = useRoute()
const article = ref<Article | null>(null)
const loading = ref(true)
const isLiked = ref(false)
const relatedArticles = ref<Article[]>([])

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true,
})

// 将 Markdown 转换为 HTML
const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  return marked(article.value.content)
})

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

async function fetchArticle() {
  try {
    loading.value = true
    const id = Number(route.params.id)
    const res = await getArticle(id, true)
    article.value = res
    relatedArticles.value = []
  } catch (error) {
    console.error('获取文章详情失败', error)
    article.value = null
  } finally {
    loading.value = false
  }
}

async function handleLike() {
  if (!article.value) return

  try {
    if (isLiked.value) {
      await unlikeArticle(article.value.id)
      article.value.likes--
      isLiked.value = false
    } else {
      await likeArticle(article.value.id)
      article.value.likes++
      isLiked.value = true
    }
  } catch (error) {
    console.error('点赞操作失败', error)
  }
}

onMounted(() => {
  fetchArticle()
})
</script>

<style scoped>
/* Article Content Styles */
.article-content {
  line-height: 1.8;
  color: #374151;
}

.article-content :deep(h1) {
  font-size: 2rem;
  font-weight: 700;
  margin-top: 2rem;
  margin-bottom: 1rem;
  color: #111827;
}

.article-content :deep(h2) {
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 1.75rem;
  margin-bottom: 0.75rem;
  color: #111827;
}

.article-content :deep(h3) {
  font-size: 1.25rem;
  font-weight: 600;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
  color: #111827;
}

.article-content :deep(p) {
  margin-bottom: 1.25rem;
}

.article-content :deep(a) {
  color: #0ea5e9;
  text-decoration: underline;
}

.article-content :deep(a:hover) {
  color: #0284c7;
}

.article-content :deep(code) {
  background-color: #f3f4f6;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
  color: #ec4899;
}

.article-content :deep(pre) {
  background-color: #1f2937;
  color: #f9fafb;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  margin: 1.5rem 0;
}

.article-content :deep(pre code) {
  background-color: transparent;
  color: inherit;
  padding: 0;
}

.article-content :deep(blockquote) {
  border-left: 4px solid #0ea5e9;
  padding-left: 1rem;
  margin: 1.5rem 0;
  font-style: italic;
  color: #6b7280;
}

.article-content :deep(ul),
.article-content :deep(ol) {
  margin: 1rem 0;
  padding-left: 1.5rem;
}

.article-content :deep(li) {
  margin-bottom: 0.5rem;
}

.article-content :deep(img) {
  max-width: 100%;
  border-radius: 0.5rem;
  margin: 1.5rem 0;
}

.article-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5rem 0;
}

.article-content :deep(th),
.article-content :deep(td) {
  border: 1px solid #e5e7eb;
  padding: 0.75rem;
  text-align: left;
}

.article-content :deep(th) {
  background-color: #f9fafb;
  font-weight: 600;
}
</style>
