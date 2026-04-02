<template>
  <div class="markdown-editor-wrapper">
    <MdEditor
      v-model="content"
      :height="editorHeight"
      :placeholder="placeholder"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

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

// 将 height prop 转换为纯数字字符串给 md-editor-v3
const editorHeight = computed(() => {
  const num = parseInt(props.height)
  return isNaN(num) ? 500 : num
})

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
</style>
