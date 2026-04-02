import request from './request'
import type {
  Article,
  ArticleListParams,
  ArticleListResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
} from './types'

// 获取文章列表
export function getArticles(params?: ArticleListParams): Promise<ArticleListResponse> {
  return request.get('/articles', { params })
}

// 获取文章详情
export function getArticle(
  id: number,
  incrementViews = true
): Promise<Article> {
  return request.get(`/articles/${id}`, {
    params: { increment_views: incrementViews },
  })
}

// 根据别名获取文章
export function getArticleBySlug(
  slug: string,
  incrementViews = true
): Promise<Article> {
  return request.get(`/articles/slug/${slug}`, {
    params: { increment_views: incrementViews },
  })
}

// 创建文章
export function createArticle(data: CreateArticleRequest): Promise<Article> {
  return request.post('/articles', data)
}

// 更新文章
export function updateArticle(data: UpdateArticleRequest): Promise<Article> {
  return request.put('/articles', data)
}

// 删除文章
export function deleteArticle(id: number): Promise<void> {
  return request.delete(`/articles/${id}`)
}

// 点赞文章
export function likeArticle(id: number): Promise<void> {
  return request.post(`/articles/${id}/like`)
}

// 取消点赞
export function unlikeArticle(id: number): Promise<void> {
  return request.delete(`/articles/${id}/like`)
}
