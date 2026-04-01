package article

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewArticleBiz,
	NewArticleService,
)
