package middlewire

import (
	"cyblog/internal/common"

	"github.com/gin-gonic/gin"
)

var (
	IsMiddleWireRegisterFinished = false
	AuthMiddleWire               func(optional bool) gin.HandlerFunc
)

func AddMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		meta := common.RequestMetadata{
			UserID:    0,
			Request:   c.Request,
			ClientIp:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			RequestId: GenerateRequestId(),
			SessionId: "",
		}
		common.SetRequestMetadata(c, &meta)
	}
}
