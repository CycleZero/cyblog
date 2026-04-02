package user

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"strconv"
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

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户公开信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body GetUserRequest true "获取用户请求参数"
// @Success 200 {object} common.Response{data=GetUserResponse} "获取成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/user/{id} [get]
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

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户的详细信息，需要认证
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.Response{data=UserInfoResponse} "获取成功"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/user/info [get]
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

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户的信息，需要认证
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateUserRequest true "更新用户请求参数"
// @Success 200 {object} common.Response{data=UserInfoResponse} "更新成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "用户未登录"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/user/update [put]
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

// AdminList 管理端获取用户列表
// @Summary 管理端获取用户列表
// @Description 获取所有用户列表，支持分页，需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=AdminListResponse} "获取成功"
// @Failure 401 {object} common.Response "未授权"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/admin/users [get]
func (s *UserService) AdminList(c *gin.Context) {
	var req AdminListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}
	resp, bizErr := s.biz.AdminList(c, &req)
	if bizErr != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取用户列表失败", bizErr))
		return
	}
	common.Success(c, resp)
}

// UpdateRole 更新用户角色
// @Summary 更新用户角色
// @Description 更新指定用户的角色，需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body UpdateRoleRequest true "更新角色请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "无效的用户ID或参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/admin/users/{id}/role [put]
func (s *UserService) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的用户ID"))
		return
	}
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}
	err = s.biz.UpdateRole(c, uint(id), req.Role)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新角色失败", err))
		return
	}
	common.Success(c, nil)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新指定用户的状态（启用/禁用），需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body UpdateStatusRequest true "更新状态请求参数"
// @Success 200 {object} common.Response "更新成功"
// @Failure 400 {object} common.Response "无效的用户ID或参数错误"
// @Failure 401 {object} common.Response "未授权"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/admin/users/{id}/status [put]
func (s *UserService) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的用户ID"))
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}
	err = s.biz.UpdateStatus(c, uint(id), req.Status)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新状态失败", err))
		return
	}
	common.Success(c, nil)
}
