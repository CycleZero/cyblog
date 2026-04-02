import request from './request'
import type {
  Tag,
  TagListParams,
  TagListResponse,
  CreateTagRequest,
  UpdateTagRequest,
} from './types'

// 获取标签列表
export function getTags(params?: TagListParams): Promise<TagListResponse> {
  return request.get('/tags', { params })
}

// 创建标签
export function createTag(data: CreateTagRequest): Promise<Tag> {
  return request.post('/tags', data)
}

// 更新标签
export function updateTag(data: UpdateTagRequest): Promise<Tag> {
  return request.put('/tags', data)
}

// 删除标签
export function deleteTag(id: number): Promise<void> {
  return request.delete(`/tags/${id}`)
}
