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

// AuthService 认证服务
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

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册接口，创建账号并返回用户信息和Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册请求参数"
// @Success 200 {object} common.Response{data=RegisterResponse} "注册成功，返回用户信息和Token"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/auth/register [post]
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

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，支持用户名或邮箱登录，返回用户信息和Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求参数"
// @Success 200 {object} common.Response{data=LoginResponse} "登录成功，返回用户信息和Token"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/auth/login [post]
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
