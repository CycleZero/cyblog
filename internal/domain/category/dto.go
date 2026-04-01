package category

import "cyblog/internal/common"

// CreateRequest 创建分类请求
type CreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Slug        string `json:"slug" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
}

// UpdateRequest 更新分类请求
type UpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"max=50"`
	Slug        string `json:"slug" binding:"max=50"`
	Description string `json:"description" binding:"max=200"`
	ParentID    *uint  `json:"parent_id"`
	Sort        *int   `json:"sort"`
}

// Response 分类响应
type Response struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	Sort        int    `json:"sort"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListResponse 分类列表响应
type ListResponse struct {
	common.PageInfo
	List []*Response `json:"list"`
}
