package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
)

type CategoryRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewCategoryRepo(data *infra.Data, logger *log.Logger) *CategoryRepo {
	err := data.DB.AutoMigrate(&model.Category{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建分类表失败：%v", err)
		panic(err)
	}
	return &CategoryRepo{
		data:   data,
		logger: logger,
	}
}

func (r *CategoryRepo) Create(ctx context.Context, category *model.Category) error {
	return r.data.DB.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepo) Update(ctx context.Context, category *model.Category) error {
	return r.data.DB.WithContext(ctx).Save(category).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Category{}, id).Error
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	err := r.data.DB.WithContext(ctx).Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	err := r.data.DB.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	return &category, err
}

func (r *CategoryRepo) List(ctx context.Context, page, pageSize int) ([]*model.Category, int64, error) {
	var categories []*model.Category
	var total int64

	err := r.data.DB.WithContext(ctx).Model(&model.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.data.DB.WithContext(ctx).Order("sort asc, created_at desc").Offset(offset).Limit(pageSize).Find(&categories).Error
	return categories, total, err
}

func (r *CategoryRepo) GetArticleCount(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.data.DB.WithContext(ctx).Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
