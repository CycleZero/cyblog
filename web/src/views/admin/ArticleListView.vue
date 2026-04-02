<template>
  <div class="article-list">
    <el-card class="filter-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="关键词">
          <el-input
            v-model="queryParams.keyword"
            placeholder="请输入关键词"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="全部" :value="undefined" />
            <el-option label="草稿" :value="1" />
            <el-option label="已发布" :value="2" />
            <el-option label="私密" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <template #header>
        <div class="table-header">
          <span>文章列表</span>
          <div class="header-actions">
            <el-button
              v-if="selectedRows.length > 0"
              type="danger"
              @click="handleBatchDelete"
            >
              批量删除 ({{ selectedRows.length }})
            </el-button>
            <el-button type="primary" @click="handleCreate">新增文章</el-button>
          </div>
        </div>
      </template>

      <el-table
        :data="tableData"
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <el-link type="primary" @click="goToDetail(row.id)">
              {{ row.title }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="author?.name" label="作者" width="120" />
        <el-table-column prop="category?.name" label="分类" width="120">
          <template #default="{ row }">
            {{ row.category?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_top" label="置顶" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_top" type="warning">是</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览" width="80" />
        <el-table-column prop="likes" label="点赞" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="handleToggleTop(row)">
              {{ row.is_top ? '取消置顶' : '置顶' }}
            </el-button>
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdminArticles,
  setArticleTop,
  batchDeleteArticles,
  // deleteArticle as deleteArticleApi,
} from '@/api/admin'
import { getCategories } from '@/api/category'
import { ArticleStatusMap, ArticleStatusTypeMap } from '@/api/admin'
import { formatDate } from '@/utils/date'
import type { Article, Category } from '@/api/types'

const router = useRouter()

// 查询参数
interface QueryParams {
  keyword?: string
  category_id?: number
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
const categories = ref<Category[]>([])
const queryParams = reactive<QueryParams>({
  page: 1,
  pageSize: 10,
})
const pagination = reactive<Pagination>({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 获取文章列表
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const res = await getAdminArticles({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
      keyword: queryParams.keyword,
      category_id: queryParams.category_id,
      status: queryParams.status,
    })
    tableData.value = res.list
    pagination.total = res.total
  } catch {
    ElMessage.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分类列表
async function fetchCategories(): Promise<void> {
  try {
    const res = await getCategories({ page: 1, pageSize: 100 })
    categories.value = res.list
  } catch {
    ElMessage.error('获取分类列表失败')
  }
}

// 搜索
function handleSearch(): void {
  pagination.page = 1
  fetchData()
}

// 重置
function handleReset(): void {
  queryParams.keyword = undefined
  queryParams.category_id = undefined
  queryParams.status = undefined
  handleSearch()
}

// 选择变化
function handleSelectionChange(rows: Article[]): void {
  selectedRows.value = rows
}

// 新增
function handleCreate(): void {
  router.push('/admin/article/create')
}

// 编辑
function handleEdit(row: Article): void {
  router.push(`/admin/article/${row.id}/edit`)
}

// 删除
async function handleDelete(row: Article): Promise<void> {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning',
    })
    await batchDeleteArticles([row.id])
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
async function handleBatchDelete(): Promise<void> {
  const ids = selectedRows.value.map((row) => row.id)
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${ids.length} 篇文章吗？`, '提示', {
      type: 'warning',
    })
    await batchDeleteArticles(ids)
    ElMessage.success('批量删除成功')
    selectedRows.value = []
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 切换置顶
async function handleToggleTop(row: Article): Promise<void> {
  try {
    await ElMessageBox.confirm(
      row.is_top ? '确定要取消置顶吗？' : '确定要置顶吗？',
      '提示',
      { type: 'warning' }
    )
    await setArticleTop(row.id, !row.is_top)
    ElMessage.success(row.is_top ? '已取消置顶' : '已置顶')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

// 分页大小变化
function handleSizeChange(size: number): void {
  pagination.pageSize = size
  queryParams.pageSize = size
  fetchData()
}

// 页码变化
function handlePageChange(page: number): void {
  pagination.page = page
  queryParams.page = page
  fetchData()
}

// 跳转详情
function goToDetail(id: number): void {
  router.push(`/article/${id}`)
}

// 获取状态类型
function getStatusType(status: number): string {
  return ArticleStatusTypeMap[status] || 'info'
}

// 获取状态文本
function getStatusText(status: number): string {
  return ArticleStatusMap[status] || '未知'
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.article-list {
  padding: 0;
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

.header-actions {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
