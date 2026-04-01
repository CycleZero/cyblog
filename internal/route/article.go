package route

import (
	"cyblog/internal/domain/article"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

func RegisterArticleRoutes(api gin.IRouter, articleService *article.ArticleService) {
	articleGroup := api.Group("/articles")
	{
		// 公开接口
		articleGroup.GET("", articleService.List)
		articleGroup.GET("/:id", articleService.GetByID)
		articleGroup.GET("/slug/:slug", articleService.GetBySlug)

		// 需要登录的接口
		authGroup := articleGroup.Group("", middlewire.AuthMiddleWire(true))
		{
			authGroup.POST("", articleService.Create)
			authGroup.PUT("", articleService.Update)
			authGroup.DELETE("/:id", articleService.Delete)
			authGroup.POST("/:id/like", articleService.Like)
			authGroup.DELETE("/:id/like", articleService.Unlike)
		}
	}
}
