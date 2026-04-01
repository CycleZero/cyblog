package domain

import (
	"cyblog/internal/domain/auth"
	"cyblog/internal/domain/user"
)

type ServiceHub struct {
	AuthService *auth.AuthService
	UserService *user.UserService
}

func NewServiceHub(
	authService *auth.AuthService,
	userService *user.UserService,
) *ServiceHub {
	return &ServiceHub{
		AuthService: authService,
		UserService: userService,
	}
}
