package common

import (
	"cyblog/constant"
	"cyblog/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseService struct {
}

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
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
