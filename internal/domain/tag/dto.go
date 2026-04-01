package tag

import "cyblog/internal/common"

// CreateRequest 创建标签请求
// swagger:model CreateTagRequest
type CreateRequest struct {
	// 标签名称
	// required: true
	// max length: 50
	Name string `json:"name" binding:"required,max=50"`

	// 标签别名
	// required: true
	// max length: 50
	Slug string `json:"slug" binding:"required,max=50"`

	// 标签颜色（十六进制）
	// format: hexcolor
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

// UpdateRequest 更新标签请求
// swagger:model UpdateTagRequest
type UpdateRequest struct {
	// 标签ID
	// required: true
	ID uint `json:"id" binding:"required"`

	// 标签名称
	// max length: 50
	Name string `json:"name" binding:"max=50"`

	// 标签别名
	// max length: 50
	Slug string `json:"slug" binding:"max=50"`

	// 标签颜色（十六进制）
	// format: hexcolor
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

// Response 标签响应
// swagger:model TagResponse
type Response struct {
	// 标签ID
	ID uint `json:"id"`

	// 标签名称
	Name string `json:"name"`

	// 标签别名
	Slug string `json:"slug"`

	// 标签颜色
	Color string `json:"color"`

	// 使用数量
	Count int `json:"count"`

	// 创建时间
	CreatedAt string `json:"created_at"`

	// 更新时间
	UpdatedAt string `json:"updated_at"`
}

// ListResponse 标签列表响应
// swagger:model TagListResponse
type ListResponse struct {
	common.PageInfo

	// 标签列表
	List []*Response `json:"list"`
}
