package article

import (
	"cyblog/internal/common"
	"cyblog/internal/domain/category"
	"cyblog/internal/domain/tag"
	"time"
)

// CreateRequest 创建文章请求
type CreateRequest struct {
	Title      string `json:"title" binding:"required,max=255"`
	Slug       string `json:"slug" binding:"omitempty,max=100"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"omitempty,url"`
	Status     int    `json:"status" binding:"oneof=1 2 3"`
	IsTop      bool   `json:"is_top"`
	IsOriginal bool   `json:"is_original"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// UpdateRequest 更新文章请求
type UpdateRequest struct {
	ID         uint   `json:"id" binding:"required"`
	Title      string `json:"title" binding:"max=255"`
	Slug       string `json:"slug" binding:"omitempty,max=100"`
	Content    string `json:"content"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"omitempty,url"`
	Status     int    `json:"status" binding:"omitempty,oneof=1 2 3"`
	IsTop      *bool  `json:"is_top"`
	IsOriginal *bool  `json:"is_original"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// ListQuery 文章列表查询参数
type ListQuery struct {
	common.PageParam
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	Status     int    `form:"status"`
	AuthorID   uint   `form:"author_id"`
	IsTop      *bool  `form:"is_top"`
	SortBy     string `form:"sort_by,default=created_at"`
	SortOrder  string `form:"sort_order,default=desc"`
}

// Getter方法供Repo层动态查询使用
func (q *ListQuery) GetKeyword() string   { return q.Keyword }
func (q *ListQuery) GetCategoryID() uint  { return q.CategoryID }
func (q *ListQuery) GetTagID() uint       { return q.TagID }
func (q *ListQuery) GetStatus() int       { return q.Status }
func (q *ListQuery) GetAuthorID() uint    { return q.AuthorID }
func (q *ListQuery) GetIsTop() *bool      { return q.IsTop }
func (q *ListQuery) GetSortBy() string    { return q.SortBy }
func (q *ListQuery) GetSortOrder() string { return q.SortOrder }
func (q *ListQuery) GetPage() int         { return q.Page }
func (q *ListQuery) GetPageSize() int     { return q.PageSize }

// Response 文章响应
type Response struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"cover_image"`
	Status     int       `json:"status"`
	Views      int       `json:"views"`
	Likes      int       `json:"likes"`
	IsTop      bool      `json:"is_top"`
	IsOriginal bool      `json:"is_original"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Author     *struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"author"`
	Category *category.Response `json:"category"`
	Tags     []*tag.Response    `json:"tags"`
}

// ListResponse 文章列表响应
type ListResponse struct {
	common.PageInfo
	List []*Response `json:"list"`
}
