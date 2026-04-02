package user

// GetUserRequest 获取用户请求
// swagger:model GetUserRequest
type GetUserRequest struct {
	// 用户ID
	Id uint `json:"id"`
}

// GetUserResponse 获取用户响应
// swagger:model GetUserResponse
type GetUserResponse struct {
	// 用户ID
	Id uint `json:"id"`

	// 用户名
	Name string `json:"name"`

	// 邮箱
	Email string `json:"email"`

	// 角色
	Role string `json:"role"`

	// 头像
	Avatar string `json:"avatar"`

	// 状态
	Status int `json:"status"`
}

// UpdateUserRequest 更新用户请求
// swagger:model UpdateUserRequest
type UpdateUserRequest struct {
	// 用户名
	// min length: 2
	// max length: 20
	Name string `json:"name" binding:"omitempty,min=2,max=20"`

	// 邮箱
	// format: email
	Email string `json:"email" binding:"omitempty,email"`

	// 头像
	// format: url
	Avatar string `json:"avatar" binding:"omitempty,url"`

	// 密码
	// min length: 6
	// max length: 20
	Password string `json:"password" binding:"omitempty,min=6,max=20"`
}

// UserInfoResponse 用户信息响应
// swagger:model UserInfoResponse
type UserInfoResponse struct {
	// 用户ID
	Id uint `json:"id"`

	// 用户名
	Name string `json:"name"`

	// 邮箱
	Email string `json:"email"`

	// 角色
	Role string `json:"role"`

	// 头像
	Avatar string `json:"avatar"`

	// 状态
	Status int `json:"status"`

	// 创建时间
	CreatedAt string `json:"created_at"`
}

// AdminListRequest 管理端用户列表请求
// swagger:model AdminUserListRequest
type AdminListRequest struct {
	// 页码
	Page int `form:"page,default=1"`

	// 每页数量
	PageSize int `form:"page_size,default=10"`

	// 关键词
	Keyword string `form:"keyword"`

	// 角色
	Role string `form:"role"`

	// 状态
	Status int `form:"status"`

	// 排序字段
	SortBy string `form:"sort_by"`

	// 排序方式
	SortOrder string `form:"sort_order"`
}

// AdminListResponse 管理端用户列表响应
// swagger:model AdminUserListResponse
type AdminListResponse struct {
	// 总数
	Total int64 `json:"total"`

	// 当前页
	Page int `json:"page"`

	// 用户列表
	List []AdminUserDTO `json:"list"`
}

// AdminUserDTO 管理端用户DTO
// swagger:model AdminUserDTO
type AdminUserDTO struct {
	// 用户ID
	Id uint `json:"id"`

	// 用户名
	Name string `json:"name"`

	// 邮箱
	Email string `json:"email"`

	// 角色
	Role string `json:"role"`

	// 头像
	Avatar string `json:"avatar"`

	// 状态
	Status int `json:"status"`

	// 创建时间
	CreatedAt string `json:"created_at"`
}

// UpdateRoleRequest 更新角色请求
// swagger:model UpdateRoleRequest
type UpdateRoleRequest struct {
	// 角色
	Role string `json:"role" binding:"required"`
}

// UpdateStatusRequest 更新状态请求
// swagger:model UpdateStatusRequest
type UpdateStatusRequest struct {
	// 状态
	Status int `json:"status" binding:"required"`
}
