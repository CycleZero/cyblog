package repo

import (
	"context"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"strings"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	data   *infra.Data
	logger *log.Logger
}

func NewArticleRepo(data *infra.Data, logger *log.Logger) *ArticleRepo {
	err := data.DB.AutoMigrate(&model.Article{})
	if err != nil {
		log.GetLogger().Sugar().Errorf("AutoMigrate 创建文章表失败：%v", err)
		panic(err)
	}
	return &ArticleRepo{
		data:   data,
		logger: logger,
	}
}

func (r *ArticleRepo) Create(ctx context.Context, article *model.Article) error {
	return r.data.DB.WithContext(ctx).Create(article).Error
}

func (r *ArticleRepo) Update(ctx context.Context, article *model.Article) error {
	return r.data.DB.WithContext(ctx).Save(article).Error
}

func (r *ArticleRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Article{}, id).Error
}

func (r *ArticleRepo) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.data.DB.WithContext(ctx).Preload("Author").Preload("Category").Preload("Tags").Where("id = ?", id).First(&article).Error
	return &article, err
}

func (r *ArticleRepo) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	err := r.data.DB.WithContext(ctx).Preload("Author").Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&article).Error
	return &article, err
}

func (r *ArticleRepo) List(ctx context.Context, query interface{}) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.Article{})

	// 动态条件查询
	if q, ok := query.(interface{ GetKeyword() string }); ok && q.GetKeyword() != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+q.GetKeyword()+"%", "%"+q.GetKeyword()+"%")
	}
	if q, ok := query.(interface{ GetCategoryID() uint }); ok && q.GetCategoryID() > 0 {
		db = db.Where("category_id = ?", q.GetCategoryID())
	}
	if q, ok := query.(interface{ GetTagID() uint }); ok && q.GetTagID() > 0 {
		db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").Where("article_tags.tag_id = ?", q.GetTagID())
	}
	if q, ok := query.(interface{ GetStatus() int }); ok && q.GetStatus() > 0 {
		db = db.Where("status = ?", q.GetStatus())
	}
	if q, ok := query.(interface{ GetAuthorID() uint }); ok && q.GetAuthorID() > 0 {
		db = db.Where("author_id = ?", q.GetAuthorID())
	}
	if q, ok := query.(interface{ GetIsTop() *bool }); ok && q.GetIsTop() != nil {
		db = db.Where("is_top = ?", *q.GetIsTop())
	}

	// 统计总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序
	if q, ok := query.(interface{ GetSortBy() string }); ok && q.GetSortBy() != "" {
		sortBy := q.GetSortBy()
		sortOrder := "desc"
		if q, ok := query.(interface{ GetSortOrder() string }); ok && strings.ToLower(q.GetSortOrder()) == "asc" {
			sortOrder = "asc"
		}
		db = db.Order(sortBy + " " + sortOrder)
	} else {
		db = db.Order("is_top desc, created_at desc")
	}

	// 分页
	if q, ok := query.(interface{ GetPage() int }); ok && q.GetPage() > 0 {
		page := q.GetPage()
		pageSize := 10
		if q, ok := query.(interface{ GetPageSize() int }); ok && q.GetPageSize() > 0 {
			pageSize = q.GetPageSize()
		}
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}

	// 关联查询
	err = db.Preload("Author").Preload("Category").Preload("Tags").Find(&articles).Error
	return articles, total, err
}

func (r *ArticleRepo) UpdateArticleTags(ctx context.Context, articleID uint, tagIDs []uint) error {
	// 先删除原有关联
	err := r.data.DB.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error
	if err != nil {
		return err
	}

	// 新增关联
	if len(tagIDs) == 0 {
		return nil
	}

	var articleTags []*model.ArticleTag
	for _, tagID := range tagIDs {
		articleTags = append(articleTags, &model.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}

	return r.data.DB.WithContext(ctx).Create(&articleTags).Error
}

func (r *ArticleRepo) GetArticleTags(ctx context.Context, articleID uint) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := r.data.DB.WithContext(ctx).Joins("JOIN article_tags ON article_tags.tag_id = tags.id").Where("article_tags.article_id = ?", articleID).Find(&tags).Error
	return tags, err
}

func (r *ArticleRepo) IncrementViews(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (r *ArticleRepo) IncrementLikes(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1")).Error
}

func (r *ArticleRepo) DecrementLikes(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes - 1")).Error
}
