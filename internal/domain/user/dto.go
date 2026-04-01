package user

type GetUserRequest struct {
	Id uint `json:"id"`
}

type GetUserResponse struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
}


type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
	Password string `json:"password" binding:"omitempty,min=6,max=20"`
}

type UserInfoResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}
