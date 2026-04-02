// 通用类型
export interface CommonResponse<T = unknown> {
  code: number
  msg: string
  data?: T
}

export interface PageParams {
  page?: number
  pageSize?: number
}

export interface PageResponse<T> {
  list: T[]
  page: number
  pageSize: number
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
  createdAt: string
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
  parentId?: number
  sort?: number
  createdAt: string
  updatedAt: string
}

export interface CategoryListParams extends PageParams {}

export interface CategoryListResponse extends PageResponse<Category> {}

export interface CreateCategoryRequest {
  name: string
  slug: string
  description?: string
  parentId?: number
  sort?: number
}

export interface UpdateCategoryRequest {
  id: number
  name?: string
  slug?: string
  description?: string
  parentId?: number
  sort?: number
}

// 标签相关类型
export interface Tag {
  id: number
  name: string
  slug: string
  color?: string
  count?: number
  createdAt: string
  updatedAt: string
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
  coverImage?: string
  status: number
  isTop?: boolean
  isOriginal?: boolean
  views: number
  likes: number
  createdAt: string
  updatedAt: string
  author: User
  category?: Category
  tags: Tag[]
}

export interface ArticleListParams extends PageParams {
  keyword?: string
  categoryId?: number
  tagId?: number
  status?: number
  authorId?: number
  isTop?: boolean
  sortBy?: string
  sortOrder?: string
}

export interface ArticleListResponse extends PageResponse<Article> {}

export interface CreateArticleRequest {
  title: string
  content: string
  slug?: string
  summary?: string
  coverImage?: string
  status: number
  categoryId?: number
  tagIds?: number[]
  isTop?: boolean
  isOriginal?: boolean
}

export interface UpdateArticleRequest {
  id: number
  title?: string
  content?: string
  slug?: string
  summary?: string
  coverImage?: string
  status?: number
  categoryId?: number
  tagIds?: number[]
  isTop?: boolean
  isOriginal?: boolean
}

// 评论相关类型
export interface Comment {
  id: number
  articleId: number
  userId: number
  parentId?: number
  content: string
  likes: number
  isLiked?: boolean
  createdAt: string
  updatedAt: string
  user: User
  replyTo?: User
  replies?: Comment[]
}

export interface CommentListParams extends PageParams {
  articleId?: number
  parentId?: number
  sortBy?: string
  sortOrder?: string
}

export interface CommentListResponse extends PageResponse<Comment> {}

export interface AdminCommentListParams extends CommentListParams {
  keyword?: string
  userId?: number
}

export interface CreateCommentRequest {
  articleId: number
  content: string
  parentId?: number
}

export interface UpdateCommentRequest {
  id: number
  content: string
}

export interface ArticleCommentCountResponse {
  articleId: number
  commentCount: number
}
