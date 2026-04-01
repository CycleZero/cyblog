package common

import (
	"context"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestMetadata struct {
	UserID uint

	// User 为nil时表示游客（未登录）
	User      model.User
	Request   *http.Request
	ClientIp  string
	RealIp    string
	UserAgent string
	RequestId string
	SessionId string

	IsOk bool
}

func GetRequestMetadata(c context.Context) *RequestMetadata {
	res, ok := c.Value("request_metadata").(*RequestMetadata)
	if !ok || res == nil {
		log.GetLogger().Error("获取请求Metadata失败")
		return &RequestMetadata{
			IsOk: false,
		}
	}
	return res
}

func SetRequestMetadata(c *gin.Context, metadata *RequestMetadata) {
	c.Set("request_metadata", metadata)
}
