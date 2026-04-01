package domain

import (
	"cyblog/internal/domain/article"
	"cyblog/internal/domain/auth"
	"cyblog/internal/domain/category"
	"cyblog/internal/domain/comment"
	"cyblog/internal/domain/tag"
	"cyblog/internal/domain/user"
)

type ServiceHub struct {
	AuthService     *auth.AuthService
	UserService     *user.UserService
	CategoryService *category.CategoryService
	TagService      *tag.TagService
	ArticleService  *article.ArticleService
	CommentService  *comment.CommentService
}

func NewServiceHub(
	authService *auth.AuthService,
	userService *user.UserService,
	categoryService *category.CategoryService,
	tagService *tag.TagService,
	articleService *article.ArticleService,
	commentService *comment.CommentService,
) *ServiceHub {
	return &ServiceHub{
		AuthService:     authService,
		UserService:     userService,
		CategoryService: categoryService,
		TagService:      tagService,
		ArticleService:  articleService,
		CommentService:  commentService,
	}
}
