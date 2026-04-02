// 通用类型
export interface CommonResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface PageParams {
  page?: number
  page_size?: number
}

export interface PageResponse<T> {
  list: T[]
  page: number
  page_size: number
  total: number
}

// 用户相关类型
export interface User {
  id: number
  name: string
  avatar: string
}

export interface UserInfo extends User {
  email: string
  role: string
  status: number
  created_at: string
}

export interface GetUserRequest {
  id: number
}

export interface GetUserResponse {
  id: number
  name: string
  email: string
  avatar: string
  role: string
  status: number
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  avatar?: string
  password?: string
}

// 认证相关类型
export interface LoginRequest {
  account: string
  password: string
}

export interface LoginResponse {
  id: number
  name: string
  email: string
  avatar: string
  role: string
  token: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface RegisterResponse {
  id: number
  name: string
  email: string
  avatar: string
  token: string
}

// 分类相关类型
export interface Category {
  id: number
  name: string
  slug: string
  description?: string
  parent_id?: number
  sort?: number
  created_at: string
  updated_at: string
}

export interface CategoryListParams extends PageParams {}

export interface CategoryListResponse extends PageResponse<Category> {}

export interface CreateCategoryRequest {
  name: string
  slug: string
  description?: string
  parent_id?: number
  sort?: number
}

export interface UpdateCategoryRequest {
  id: number
  name?: string
  slug?: string
  description?: string
  parent_id?: number
  sort?: number
}

// 标签相关类型
export interface Tag {
  id: number
  name: string
  slug: string
  color?: string
  count?: number
  created_at: string
  updated_at: string
}

export interface TagListParams extends PageParams {}

export interface TagListResponse extends PageResponse<Tag> {}

export interface CreateTagRequest {
  name: string
  slug: string
  color?: string
}

export interface UpdateTagRequest {
  id: number
  name?: string
  slug?: string
  color?: string
}

// 文章相关类型
export interface Article {
  id: number
  title: string
  slug?: string
  summary?: string
  content?: string
  cover_image?: string
  status: number
  is_top?: boolean
  is_original?: boolean
  views: number
  likes: number
  created_at: string
  updated_at: string
  author: User
  category?: Category
  tags: Tag[]
}

export interface ArticleListParams extends PageParams {
  keyword?: string
  category_id?: number
  tag_id?: number
  status?: number
  author_id?: number
  is_top?: boolean
  sort_by?: string
  sort_order?: string
}

export interface ArticleListResponse extends PageResponse<Article> {}

export interface CreateArticleRequest {
  title: string
  content: string
  slug?: string
  summary?: string
  cover_image?: string
  status: number
  category_id?: number
  tag_ids?: number[]
  is_top?: boolean
  is_original?: boolean
}

export interface UpdateArticleRequest {
  id: number
  title?: string
  content?: string
  slug?: string
  summary?: string
  cover_image?: string
  status?: number
  category_id?: number
  tag_ids?: number[]
  is_top?: boolean
  is_original?: boolean
}

// 评论相关类型
export interface Comment {
  id: number
  article_id: number
  user_id: number
  parent_id?: number
  content: string
  likes: number
  is_liked?: boolean
  created_at: string
  updated_at: string
  user: User
  reply_to?: User
  replies?: Comment[]
}

export interface CommentListParams extends PageParams {
  article_id?: number
  parent_id?: number
  sort_by?: string
  sort_order?: string
}

export interface CommentListResponse extends PageResponse<Comment> {}

export interface AdminCommentListParams extends CommentListParams {
  keyword?: string
  user_id?: number
}

export interface CreateCommentRequest {
  article_id: number
  content: string
  parent_id?: number
}

export interface UpdateCommentRequest {
  id: number
  content: string
}

export interface ArticleCommentCountResponse {
  article_id: number
  comment_count: number
}
