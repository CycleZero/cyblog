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

// ListUsers 获取用户列表（管理端）
func (r *UserRepo) ListUsers(ctx context.Context, query interface{}) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.User{})

	// 关键词搜索
	if q, ok := query.(interface{ GetKeyword() string }); ok && q.GetKeyword() != "" {
		keyword := "%" + q.GetKeyword() + "%"
		db = db.Where("name LIKE ? OR email LIKE ?", keyword, keyword)
	}

	// 角色筛选
	if q, ok := query.(interface{ GetRole() string }); ok && q.GetRole() != "" {
		db = db.Where("role = ?", q.GetRole())
	}

	// 状态筛选
	if q, ok := query.(interface{ GetStatus() int }); ok && q.GetStatus() > 0 {
		db = db.Where("status = ?", q.GetStatus())
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	sortBy := "created_at"
	sortOrder := "desc"
	if q, ok := query.(interface{ GetSortBy() string }); ok && q.GetSortBy() != "" {
		sortBy = q.GetSortBy()
	}
	if q, ok := query.(interface{ GetSortOrder() string }); ok && q.GetSortOrder() != "" {
		sortOrder = q.GetSortOrder()
	}
	db = db.Order(sortBy + " " + sortOrder)

	// 分页
	page := 1
	pageSize := 10
	if q, ok := query.(interface{ GetPage() int }); ok && q.GetPage() > 0 {
		page = q.GetPage()
	}
	if q, ok := query.(interface{ GetPageSize() int }); ok && q.GetPageSize() > 0 {
		pageSize = q.GetPageSize()
	}
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)

	if err := db.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
