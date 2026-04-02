<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon article-icon">
              <el-icon :size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.articleCount }}</div>
              <div class="stat-label">文章总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon view-icon">
              <el-icon :size="32"><View /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.todayViews }}</div>
              <div class="stat-label">今日访问</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>近期文章</span>
              <el-button type="primary" link @click="goToArticle">
                查看全部
              </el-button>
            </div>
          </template>
          <el-table :data="stats.recentArticles" v-loading="loading">
            <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="hot-articles-card">
          <template #header>
            <div class="card-header">
              <span>热门文章</span>
            </div>
          </template>
          <div v-loading="loading" class="hot-list">
            <div
              v-for="(article, index) in stats.hotArticles"
              :key="article.id"
              class="hot-item"
              @click="goToArticleDetail(article.id)"
            >
              <span class="hot-rank" :class="`rank-${index + 1}`">{{ index + 1 }}</span>
              <span class="hot-title">{{ article.title }}</span>
              <span class="hot-views">{{ article.views }} 阅读</span>
            </div>
            <el-empty v-if="!loading && stats.hotArticles.length === 0" description="暂无数据" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Document,
  View,
} from '@element-plus/icons-vue'
import { getDashboard, ArticleStatusMap, ArticleStatusTypeMap } from '@/api/admin'

const router = useRouter()

const loading = ref(false)
const stats = reactive({
  articleCount: 0,
  todayViews: 0,
  recentArticles: [] as { id: number; title: string; status: number }[],
  hotArticles: [] as { id: number; title: string; views: number }[],
})

async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const data = await getDashboard()
    stats.articleCount = data.articleCount
    stats.todayViews = data.todayViews
    stats.recentArticles = data.recentArticles
    stats.hotArticles = data.hotArticles
  } catch {
    ElMessage.error('获取仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

function getStatusType(status: number): string {
  return ArticleStatusTypeMap[status] || 'info'
}

function getStatusText(status: number): string {
  return ArticleStatusMap[status] || '未知'
}

function goToArticle(): void {
  router.push('/admin/article')
}

function goToArticleDetail(id: number): void {
  router.push(`/article/${id}`)
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  height: 120px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.article-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.view-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: var(--el-text-color-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
}

.content-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.hot-articles-card {
  height: 100%;
}

.hot-list {
  min-height: 200px;
}

.hot-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
  cursor: pointer;
  transition: background-color 0.3s;
}

.hot-item:last-child {
  border-bottom: none;
}

.hot-item:hover {
  background-color: var(--el-fill-color-light);
}

.hot-rank {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 12px;
  color: #fff;
  background-color: var(--el-color-info-light);
  color: var(--el-color-info);
}

.hot-rank.rank-1 {
  background-color: #ffd700;
  color: #fff;
}

.hot-rank.rank-2 {
  background-color: #c0c0c0;
  color: #fff;
}

.hot-rank.rank-3 {
  background-color: #cd7f32;
  color: #fff;
}

.hot-title {
  flex: 1;
  font-size: 14px;
  color: var(--el-text-color-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hot-views {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
}
</style>
