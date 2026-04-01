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

func (s *BaseService) Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
		"msg":  constant.MsgSuccess,
	})
}

func (s *BaseService) Fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code": code,
		"data": nil,
		"msg":  msg,
	})
}

func (s *BaseService) Error(c *gin.Context, err *errs.CyBlogError) {
	s.Fail(c, errs.Code(err), err.Error())
}
