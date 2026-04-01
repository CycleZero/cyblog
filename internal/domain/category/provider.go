package category

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCategoryService,
	NewCategoryBiz,
)
