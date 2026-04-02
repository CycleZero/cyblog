<template>
  <div class="markdown-editor-wrapper">
    <v-md-editor
      v-model="content"
      :height="height"
      mode="edit"
      :placeholder="placeholder"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
// 使用已配置好的 VMdEditor
import VMdEditor from '@kangc/v-md-editor'

interface Props {
  modelValue?: string
  height?: string
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  height: '500px',
  placeholder: '请输入内容...',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: string]
}>()

const content = ref(props.modelValue)

// 监听 props 变化
watch(() => props.modelValue, (newVal) => {
  content.value = newVal
})

function handleChange(value: string) {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.markdown-editor-wrapper {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.markdown-editor-wrapper :deep(.v-md-editor) {
  border: none;
  box-shadow: none;
}
</style>
