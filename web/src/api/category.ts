import request from './request'
import type {
  Category,
  CategoryListParams,
  CategoryListResponse,
  CreateCategoryRequest,
  UpdateCategoryRequest,
} from './types'

// 获取分类列表
export function getCategories(
  params?: CategoryListParams
): Promise<CategoryListResponse> {
  return request.get('/categories', { params })
}

// 创建分类
export function createCategory(
  data: CreateCategoryRequest
): Promise<Category> {
  return request.post('/categories', data)
}

// 更新分类
export function updateCategory(
  data: UpdateCategoryRequest
): Promise<Category> {
  return request.put('/categories', data)
}

// 删除分类
export function deleteCategory(id: number): Promise<void> {
  return request.delete(`/categories/${id}`)
}
