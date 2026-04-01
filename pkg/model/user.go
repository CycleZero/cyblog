package model

import "gorm.io/gorm"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

const (
	StatusActive = 1
	StatusBan    = 0
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:50;not null;uniqueIndex"`
	Email    string `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Password string `json:"-" gorm:"size:255;not null"`
	Role     string `json:"role" gorm:"size:20;default:'user'"`
	Status   int    `json:"status" gorm:"default:1"`
	Avatar   string `json:"avatar" gorm:"size:255"`
}
