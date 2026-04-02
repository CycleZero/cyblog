<template>
  <div class="user-list">
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
        <el-form-item label="角色">
          <el-select v-model="queryParams.role" placeholder="请选择角色" clearable>
            <el-option label="全部" :value="undefined" />
            <el-option label="管理员" value="admin" />
            <el-option label="编辑" value="editor" />
            <el-option label="用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" width="200">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="40" :src="row.avatar">
                {{ row.name?.charAt(0) }}
              </el-avatar>
              <div class="user-info">
                <div class="name">{{ row.name }}</div>
                <div class="email">{{ row.email }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-dropdown trigger="click" @command="(cmd: string) => handleCommand(cmd, row)">
              <el-button link type="primary">
                更多 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑角色</el-dropdown-item>
                  <el-dropdown-item
                    :command="row.status === 1 ? 'disable' : 'enable'"
                  >
                    {{ row.status === 1 ? '禁用用户' : '启用用户' }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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

    <!-- 编辑角色对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑用户角色" width="400px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.name" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="editForm.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="编辑" value="editor" />
            <el-option label="用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleUpdateRole">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import {
  getAdminUsers,
  updateUserRole,
  updateUserStatus,
  UserRoleMap,
  UserRoleTypeMap,
  type AdminUser,
} from '@/api/admin'
import { formatDate } from '@/utils/date'

// 查询参数
interface QueryParams {
  keyword?: string
  role?: string
  page?: number
  pageSize?: number
}

interface Pagination {
  page: number
  pageSize: number
  total: number
}

// 状态
const loading = ref(false)
const tableData = ref<AdminUser[]>([])
const queryParams = reactive<QueryParams>({
  page: 1,
  pageSize: 10,
})
const pagination = reactive<Pagination>({
  page: 1,
  pageSize: 10,
  total: 0,
})
const editDialogVisible = ref(false)
const submitLoading = ref(false)
const currentUser = ref<AdminUser | null>(null)
const editForm = reactive({
  name: '',
  role: 'user',
})

// 获取用户列表
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const res = await getAdminUsers({
      page: queryParams.page,
      pageSize: queryParams.pageSize,
      keyword: queryParams.keyword,
      role: queryParams.role,
    })
    tableData.value = res.list
    pagination.total = res.total
  } catch {
    ElMessage.error('获取用户列表失败')
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
  queryParams.role = undefined
  handleSearch()
}

// 操作处理
async function handleCommand(command: string, row: AdminUser): Promise<void> {
  currentUser.value = row

  switch (command) {
    case 'edit':
      editForm.name = row.name
      editForm.role = row.role
      editDialogVisible.value = true
      break
    case 'enable':
    case 'disable':
      await handleToggleStatus(row)
      break
  }
}

// 更新角色
async function handleUpdateRole(): Promise<void> {
  if (!currentUser.value) return

  submitLoading.value = true
  try {
    await updateUserRole(currentUser.value.id, editForm.role)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchData()
  } catch {
    ElMessage.error('更新失败')
  } finally {
    submitLoading.value = false
  }
}

// 切换状态
async function handleToggleStatus(row: AdminUser): Promise<void> {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(`确定要${action}用户 ${row.name} 吗？`, '提示', { type: 'warning' })
    await updateUserStatus(row.id, newStatus)
    ElMessage.success(`${action}成功`)
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action}失败`)
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

// 获取角色类型
function getRoleType(role: string): string {
  return UserRoleTypeMap[role] || 'info'
}

// 获取角色文本
function getRoleText(role: string): string {
  return UserRoleMap[role] || '未知'
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.user-list {
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
  gap: 12px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-info .name {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.user-info .email {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
