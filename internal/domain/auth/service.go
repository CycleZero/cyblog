package auth

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"cyblog/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthService struct {
	common.BaseService
	biz    *AuthBiz
	logger *log.Logger
}

func NewAuthService(biz *AuthBiz, logger *log.Logger) *AuthService {
	return &AuthService{
		biz:    biz,
		logger: logger,
	}
}

// 用户注册
func (s *AuthService) Register(c *gin.Context) {
	var req RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Sugar().Error(err)
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	user, err := s.biz.Register(c, req.Name, req.Email, req.Password)
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.Wrap(http.StatusInternalServerError, "注册失败", err))
		return
	}

	// 生成Token
	token, err := util.GenerateToken(user)
	if err != nil {
		s.logger.Error("生成Token失败", zap.Error(err))
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "注册失败", err))
		return
	}

	s.logger.Info("用户注册成功", zap.String("name", user.Name))
	common.Success(c, RegisterResponse{
		Id:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.Avatar,
		Token:  token,
	})
	return
}

// 用户登录
func (s *AuthService) Login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Sugar().Error(err)
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	user, err := s.biz.Login(c, req.Account, req.Password)
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.Wrap(http.StatusInternalServerError, "登录失败", err))
		return
	}

	// 生成Token
	token, err := util.GenerateToken(user)
	if err != nil {
		s.logger.Error("生成Token失败", zap.Error(err))
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "登录失败", err))
		return
	}

	s.logger.Info("用户登录成功", zap.String("name", user.Name))
	common.Success(c, LoginResponse{
		Id:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Avatar: user.Avatar,
		Token:  token,
	})
	return
}
