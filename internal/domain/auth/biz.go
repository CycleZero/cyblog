package auth

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

type AuthBiz struct {
	userRepo *repo.UserRepo
	logger   *log.Logger
}

func NewAuthBiz(
	userRepo *repo.UserRepo,
	logger *log.Logger,
) *AuthBiz {
	return &AuthBiz{
		userRepo: userRepo,
		logger:   logger,
	}
}

// 密码加密
func (biz *AuthBiz) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 密码验证
func (biz *AuthBiz) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 用户注册
func (biz *AuthBiz) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	_, err := biz.userRepo.GetUserByName(ctx, name)
	if err == nil {
		return nil, errs.New(http.StatusBadRequest, "用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(http.StatusInternalServerError, "检查用户名失败", err)
	}

	// 检查邮箱是否已存在
	_, err = biz.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errs.New(http.StatusBadRequest, "邮箱已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(http.StatusInternalServerError, "检查邮箱失败", err)
	}

	// 加密密码
	hashPwd, err := biz.hashPassword(password)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "密码加密失败", err)
	}

	// 创建用户
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashPwd,
		Role:     model.RoleUser,
		Status:   model.StatusActive,
	}

	err = biz.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "创建用户失败", err)
	}

	return user, nil
}

// 用户登录
func (biz *AuthBiz) Login(ctx context.Context, account, password string) (*model.User, error) {
	// 先尝试用用户名查找
	user, err := biz.userRepo.GetUserByName(ctx, account)
	if err != nil {
		// 用户名没找到，尝试用邮箱查找
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = biz.userRepo.GetUserByEmail(ctx, account)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errs.New(http.StatusBadRequest, "用户名或密码错误")
				}
				return nil, errs.Wrap(http.StatusInternalServerError, "查询用户失败", err)
			}
		} else {
			return nil, errs.Wrap(http.StatusInternalServerError, "查询用户失败", err)
		}
	}

	// 检查用户状态
	if user.Status != model.StatusActive {
		return nil, errs.New(http.StatusForbidden, "账号已被禁用")
	}

	// 验证密码
	if !biz.checkPassword(password, user.Password) {
		return nil, errs.New(http.StatusBadRequest, "用户名或密码错误")
	}

	return user, nil
}