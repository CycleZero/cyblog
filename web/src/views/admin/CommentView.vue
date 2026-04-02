<template>
  <div class="comment-list">
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
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="content" label="评论内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="user.name" label="评论者" width="120">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="24" :src="row.user.avatar">
                {{ row.user.name?.charAt(0) }}
              </el-avatar>
              <span>{{ row.user.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="article_id" label="文章ID" width="100" />
        <el-table-column prop="likes" label="点赞" width="80" />
        <el-table-column prop="created_at" label="评论时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="goToArticle(row.article_id)">
              查看文章
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
import { getComments, deleteComment } from '@/api/comment'
import { formatDateTime } from '@/utils/date'
import type { Comment } from '@/api/types'

const router = useRouter()

// 查询参数
interface QueryParams {
  keyword?: string
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
const tableData = ref<Comment[]>([])
const queryParams = reactive<QueryParams>({
  page: 1,
  pageSize: 10,
})
const pagination = reactive<Pagination>({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 获取评论列表
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    // 由于管理端评论接口可能不存在，先用普通评论接口
    const res = await getComments({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
    })
    tableData.value = res.list
    pagination.total = res.total
  } catch {
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
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
  handleSearch()
}

// 删除
async function handleDelete(row: Comment): Promise<void> {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await deleteComment(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 跳转文章
function goToArticle(articleId: number): void {
  router.push(`/article/${articleId}`)
}

// 分页大小变化
function handleSizeChange(size: number): void {
  pagination.pageSize = size
  fetchData()
}

// 页码变化
function handlePageChange(page: number): void {
  pagination.page = page
  fetchData()
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.comment-list {
  padding: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
