package user

import (
	"context"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"cyblog/pkg/repo"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserBiz struct {
	userRepo *repo.UserRepo
	logger   *log.Logger
}

func NewUserBiz(
	userRepo *repo.UserRepo,
	logger *log.Logger,
) *UserBiz {
	return &UserBiz{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (biz *UserBiz) GetUser(ctx context.Context, id uint) (*model.User, error) {
	u, err := biz.userRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "用户不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取用户失败", err)
	}
	return u, nil
}

// 更新用户信息
func (biz *UserBiz) UpdateUser(ctx context.Context, user *model.User, name, email, avatar, password string) (*model.User, error) {
	if name != "" && name != user.Name {
		// 检查新用户名是否已存在
		existUser, err := biz.userRepo.GetUserByName(ctx, name)
		if err == nil && existUser.ID != user.ID {
			return nil, errs.New(http.StatusBadRequest, "用户名已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查用户名失败", err)
		}
		user.Name = name
	}

	if email != "" && email != user.Email {
		// 检查新邮箱是否已存在
		existUser, err := biz.userRepo.GetUserByEmail(ctx, email)
		if err == nil && existUser.ID != user.ID {
			return nil, errs.New(http.StatusBadRequest, "邮箱已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查邮箱失败", err)
		}
		user.Email = email
	}

	if avatar != "" {
		user.Avatar = avatar
	}

	if password != "" {
		// 加密新密码
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errs.Wrap(http.StatusInternalServerError, "密码加密失败", err)
		}
		user.Password = string(bytes)
	}

	err := biz.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "更新用户信息失败", err)
	}

	return user, nil
}

// AdminList 管理端获取用户列表
func (biz *UserBiz) AdminList(ctx context.Context, req *AdminListRequest) (*AdminListResponse, error) {
	users, total, err := biz.userRepo.ListUsers(ctx, req)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取用户列表失败", err)
	}

	list := make([]AdminUserDTO, 0, len(users))
	for _, u := range users {
		list = append(list, AdminUserDTO{
			Id:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &AdminListResponse{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, nil
}

// UpdateRole 更新用户角色
func (biz *UserBiz) UpdateRole(ctx context.Context, id uint, role string) error {
	user, err := biz.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errs.New(http.StatusNotFound, "用户不存在")
	}
	user.Role = role
	return biz.userRepo.UpdateUser(ctx, user)
}

// UpdateStatus 更新用户状态
func (biz *UserBiz) UpdateStatus(ctx context.Context, id uint, status int) error {
	user, err := biz.userRepo.GetUserById(ctx, id)
	if err != nil {
		return errs.New(http.StatusNotFound, "用户不存在")
	}
	user.Status = status
	return biz.userRepo.UpdateUser(ctx, user)
}
