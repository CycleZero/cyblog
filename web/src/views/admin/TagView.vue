<template>
  <div class="tag-list">
    <el-card class="filter-card">
      <el-form :inline="true">
        <el-form-item>
          <el-button type="primary" @click="handleCreate">新增标签</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <el-row :gutter="20">
        <el-col v-for="tag in tableData" :key="tag.id" :span="6">
          <div class="tag-item">
            <div class="tag-content">
              <el-tag :color="tag.color || '#409eff'" effect="dark">
                {{ tag.name }}
              </el-tag>
              <span class="tag-slug">{{ tag.slug }}</span>
              <span class="tag-count">{{ tag.count || 0 }} 篇文章</span>
            </div>
            <div class="tag-actions">
              <el-button link type="primary" @click="handleEdit(tag)">编辑</el-button>
              <el-button link type="danger" @click="handleDelete(tag)">删除</el-button>
            </div>
          </div>
        </el-col>
      </el-row>
      <el-empty v-if="!loading && tableData.length === 0" description="暂无数据" />
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '新增标签'"
      width="500px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="别名" prop="slug">
          <el-input v-model="form.slug" placeholder="请输入标签别名" />
        </el-form-item>
        <el-form-item label="颜色" prop="color">
          <el-color-picker v-model="form.color" />
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
import { getTags, createTag, updateTag, deleteTag as deleteTagApi } from '@/api/tag'
import type { Tag } from '@/api/types'

// 状态
const loading = ref(false)
const tableData = ref<Tag[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

// 表单数据
const form = reactive({
  name: '',
  slug: '',
  color: '#409eff',
})

// 表单验证规则
const rules: FormRules = {
  name: [{ required: true, message: '请输入标签名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标签别名', trigger: 'blur' }],
}

// 获取标签列表
async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const res = await getTags({ page: 1, pageSize: 100 })
    tableData.value = res.list
  } catch {
    ElMessage.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
function resetForm(): void {
  form.name = ''
  form.slug = ''
  form.color = '#409eff'
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
function handleEdit(row: Tag): void {
  isEdit.value = true
  currentId.value = row.id
  form.name = row.name
  form.slug = row.slug
  form.color = row.color || '#409eff'
  dialogVisible.value = true
}

// 提交
async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    if (isEdit.value && currentId.value) {
      await updateTag({ id: currentId.value, ...form })
      ElMessage.success('更新成功')
    } else {
      await createTag({ name: form.name, slug: form.slug, color: form.color })
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
async function handleDelete(row: Tag): Promise<void> {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', { type: 'warning' })
    await deleteTagApi(row.id)
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
.tag-list {
  padding: 0;
}

.filter-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
  min-height: 400px;
}

.tag-item {
  padding: 16px;
  margin-bottom: 16px;
  background-color: var(--el-fill-color-light);
  border-radius: 8px;
  transition: all 0.3s;
}

.tag-item:hover {
  background-color: var(--el-fill-color);
  transform: translateY(-2px);
}

.tag-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.tag-slug {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.tag-count {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.tag-actions {
  display: flex;
  gap: 8px;
}
</style>
