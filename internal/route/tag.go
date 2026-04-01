package route

import (
	"cyblog/internal/domain/tag"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(api gin.IRouter, tagService *tag.TagService) {
	tagGroup := api.Group("/tags")
	{
		// 公开接口
		tagGroup.GET("", tagService.List)

		// 需要管理员权限的接口
		adminGroup := tagGroup.Group("", middlewire.AuthMiddleWire(true), middlewire.AdminMiddleWire())
		{
			adminGroup.POST("", tagService.Create)
			adminGroup.PUT("", tagService.Update)
			adminGroup.DELETE("/:id", tagService.Delete)
		}
	}
}
