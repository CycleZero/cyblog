package middlewire

import (
	"cyblog/internal/common"
	"cyblog/pkg/log"
	"cyblog/pkg/model"

	"github.com/gin-gonic/gin"
)

func AdminMiddleWire() gin.HandlerFunc {
	return func(c *gin.Context) {
		meta := common.GetRequestMetadata(c)
		if meta.UserID == 0 {
			log.SugaredLogger().Error("用户未登录")
			c.Abort()
			return
		}
		if meta.User.Role != model.RoleAdmin {
			log.SugaredLogger().Error("用户权限不足")
			common.Fail(c, 403, "用户权限不足")
			c.Abort()
			return
		}
		log.SugaredLogger().Info("管理员登录", meta.User.Name)
		c.Next()
		return
	}
}
