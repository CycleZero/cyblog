package auth

// RegisterRequest 注册请求
// swagger:model RegisterRequest
type RegisterRequest struct {
	// 用户名
	// required: true
	// min length: 2
	// max length: 20
	Name string `json:"name" binding:"required,min=2,max=20"`

	// 邮箱
	// required: true
	// format: email
	Email string `json:"email" binding:"required,email"`

	// 密码
	// required: true
	// min length: 6
	// max length: 20
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// RegisterResponse 注册响应
// swagger:model RegisterResponse
type RegisterResponse struct {
	// 用户ID
	Id uint `json:"id"`

	// 用户名
	Name string `json:"name"`

	// 邮箱
	Email string `json:"email"`

	// 头像
	Avatar string `json:"avatar"`

	// 认证Token
	Token string `json:"token"`
}

// LoginRequest 登录请求
// swagger:model LoginRequest
type LoginRequest struct {
	// 账号（用户名或邮箱）
	// required: true
	Account string `json:"account" binding:"required"`

	// 密码
	// required: true
	// min length: 6
	// max length: 20
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// LoginResponse 登录响应
// swagger:model LoginResponse
type LoginResponse struct {
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

	// 认证Token
	Token string `json:"token"`
}
