package tag

import "cyblog/internal/common"

// CreateRequest 创建标签请求
type CreateRequest struct {
	Name  string `json:"name" binding:"required,max=50"`
	Slug  string `json:"slug" binding:"required,max=50"`
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

// UpdateRequest 更新标签请求
type UpdateRequest struct {
	ID    uint   `json:"id" binding:"required"`
	Name  string `json:"name" binding:"max=50"`
	Slug  string `json:"slug" binding:"max=50"`
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

// Response 标签响应
type Response struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListResponse 标签列表响应
type ListResponse struct {
	common.PageInfo
	List []*Response `json:"list"`
}
