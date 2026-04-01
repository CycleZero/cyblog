package tag

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewTagBiz,
	NewTagService,
)
