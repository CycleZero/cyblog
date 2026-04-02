import request from './request'
import type {
  Comment,
  CommentListParams,
  CommentListResponse,
  AdminCommentListParams,
  CreateCommentRequest,
  UpdateCommentRequest,
  ArticleCommentCountResponse,
} from './types'

// 获取文章评论列表
export function getComments(
  params: CommentListParams
): Promise<CommentListResponse> {
  return request.get('/comments', { params })
}

// 获取评论的回复列表
export function getCommentReplies(
  id: number,
  params?: { page?: number; page_size?: number }
): Promise<CommentListResponse> {
  return request.get(`/comments/${id}/replies`, { params })
}

// 创建评论
export function createComment(data: CreateCommentRequest): Promise<Comment> {
  return request.post('/comments', data)
}

// 更新评论
export function updateComment(data: UpdateCommentRequest): Promise<Comment> {
  return request.put(`/comments/${data.id}`, { content: data.content })
}

// 删除评论
export function deleteComment(id: number): Promise<void> {
  return request.delete(`/comments/${id}`)
}

// 点赞评论
export function likeComment(id: number): Promise<void> {
  return request.post(`/comments/${id}/like`)
}

// 取消点赞评论
export function unlikeComment(id: number): Promise<void> {
  return request.delete(`/comments/${id}/like`)
}

// 获取文章评论数
export function getArticleCommentCount(
  articleId: number
): Promise<ArticleCommentCountResponse> {
  return request.get(`/articles/${articleId}/comment-count`)
}

// 管理端获取评论列表
export function getAdminComments(
  params?: AdminCommentListParams
): Promise<CommentListResponse> {
  return request.get('/admin/comments', { params })
}

// 管理端删除评论
export function adminDeleteComment(id: number): Promise<void> {
  return request.delete(`/admin/comments/${id}`)
}
