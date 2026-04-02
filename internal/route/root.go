package route

import (
	"cyblog/internal/domain"
	"cyblog/internal/domain/auth"
	"cyblog/internal/domain/user"
	"cyblog/internal/route/middlewire"

	"github.com/gin-gonic/gin"
)

type RegisterFunc func(root gin.IRouter, serviceHub *domain.ServiceHub)

func NewRegisterFunc() RegisterFunc {
	return RegisterRouter
}

func RegisterRouter(root gin.IRouter, serviceHub *domain.ServiceHub) {
	root.Use(middlewire.CORS())
	root.Use(middlewire.AddMetaData())
	root.Use(middlewire.AuthMiddleWire(false))

	// 公开路由组
	apiGroup := root.Group("/api/")

	// 注册路由
	RegisterAuthRoutes(apiGroup, serviceHub.AuthService)
	RegisterUserRoutes(apiGroup, serviceHub.UserService)
	RegisterCategoryRoutes(apiGroup, serviceHub.CategoryService)
	RegisterTagRoutes(apiGroup, serviceHub.TagService)
	RegisterArticleRoutes(apiGroup, serviceHub.ArticleService)
	RegisterCommentRoutes(apiGroup, serviceHub.CommentService)

	// 注册管理端路由
	RegisterAdminRoutes(apiGroup, serviceHub)
}

func RegisterUserRoutes(api gin.IRouter, userService *user.UserService) {
	userGroup := api.Group("/user")
	{
		// 公开用户接口
		userGroup.GET("/:id", userService.GetUser)

		// 需要认证的用户接口（统一应用鉴权中间件）
		authUserGroup := userGroup.Group("", middlewire.AuthMiddleWire(true))
		{
			authUserGroup.GET("/info", userService.GetCurrentUser)
			authUserGroup.PUT("/update", userService.UpdateUser)
		}
	}
}

func RegisterAuthRoutes(publicApi gin.IRouter, authService *auth.AuthService) {
	authGroup := publicApi.Group("/auth")
	{
		authGroup.POST("/register", authService.Register)
		authGroup.POST("/login", authService.Login)
	}
}
