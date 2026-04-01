package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"

	"gorm.io/gorm"
)

type TagRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewTagRepo(data *infra.Data, logger *log.Logger) *TagRepo {
	err := data.DB.AutoMigrate(&model.Tag{}, &model.ArticleTag{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建标签表失败：%v", err)
		panic(err)
	}
	return &TagRepo{
		data:   data,
		logger: logger,
	}
}

func (r *TagRepo) Create(ctx context.Context, tag *model.Tag) error {
	return r.data.DB.WithContext(ctx).Create(tag).Error
}

func (r *TagRepo) Update(ctx context.Context, tag *model.Tag) error {
	return r.data.DB.WithContext(ctx).Save(tag).Error
}

func (r *TagRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Tag{}, id).Error
}

func (r *TagRepo) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.data.DB.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	return &tag, err
}

func (r *TagRepo) GetBySlug(ctx context.Context, slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.data.DB.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	return &tag, err
}

func (r *TagRepo) List(ctx context.Context, page, pageSize int) ([]*model.Tag, int64, error) {
	var tags []*model.Tag
	var total int64

	err := r.data.DB.WithContext(ctx).Model(&model.Tag{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.data.DB.WithContext(ctx).Order("count desc, created_at desc").Offset(offset).Limit(pageSize).Find(&tags).Error
	return tags, total, err
}

func (r *TagRepo) GetByIDs(ctx context.Context, ids []uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.data.DB.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *TagRepo) IncrementCount(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("count", gorm.Expr("count + 1")).Error
}

func (r *TagRepo) DecrementCount(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Tag{}).Where("id = ?", id).UpdateColumn("count", gorm.Expr("count - 1")).Error
}
