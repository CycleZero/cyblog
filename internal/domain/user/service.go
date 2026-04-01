package user

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserService struct {
	common.BaseService
	biz    *UserBiz
	logger *log.Logger
}

func NewUserService(biz *UserBiz, logger *log.Logger) *UserService {
	return &UserService{
		biz:    biz,
		logger: logger,
	}
}

func (s *UserService) GetUser(c *gin.Context) {
	var req GetUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Sugar().Error(err)
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}
	user, err := s.biz.GetUser(c, req.Id)
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取用户失败", err))
		return
	}
	s.logger.Info("获取用户成功")
	common.Success(c, GetUserResponse{
		Id:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Avatar: user.Avatar,
		Status: user.Status,
	})
	return

}

// 获取当前登录用户信息
func (s *UserService) GetCurrentUser(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		common.Error(c, errs.New(http.StatusUnauthorized, "用户未登录"))
		return
	}

	user, err := s.biz.GetUser(c, userId.(uint))
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取用户信息失败", err))
		return
	}

	common.Success(c, UserInfoResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format(time.DateTime),
	})
	return
}

// 更新用户信息
func (s *UserService) UpdateUser(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		common.Error(c, errs.New(http.StatusUnauthorized, "用户未登录"))
		return
	}

	var req UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Sugar().Error(err)
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	user, err := s.biz.GetUser(c, userId.(uint))
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取用户信息失败", err))
		return
	}

	updatedUser, err := s.biz.UpdateUser(c, user, req.Name, req.Email, req.Avatar, req.Password)
	if err != nil {
		s.logger.Error(err.Error())
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新用户信息失败", err))
		return
	}

	s.logger.Info("用户信息更新成功", zap.String("name", updatedUser.Name))
	common.Success(c, UserInfoResponse{
		Id:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		Avatar:    updatedUser.Avatar,
		Status:    updatedUser.Status,
		CreatedAt: updatedUser.CreatedAt.Format(time.DateTime),
	})
	return
}
