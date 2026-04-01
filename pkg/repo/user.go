package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
)

type UserRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewUserRepo(data *infra.Data, logger *log.Logger) *UserRepo {
	err := data.DB.AutoMigrate(&model.User{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建表失败：%v", err)
		panic(err)
	}
	return &UserRepo{
		data:   data,
		logger: logger,
	}
}

func (r *UserRepo) GetUserById(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.data.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.data.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := r.data.DB.WithContext(ctx).Where("name = ?", name).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.data.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	return r.data.DB.WithContext(ctx).Save(user).Error
}
