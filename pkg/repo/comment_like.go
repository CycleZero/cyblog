package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
)

type CommentLikeRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewCommentLikeRepo(data *infra.Data, logger *log.Logger) *CommentLikeRepo {
	err := data.DB.AutoMigrate(&model.CommentLike{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建评论点赞记录表失败：%v", err)
		panic(err)
	}
	return &CommentLikeRepo{
		data:   data,
		logger: logger,
	}
}

func (r *CommentLikeRepo) Create(ctx context.Context, like *model.CommentLike) error {
	return r.data.DB.WithContext(ctx).Create(like).Error
}

func (r *CommentLikeRepo) Delete(ctx context.Context, commentID, userID uint) error {
	return r.data.DB.WithContext(ctx).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&model.CommentLike{}).Error
}

func (r *CommentLikeRepo) GetByCommentAndUser(ctx context.Context, commentID, userID uint) (*model.CommentLike, error) {
	var like model.CommentLike
	err := r.data.DB.WithContext(ctx).Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like).Error
	return &like, err
}

func (r *CommentLikeRepo) Exists(ctx context.Context, commentID, userID uint) (bool, error) {
	var count int64
	err := r.data.DB.WithContext(ctx).Model(&model.CommentLike{}).Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error
	return count > 0, err
}

func (r *CommentLikeRepo) GetByUser(ctx context.Context, userID uint, page, pageSize int) ([]*model.CommentLike, int64, error) {
	var likes []*model.CommentLike
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.CommentLike{}).Where("user_id = ?", userID)

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = db.Preload("Comment").Preload("Comment.User").
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&likes).Error

	return likes, total, err
}

func (r *CommentLikeRepo) GetLikedCommentIDs(ctx context.Context, userID uint, commentIDs []uint) (map[uint]bool, error) {
	var results []struct {
		CommentID uint
	}

	err := r.data.DB.WithContext(ctx).
		Model(&model.CommentLike{}).
		Select("comment_id").
		Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	likedMap := make(map[uint]bool)
	for _, r := range results {
		likedMap[r.CommentID] = true
	}

	return likedMap, nil
}
