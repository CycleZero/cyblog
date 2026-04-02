package route

import (
	"cyblog/internal/domain/comment"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(api gin.IRouter, commentService *comment.CommentService) {
	// 评论相关路由
	commentGroup := api.Group("/comments")
	{
		// 公开接口
		commentGroup.GET("", commentService.List)
		commentGroup.GET("/:id/replies", commentService.GetReplies)

		// 需要登录的接口
		authGroup := commentGroup.Group("", middlewire.AuthMiddleWire(true))
		{
			authGroup.POST("", commentService.Create)
			authGroup.PUT("/:id", commentService.Update)
			authGroup.DELETE("/:id", commentService.Delete)
			authGroup.POST("/:id/like", commentService.Like)
			authGroup.DELETE("/:id/like", commentService.Unlike)
		}
	}

	// 文章评论数路由（挂载在articles下）
	articleGroup := api.Group("/articles")
	{
		articleGroup.GET("/:id/comment-count", commentService.GetCommentCount)
	}

}
