package route

import (
	"cyblog/internal/domain"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(root gin.IRouter, serviceHub *domain.ServiceHub) {
	adminGroup := root.Group("/admin")
	adminGroup.Use(middlewire.AdminMiddleWire())
	{
		// 仪表盘
		adminGroup.GET("/dashboard", serviceHub.ArticleService.GetDashboard)

		// 文章管理
		articleGroup := adminGroup.Group("/articles")
		{
			articleGroup.GET("", serviceHub.ArticleService.AdminList)
			articleGroup.POST("/:id/top", serviceHub.ArticleService.SetTop)
			articleGroup.POST("/batch-delete", serviceHub.ArticleService.BatchDelete)
			articleGroup.POST("/batch-status", serviceHub.ArticleService.BatchUpdateStatus)
		}

		// 评论管理
		commentGroup := adminGroup.Group("/comments")
		{
			commentGroup.GET("", serviceHub.CommentService.AdminList)
			commentGroup.DELETE("/:id", serviceHub.CommentService.AdminDelete)
		}

		// 用户管理
		userGroup := adminGroup.Group("/users")
		{
			userGroup.GET("", serviceHub.UserService.AdminList)
			userGroup.PUT("/:id/role", serviceHub.UserService.UpdateRole)
			userGroup.PUT("/:id/status", serviceHub.UserService.UpdateStatus)
		}
	}
}
