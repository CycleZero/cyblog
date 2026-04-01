package common

import (
	"cyblog/constant"
	"cyblog/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseService struct {
}

// Response 通用响应结构
// swagger:model CommonResponse
type Response struct {
	// 响应码
	Code int `json:"code"`

	// 响应数据
	Data any `json:"data"`

	// 响应消息
	Msg string `json:"msg"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
		"msg":  constant.MsgSuccess,
	})
}

func Fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code": code,
		"data": nil,
		"msg":  msg,
	})
}

func Error(c *gin.Context, err *errs.CyBlogError) {
	Fail(c, errs.Code(err), err.Error())
}
