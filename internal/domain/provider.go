package domain

import (
	"cyblog/internal/domain/auth"
	"cyblog/internal/domain/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewServiceHub,
	auth.ProviderSet,
	user.ProviderSet,
)
