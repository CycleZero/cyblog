package domain

import (
	"cyblog/internal/domain/article"
	"cyblog/internal/domain/auth"
	"cyblog/internal/domain/category"
	"cyblog/internal/domain/comment"
	"cyblog/internal/domain/tag"
	"cyblog/internal/domain/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewServiceHub,
	auth.ProviderSet,
	user.ProviderSet,
	category.ProviderSet,
	tag.ProviderSet,
	article.ProviderSet,
	comment.ProviderSet,
)
