package route

import (
	"cyblog/internal/domain/category"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(api gin.IRouter, categoryService *category.CategoryService) {
	categoryGroup := api.Group("/categories")
	{
		// 公开接口
		categoryGroup.GET("", categoryService.List)

		// 需要管理员权限的接口
		adminGroup := categoryGroup.Group("", middlewire.AuthMiddleWire(true), middlewire.AdminMiddleWire())
		{
			adminGroup.POST("", categoryService.Create)
			adminGroup.PUT("", categoryService.Update)
			adminGroup.DELETE("/:id", categoryService.Delete)
		}
	}
}
