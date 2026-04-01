package category

import (
	"cyblog/internal/common/dto"
)

// CreateRequest 创建分类请求
// swagger:model CreateCategoryRequest
type CreateRequest struct {
	// 分类名称
	// required: true
	// max length: 50
	Name string `json:"name" binding:"required,max=50"`

	// 分类别名
	// required: true
	// max length: 50
	Slug string `json:"slug" binding:"required,max=50"`

	// 分类描述
	// max length: 200
	Description string `json:"description" binding:"max=200"`

	// 父分类ID
	ParentID uint `json:"parent_id"`

	// 排序
	Sort int `json:"sort"`
}

// UpdateRequest 更新分类请求
// swagger:model UpdateCategoryRequest
type UpdateRequest struct {
	// 分类ID
	// required: true
	ID uint `json:"id" binding:"required"`

	// 分类名称
	// max length: 50
	Name string `json:"name" binding:"max=50"`

	// 分类别名
	// max length: 50
	Slug string `json:"slug" binding:"max=50"`

	// 分类描述
	// max length: 200
	Description string `json:"description" binding:"max=200"`

	// 父分类ID
	ParentID *uint `json:"parent_id"`

	// 排序
	Sort *int `json:"sort"`
}

// Response 分类响应
// swagger:model CategoryResponse
type Response struct {
	// 分类ID
	ID uint `json:"id"`

	// 分类名称
	Name string `json:"name"`

	// 分类别名
	Slug string `json:"slug"`

	// 分类描述
	Description string `json:"description"`

	// 父分类ID
	ParentID uint `json:"parent_id"`

	// 排序
	Sort int `json:"sort"`

	// 创建时间
	CreatedAt string `json:"created_at"`

	// 更新时间
	UpdatedAt string `json:"updated_at"`
}

// ListResponse 分类列表响应
// swagger:model CategoryListResponse
type ListResponse struct {
	dto.PageInfo

	// 分类列表
	List []*Response `json:"list"`
}
