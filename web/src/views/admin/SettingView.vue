<template>
  <div class="setting-view">
    <el-card>
      <template #header>
        <span>系统设置</span>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-divider content-position="left">基本信息</el-divider>
        <el-form-item label="网站名称" prop="siteName">
          <el-input v-model="form.siteName" placeholder="请输入网站名称" />
        </el-form-item>
        <el-form-item label="网站描述" prop="siteDescription">
          <el-input
            v-model="form.siteDescription"
            type="textarea"
            :rows="3"
            placeholder="请输入网站描述"
          />
        </el-form-item>

        <el-divider content-position="left">SEO 设置</el-divider>
        <el-form-item label="SEO 关键词" prop="seoKeywords">
          <el-input
            v-model="form.seoKeywords"
            type="textarea"
            :rows="2"
            placeholder="请输入 SEO 关键词，用逗号分隔"
          />
        </el-form-item>
        <el-form-item label="SEO 描述" prop="seoDescription">
          <el-input
            v-model="form.seoDescription"
            type="textarea"
            :rows="3"
            placeholder="请输入 SEO 描述"
          />
        </el-form-item>

        <el-divider content-position="left">社交链接</el-divider>
        <el-form-item label="GitHub 链接">
          <el-input v-model="form.githubUrl" placeholder="请输入 GitHub 链接" />
        </el-form-item>
        <el-form-item label="微信链接">
          <el-input v-model="form.wechatUrl" placeholder="请输入微信链接或二维码" />
        </el-form-item>

        <el-divider content-position="left">功能设置</el-divider>
        <el-form-item label="评论审核">
          <el-switch v-model="form.commentAudit" />
          <span class="form-tip">开启后，用户评论需要管理员审核后才能显示</span>
        </el-form-item>
        <el-form-item label="分页大小">
          <el-input-number v-model="form.pageSize" :min="5" :max="50" />
          <span class="form-tip">文章列表每页显示的数量</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
            保存设置
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

// 表单引用
const formRef = ref<FormInstance>()
const submitLoading = ref(false)

// 表单数据
const form = reactive({
  siteName: 'Cyblog',
  siteDescription: '',
  seoKeywords: '',
  seoDescription: '',
  githubUrl: '',
  wechatUrl: '',
  pageSize: 10,
  commentAudit: false,
})

// 表单验证规则
const rules: FormRules = {
  siteName: [{ required: true, message: '请输入网站名称', trigger: 'blur' }],
}

// 获取设置
async function fetchSettings(): Promise<void> {
  try {
    // TODO: 调用获取设置 API
    // const data = await getSettings()
    // Object.assign(form, data)

    // 模拟数据
    form.siteName = 'Cyblog'
    form.siteDescription = '一个简洁优雅的博客系统'
    form.seoKeywords = '博客,技术分享,个人网站'
    form.seoDescription = 'Cyblog 是一个使用 Go 和 Vue3 构建的现代化博客系统'
    form.pageSize = 10
    form.commentAudit = false
  } catch {
    ElMessage.error('获取设置失败')
  }
}

// 提交
async function handleSubmit(): Promise<void> {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    // TODO: 调用更新设置 API
    // await updateSettings(form)
    ElMessage.success('保存成功')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    submitLoading.value = false
  }
}

// 重置
function handleReset(): void {
  formRef.value?.resetFields()
  fetchSettings()
}

// 生命周期
onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
.setting-view {
  padding: 0;
  max-width: 800px;
}

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
