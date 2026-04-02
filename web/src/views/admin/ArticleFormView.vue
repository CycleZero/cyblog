<template>
  <div class="article-form">
    <el-card>
      <template #header>
        <span>{{ isEdit ? '编辑文章' : '创建文章' }}</span>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>
        <el-form-item label="封面图" prop="coverImage">
          <el-input v-model="form.coverImage" placeholder="请输入封面图URL" />
        </el-form-item>
        <el-form-item label="摘要" prop="summary">
          <el-input
            v-model="form.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
          />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="15"
            placeholder="请输入文章内容"
          />
        </el-form-item>
        <el-form-item label="分类" prop="categoryId">
          <el-select v-model="form.categoryId" placeholder="请选择分类">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="标签" prop="tagIds">
          <el-select v-model="form.tagIds" multiple placeholder="请选择标签">
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">草稿</el-radio>
            <el-radio :label="2">已发布</el-radio>
            <el-radio :label="3">私密</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="置顶">
          <el-switch v-model="form.isTop" />
        </el-form-item>
        <el-form-item label="原创">
          <el-switch v-model="form.isOriginal" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
            {{ isEdit ? '保存' : '发布' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getArticle, createArticle as createArticleApi, updateArticle as updateArticleApi } from '@/api/article'
import { getCategories } from '@/api/category'
import { getTags } from '@/api/tag'
import type { Article, Category, Tag } from '@/api/types'

const route = useRoute()
const router = useRouter()

const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

const isEdit = computed(() => !!route.params.id)
const articleId = computed(() => route.params.id as string | undefined)

// 表单数据
const form = reactive({
  title: '',
  coverImage: '',
  summary: '',
  content: '',
  categoryId: undefined as number | undefined,
  tagIds: [] as number[],
  status: 1,
  isTop: false,
  isOriginal: false,
})

// 表单验证规则
const rules: FormRules = {
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入文章内容', trigger: 'blur' }],
}

// 获取文章详情
async function fetchArticle(): Promise<void> {
  if (!articleId.value) return

  try {
    const data = await getArticle(Number(articleId.value))
    form.title = data.title
    form.coverImage = data.coverImage || ''
    form.summary = data.summary || ''
    form.content = data.content || ''
    form.categoryId = data.category?.id
    form.tagIds = data.tags?.map((t) => t.id) || []
    form.status = data.status
    form.isTop = data.isTop || false
    form.isOriginal = data.isOriginal || false
  } catch {
    ElMessage.error('获取文章详情失败')
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

// 获取标签列表
async function fetchTags(): Promise<void> {
  try {
    const res = await getTags({ page: 1, pageSize: 100 })
    tags.value = res.list
  } catch {
    ElMessage.error('获取标签列表失败')
  }
}

// 提交
async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const data = {
      title: form.title,
      content: form.content,
      cover_image: form.coverImage || undefined,
      summary: form.summary || undefined,
      category_id: form.categoryId,
      tag_ids: form.tagIds,
      status: form.status,
      is_top: form.isTop,
      is_original: form.isOriginal,
    }

    if (isEdit.value && articleId.value) {
      await updateArticleApi({ id: Number(articleId.value), ...data })
      ElMessage.success('更新成功')
    } else {
      await createArticleApi(data)
      ElMessage.success('创建成功')
    }
    router.push('/admin/article')
  } catch {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitLoading.value = false
  }
}

// 取消
function handleCancel(): void {
  router.back()
}

// 生命周期
onMounted(() => {
  fetchCategories()
  fetchTags()
  if (isEdit.value) {
    fetchArticle()
  }
})
</script>

<style scoped>
.article-form {
  padding: 0;
}
</style>
