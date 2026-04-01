package article

import (
	"cyblog/internal/common"
	"cyblog/internal/domain/category"
	"cyblog/internal/domain/tag"
	"time"
)

// CreateRequest 创建文章请求
// swagger:model CreateArticleRequest
type CreateRequest struct {
	// 标题
	// required: true
	// max length: 255
	Title string `json:"title" binding:"required,max=255"`

	// 别名
	// max length: 100
	Slug string `json:"slug" binding:"omitempty,max=100"`

	// 内容
	// required: true
	Content string `json:"content" binding:"required"`

	// 摘要
	// max length: 500
	Summary string `json:"summary" binding:"max=500"`

	// 封面图片
	// format: url
	CoverImage string `json:"cover_image" binding:"omitempty,url"`

	// 状态 1:草稿 2:已发布 3:私密
	// required: true
	// enum: 1,2,3
	Status int `json:"status" binding:"oneof=1 2 3"`

	// 是否置顶
	IsTop bool `json:"is_top"`

	// 是否原创
	IsOriginal bool `json:"is_original"`

	// 分类ID
	CategoryID uint `json:"category_id"`

	// 标签ID列表
	TagIDs []uint `json:"tag_ids"`
}

// UpdateRequest 更新文章请求
// swagger:model UpdateArticleRequest
type UpdateRequest struct {
	// 文章ID
	// required: true
	ID uint `json:"id" binding:"required"`

	// 标题
	// max length: 255
	Title string `json:"title" binding:"max=255"`

	// 别名
	// max length: 100
	Slug string `json:"slug" binding:"omitempty,max=100"`

	// 内容
	Content string `json:"content"`

	// 摘要
	// max length: 500
	Summary string `json:"summary" binding:"max=500"`

	// 封面图片
	// format: url
	CoverImage string `json:"cover_image" binding:"omitempty,url"`

	// 状态 1:草稿 2:已发布 3:私密
	// enum: 1,2,3
	Status int `json:"status" binding:"omitempty,oneof=1 2 3"`

	// 是否置顶
	IsTop *bool `json:"is_top"`

	// 是否原创
	IsOriginal *bool `json:"is_original"`

	// 分类ID
	CategoryID *uint `json:"category_id"`

	// 标签ID列表
	TagIDs []uint `json:"tag_ids"`
}

// ListQuery 文章列表查询参数
// swagger:model ArticleListQuery
type ListQuery struct {
	common.PageParam

	// 关键词
	Keyword string `form:"keyword"`

	// 分类ID
	CategoryID uint `form:"category_id"`

	// 标签ID
	TagID uint `form:"tag_id"`

	// 状态
	Status int `form:"status"`

	// 作者ID
	AuthorID uint `form:"author_id"`

	// 是否置顶
	IsTop *bool `form:"is_top"`

	// 排序字段
	SortBy string `form:"sort_by,default=created_at"`

	// 排序方式
	SortOrder string `form:"sort_order,default=desc"`
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
// swagger:model ArticleResponse
type Response struct {
	// 文章ID
	ID uint `json:"id"`

	// 标题
	Title string `json:"title"`

	// 别名
	Slug string `json:"slug"`

	// 摘要
	Summary string `json:"summary"`

	// 封面图片
	CoverImage string `json:"cover_image"`

	// 状态
	Status int `json:"status"`

	// 浏览量
	Views int `json:"views"`

	// 点赞数
	Likes int `json:"likes"`

	// 是否置顶
	IsTop bool `json:"is_top"`

	// 是否原创
	IsOriginal bool `json:"is_original"`

	// 创建时间
	CreatedAt time.Time `json:"created_at"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// 作者
	Author *struct {
		// 作者ID
		ID uint `json:"id"`

		// 作者名
		Name string `json:"name"`

		// 头像
		Avatar string `json:"avatar"`
	} `json:"author"`

	// 分类
	Category *category.Response `json:"category"`

	// 标签列表
	Tags []*tag.Response `json:"tags"`
}

// ListResponse 文章列表响应
// swagger:model ArticleListResponse
type ListResponse struct {
	common.PageInfo

	// 文章列表
	List []*Response `json:"list"`
}
