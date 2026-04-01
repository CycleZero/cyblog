package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"

	"gorm.io/gorm"
)

type CommentRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewCommentRepo(data *infra.Data, logger *log.Logger) *CommentRepo {
	err := data.DB.AutoMigrate(&model.Comment{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建评论表失败：%v", err)
		panic(err)
	}
	return &CommentRepo{
		data:   data,
		logger: logger,
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	return r.data.DB.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) error {
	return r.data.DB.WithContext(ctx).Save(comment).Error
}

func (r *CommentRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (r *CommentRepo) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.data.DB.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	return &comment, err
}

func (r *CommentRepo) List(ctx context.Context, query interface{}) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.Comment{})

	if q, ok := query.(interface{ GetArticleID() uint }); ok {
		if articleID := q.GetArticleID(); articleID > 0 {
			db = db.Where("article_id = ?", articleID)
		}
	}

	if q, ok := query.(interface{ GetParentID() *uint }); ok {
		if parentID := q.GetParentID(); parentID != nil {
			db = db.Where("parent_id = ?", *parentID)
		} else {
			db = db.Where("parent_id = 0")
		}
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := 0
	pageSize := 10

	if q, ok := query.(interface{ GetPage() int }); ok {
		if page := q.GetPage(); page > 0 {
			offset = (page - 1) * pageSize
		}
	}

	if q, ok := query.(interface{ GetPageSize() int }); ok {
		if ps := q.GetPageSize(); ps > 0 {
			pageSize = ps
		}
	}

	sortBy := "created_at"
	sortOrder := "desc"

	if q, ok := query.(interface{ GetSortBy() string }); ok {
		if sb := q.GetSortBy(); sb != "" {
			sortBy = sb
		}
	}

	if q, ok := query.(interface{ GetSortOrder() string }); ok {
		if so := q.GetSortOrder(); so != "" {
			sortOrder = so
		}
	}

	err = db.Preload("User").
		Order(gorm.Expr("? ?", gorm.Expr(sortBy), gorm.Expr(sortOrder))).
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepo) GetCommentCountByArticleIDs(ctx context.Context, articleIDs []uint) (map[uint]int, error) {
	var results []struct {
		ArticleID uint
		Count     int
	}

	err := r.data.DB.WithContext(ctx).
		Model(&model.Comment{}).
		Select("article_id, count(*) as count").
		Where("article_id IN ?", articleIDs).
		Group("article_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]int)
	for _, r := range results {
		countMap[r.ArticleID] = r.Count
	}

	return countMap, nil
}

func (r *CommentRepo) GetCommentCountByArticleID(ctx context.Context, articleID uint) (int, error) {
	var count int64
	err := r.data.DB.WithContext(ctx).Model(&model.Comment{}).Where("article_id = ?", articleID).Count(&count).Error
	return int(count), err
}

func (r *CommentRepo) GetReplies(ctx context.Context, parentID uint, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.Comment{}).Where("parent_id = ?", parentID)

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = db.Preload("User").Preload("Parent").Preload("Parent.User").
		Order("created_at asc").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepo) IncrementLikes(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

func (r *CommentRepo) DecrementLikes(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes - 1")).Error
}

func (r *CommentRepo) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return r.data.DB.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.Comment{}).Error
}
