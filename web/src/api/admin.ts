/**
 * 后台管理 API
 */
import request from './request'
import type { Article, PageResponse } from './types'

// ============== 类型定义 ==============

/** 仪表盘响应 */
export interface DashboardResponse {
  articleCount: number
  todayViews: number
  recentArticles: RecentArticle[]
  hotArticles: HotArticle[]
}

/** 近期文章（简化） */
export interface RecentArticle {
  id: number
  title: string
  status: number
}

/** 热门文章（简化） */
export interface HotArticle {
  id: number
  title: string
  views: number
}

/** 文章状态枚举 */
export enum ArticleStatus {
  Draft = 1,
  Published = 2,
  Private = 3,
}

/** 文章状态映射 */
export const ArticleStatusMap: Record<number, string> = {
  [ArticleStatus.Draft]: '草稿',
  [ArticleStatus.Published]: '已发布',
  [ArticleStatus.Private]: '私密',
}

/** 文章状态类型映射 */
export const ArticleStatusTypeMap: Record<number, string> = {
  [ArticleStatus.Draft]: 'info',
  [ArticleStatus.Published]: 'success',
  [ArticleStatus.Private]: 'warning',
}

/** 用户角色枚举 */
export enum UserRole {
  Admin = 'admin',
  Editor = 'editor',
  User = 'user',
}

/** 用户角色映射 */
export const UserRoleMap: Record<string, string> = {
  [UserRole.Admin]: '管理员',
  [UserRole.Editor]: '编辑',
  [UserRole.User]: '用户',
}

/** 用户角色类型映射 */
export const UserRoleTypeMap: Record<string, string> = {
  [UserRole.Admin]: 'danger',
  [UserRole.Editor]: 'warning',
  [UserRole.User]: 'info',
}

/** 用户状态枚举 */
export enum UserStatus {
  Disabled = 0,
  Enabled = 1,
}

// ============== 仪表盘 ==============

/**
 * 获取仪表盘数据
 */
export function getDashboard(): Promise<DashboardResponse> {
  return request.get('/admin/dashboard')
}

// ============== 文章管理 ==============

/** 文章列表查询参数 */
export interface AdminArticleParams {
  keyword?: string
  category_id?: number
  status?: number
  page?: number
  pageSize?: number
}

/**
 * 管理端获取文章列表
 */
export function getAdminArticles(params: AdminArticleParams): Promise<PageResponse<Article>> {
  return request.get('/admin/articles', { params })
}

/** 置顶请求 */
export interface SetTopRequest {
  is_top: boolean
}

/**
 * 置顶/取消置顶文章
 */
export function setArticleTop(id: number, isTop: boolean): Promise<void> {
  return request.put(`/admin/articles/${id}/top`, { is_top: isTop })
}

/** 批量删除请求 */
export interface BatchDeleteRequest {
  ids: number[]
}

/**
 * 批量删除文章
 */
export function batchDeleteArticles(ids: number[]): Promise<void> {
  return request.post('/admin/articles/batch-delete', { ids })
}

/** 批量更新状态请求 */
export interface BatchUpdateStatusRequest {
  ids: number[]
  status: number
}

/**
 * 批量更新文章状态
 */
export function batchUpdateArticleStatus(ids: number[], status: number): Promise<void> {
  return request.put('/admin/articles/batch-status', { ids, status })
}

// ============== 评论管理 ==============

/** 评论列表查询参数 */
export interface AdminCommentParams {
  keyword?: string
  article_id?: number
  user_id?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  page?: number
  pageSize?: number
}

/**
 * 管理端获取评论列表
 */
export function getAdminComments(params: AdminCommentParams): Promise<unknown> {
  return request.get('/admin/comments', { params })
}

/**
 * 管理端删除评论
 */
export function deleteAdminComment(id: number): Promise<void> {
  return request.delete(`/admin/comments/${id}`)
}

// ============== 用户管理 ==============

/** 管理员用户信息 */
export interface AdminUser {
  id: number
  name: string
  email: string
  avatar: string
  role: string
  status: number
  created_at: string
}

/** 用户列表响应 */
export interface AdminUserListResponse {
  list: AdminUser[]
  page: number
  total: number
}

/** 用户列表查询参数 */
export interface AdminUserParams {
  keyword?: string
  role?: string
  status?: number
  page?: number
  pageSize?: number
}

/** 更新角色请求 */
export interface UpdateRoleRequest {
  role: string
}

/** 更新状态请求 */
export interface UpdateStatusRequest {
  status: number
}

/**
 * 获取用户列表
 */
export function getAdminUsers(params: AdminUserParams): Promise<AdminUserListResponse> {
  return request.get('/admin/users', { params })
}

/**
 * 更新用户角色
 */
export function updateUserRole(id: number, role: string): Promise<void> {
  return request.put(`/admin/users/${id}/role`, { role })
}

/**
 * 更新用户状态
 */
export function updateUserStatus(id: number, status: number): Promise<void> {
  return request.put(`/admin/users/${id}/status`, { status })
}

// ============== 系统设置 ==============

/** 系统设置 */
export interface SystemSettings {
  siteName: string
  siteDescription: string
  seoKeywords: string
  seoDescription: string
  githubUrl?: string
  wechatUrl?: string
  pageSize: number
  commentAudit: boolean
}

/**
 * 获取系统设置
 */
export function getSettings(): Promise<SystemSettings> {
  return request.get('/admin/settings')
}

/**
 * 更新系统设置
 */
export function updateSettings(data: SystemSettings): Promise<void> {
  return request.put('/admin/settings', data)
}
