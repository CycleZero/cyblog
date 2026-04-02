<template>
  <div class="category-list">
    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item>
          <el-button type="primary" @click="handleCreate">新增分类</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" row-key="id">
        <el-table-column prop="name" label="分类名称" min-width="150" />
        <el-table-column prop="slug" label="别名" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.description || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑分类' : '新增分类'"
      width="500px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" placeholder="请输入分类别名" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getCategories, createCategory, updateCategory, deleteCategory as deleteCategoryApi } from '@/api/category'
import { formatDate } from '@/utils/date'
import type { Category } from '@/api/types'

// 状态
const loading = ref(false)
const tableData = ref<Category[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

// 表单数据
const form = reactive({
  name: '',
  slug: '',
  description: '',
  sort: 0,
})

// 表单验证规则
const rules: FormRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入分类别名', trigger: 'blur' }],
}

// 获取分类列表
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const res = await getCategories({ page: 1, pageSize: 100 })
    tableData.value = res.list
  } catch {
    ElMessage.error('获取分类列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
function resetForm(): void {
  form.name = ''
  form.slug = ''
  form.description = ''
  form.sort = 0
  currentId.value = null
  formRef.value?.resetFields()
}

// 新增
function handleCreate(): void {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: Category): void {
  isEdit.value = true
  currentId.value = row.id
  form.name = row.name
  form.slug = row.slug
  form.description = row.description || ''
  form.sort = row.sort || 0
  dialogVisible.value = true
}

// 提交
async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    if (isEdit.value && currentId.value) {
      await updateCategory({ id: currentId.value, ...form })
      ElMessage.success('更新成功')
    } else {
      await createCategory({ name: form.name, slug: form.slug, description: form.description, sort: form.sort })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitLoading.value = false
  }
}

// 删除
async function handleDelete(row: Category): Promise<void> {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', { type: 'warning' })
    await deleteCategoryApi(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 生命周期
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.category-list {
  padding: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}
</style>
