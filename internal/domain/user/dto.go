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
