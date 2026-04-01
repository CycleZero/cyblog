package internal

import (
	"cyblog/internal/domain"
	"cyblog/internal/route"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	domain.ProviderSet,
	route.ProviderSet,
)
