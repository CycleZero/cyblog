package article

import (
	"context"
	"cyblog/internal/common"
	"cyblog/internal/common/dto"
	"cyblog/internal/domain/category"
	"cyblog/internal/domain/tag"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"cyblog/pkg/repo"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArticleBiz struct {
	articleRepo  *repo.ArticleRepo
	categoryRepo *repo.CategoryRepo
	tagRepo      *repo.TagRepo
	logger       *log.Logger
}

func NewArticleBiz(
	articleRepo *repo.ArticleRepo,
	categoryRepo *repo.CategoryRepo,
	tagRepo *repo.TagRepo,
	logger *log.Logger,
) *ArticleBiz {
	return &ArticleBiz{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		logger:       logger,
	}
}

// Create 创建文章
func (biz *ArticleBiz) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	meta := common.GetRequestMetadata(ctx)
	userID := meta.UserID
	if userID == 0 {
		return nil, errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 生成文章slug（简单实现，后续可优化为拼音转换）
	slug := req.Slug
	if slug == "" {
		slug = fmt.Sprintf("%d-%s", time.Now().Unix(), strings.ReplaceAll(req.Title, " ", "-"))
	}

	// 检查slug是否唯一
	_, err := biz.articleRepo.GetBySlug(ctx, slug)
	if err == nil {
		return nil, errs.New(http.StatusBadRequest, "文章别名已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(http.StatusInternalServerError, "检查文章别名失败", err)
	}

	// 检查分类是否存在
	//var category *model.Category
	if req.CategoryID > 0 {
		_, err = biz.categoryRepo.GetByID(ctx, req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errs.New(http.StatusBadRequest, "分类不存在")
			}
			return nil, errs.Wrap(http.StatusInternalServerError, "检查分类失败", err)
		}
	}

	// 检查标签是否存在
	var tags []*model.Tag
	if len(req.TagIDs) > 0 {
		tags, err = biz.tagRepo.GetByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查标签失败", err)
		}
		if len(tags) != len(req.TagIDs) {
			return nil, errs.New(http.StatusBadRequest, "存在无效的标签ID")
		}
	}

	// 自动生成摘要
	summary := req.Summary
	if summary == "" && len(req.Content) > 200 {
		summary = string([]rune(req.Content)[:200]) + "..."
	}

	// 创建文章（直接存储 Markdown 内容，前端负责渲染）
	article := &model.Article{
		Title:      req.Title,
		Slug:       slug,
		Content:    req.Content,
		Summary:    summary,
		CoverImage: req.CoverImage,
		Status:     req.Status,
		IsTop:      req.IsTop,
		IsOriginal: req.IsOriginal,
		AuthorID:   userID,
		CategoryID: req.CategoryID,
	}

	err = biz.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "创建文章失败", err)
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		err = biz.articleRepo.UpdateArticleTags(ctx, article.ID, req.TagIDs)
		if err != nil {
			return nil, errs.Wrap(http.StatusInternalServerError, "关联标签失败", err)
		}

		// 更新标签计数
		for _, tagID := range req.TagIDs {
			_ = biz.tagRepo.IncrementCount(ctx, tagID)
		}
	}

	// 重新查询完整信息
	article, err = biz.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章详情失败", err)
	}

	return convertToResponse(article), nil
}

// Update 更新文章
func (biz *ArticleBiz) Update(ctx context.Context, req *UpdateRequest) (*Response, error) {
	meta := common.GetRequestMetadata(ctx)
	userID := meta.UserID
	if userID == 0 {
		return nil, errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查文章是否存在
	article, err := biz.articleRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "文章不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	// 检查权限：仅作者或管理员可修改
	if article.AuthorID != userID && meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限修改此文章")
	}

	// 如果更新了slug，检查是否唯一
	if req.Slug != "" && req.Slug != article.Slug {
		existArticle, err := biz.articleRepo.GetBySlug(ctx, req.Slug)
		if err == nil && existArticle.ID != req.ID {
			return nil, errs.New(http.StatusBadRequest, "文章别名已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查文章别名失败", err)
		}
		article.Slug = req.Slug
	}

	// 如果更新了分类，检查分类是否存在
	if req.CategoryID != nil && *req.CategoryID != article.CategoryID {
		if *req.CategoryID > 0 {
			_, err := biz.categoryRepo.GetByID(ctx, *req.CategoryID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errs.New(http.StatusBadRequest, "分类不存在")
				}
				return nil, errs.Wrap(http.StatusInternalServerError, "检查分类失败", err)
			}
		}
		article.CategoryID = *req.CategoryID
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
		// 重新生成摘要
		if req.Summary == "" && len(req.Content) > 200 {
			article.Summary = string([]rune(req.Content)[:200]) + "..."
		}
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	if req.Status > 0 {
		article.Status = req.Status
	}
	if req.IsTop != nil {
		article.IsTop = *req.IsTop
	}
	if req.IsOriginal != nil {
		article.IsOriginal = *req.IsOriginal
	}

	// 保存基本信息
	err = biz.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "更新文章失败", err)
	}

	// 更新标签关联
	if req.TagIDs != nil {
		// 获取原有的标签IDs
		oldTags, err := biz.articleRepo.GetArticleTags(ctx, article.ID)
		if err != nil {
			return nil, errs.Wrap(http.StatusInternalServerError, "获取原有标签失败", err)
		}
		oldTagIDMap := make(map[uint]bool)
		for _, tag := range oldTags {
			oldTagIDMap[tag.ID] = true
		}

		// 检查新标签是否都存在
		if len(req.TagIDs) > 0 {
			newTags, err := biz.tagRepo.GetByIDs(ctx, req.TagIDs)
			if err != nil {
				return nil, errs.Wrap(http.StatusInternalServerError, "检查标签失败", err)
			}
			if len(newTags) != len(req.TagIDs) {
				return nil, errs.New(http.StatusBadRequest, "存在无效的标签ID")
			}
		}

		// 更新关联
		err = biz.articleRepo.UpdateArticleTags(ctx, article.ID, req.TagIDs)
		if err != nil {
			return nil, errs.Wrap(http.StatusInternalServerError, "更新标签关联失败", err)
		}

		// 更新标签计数
		// 旧标签计数减1
		for tagID := range oldTagIDMap {
			_ = biz.tagRepo.DecrementCount(ctx, tagID)
		}
		// 新标签计数加1
		for _, tagID := range req.TagIDs {
			_ = biz.tagRepo.IncrementCount(ctx, tagID)
		}
	}

	// 重新查询完整信息
	article, err = biz.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章详情失败", err)
	}

	return convertToResponse(article), nil
}

// Delete 删除文章
func (biz *ArticleBiz) Delete(ctx context.Context, id uint) error {
	meta := common.GetRequestMetadata(ctx)
	userID := meta.UserID
	if userID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查文章是否存在
	article, err := biz.articleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "文章不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	// 检查权限：仅作者或管理员可删除
	if article.AuthorID != userID && meta.User.Role != model.RoleAdmin {
		return errs.New(http.StatusForbidden, "无权限删除此文章")
	}

	// 获取关联的标签，更新计数
	tags, err := biz.articleRepo.GetArticleTags(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "获取文章标签失败", err)
	}

	// 删除文章
	err = biz.articleRepo.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "删除文章失败", err)
	}

	// 更新标签计数
	for _, tag := range tags {
		_ = biz.tagRepo.DecrementCount(ctx, tag.ID)
	}

	return nil
}

// GetByID 获取文章详情
func (biz *ArticleBiz) GetByID(ctx context.Context, id uint, incrementViews bool) (*Response, error) {
	article, err := biz.articleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "文章不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	// 增加浏览量
	if incrementViews {
		_ = biz.articleRepo.IncrementViews(ctx, id)
		article.Views += 1
	}

	return convertToResponse(article), nil
}

// GetBySlug 根据别名获取文章
func (biz *ArticleBiz) GetBySlug(ctx context.Context, slug string, incrementViews bool) (*Response, error) {
	article, err := biz.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "文章不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	// 增加浏览量
	if incrementViews {
		_ = biz.articleRepo.IncrementViews(ctx, article.ID)
		article.Views += 1
	}

	return convertToResponse(article), nil
}

// List 获取文章列表
func (biz *ArticleBiz) List(ctx context.Context, query *ListQuery) (*ListResponse, error) {
	articles, total, err := biz.articleRepo.List(ctx, query)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章列表失败", err)
	}

	// 转换为响应格式
	var list []*Response
	for _, article := range articles {
		list = append(list, convertToResponse(article))
	}

	return &ListResponse{
		PageInfo: dto.PageInfo{
			Total:    int(total),
			Page:     query.Page,
			PageSize: query.PageSize,
		},
		List: list,
	}, nil
}

// Like 点赞文章
func (biz *ArticleBiz) Like(ctx context.Context, articleID uint) error {
	meta := common.GetRequestMetadata(ctx)
	userID := meta.UserID
	if userID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查文章是否存在
	_, err := biz.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "文章不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	// TODO: 后续添加点赞记录表，防止重复点赞
	err = biz.articleRepo.IncrementLikes(ctx, articleID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "点赞失败", err)
	}

	return nil
}

// Unlike 取消点赞
func (biz *ArticleBiz) Unlike(ctx context.Context, articleID uint) error {
	meta := common.GetRequestMetadata(ctx)
	userID := meta.UserID
	if userID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查文章是否存在
	_, err := biz.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "文章不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取文章失败", err)
	}

	err = biz.articleRepo.DecrementLikes(ctx, articleID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "取消点赞失败", err)
	}

	return nil
}

// convertToResponse 转换为响应格式
func convertToResponse(article *model.Article) *Response {
	resp := &Response{
		ID:         article.ID,
		Title:      article.Title,
		Slug:       article.Slug,
		Content:    article.Content,
		Summary:    article.Summary,
		CoverImage: article.CoverImage,
		Status:     article.Status,
		Views:      article.Views,
		Likes:      article.Likes,
		IsTop:      article.IsTop,
		IsOriginal: article.IsOriginal,
		CreatedAt:  article.CreatedAt,
		UpdatedAt:  article.UpdatedAt,
	}

	// 转换作者信息
	if article.Author != nil {
		resp.Author = &struct {
			ID     uint   `json:"id"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		}{
			ID:     article.Author.ID,
			Name:   article.Author.Name,
			Avatar: article.Author.Avatar,
		}
	}

	// 转换分类信息
	if article.Category != nil {
		resp.Category = &category.Response{
			ID:          article.Category.ID,
			Name:        article.Category.Name,
			Slug:        article.Category.Slug,
			Description: article.Category.Description,
			ParentID:    article.Category.ParentID,
			Sort:        article.Category.Sort,
			CreatedAt:   article.Category.CreatedAt.Format(time.DateTime),
			UpdatedAt:   article.Category.UpdatedAt.Format(time.DateTime),
		}
	}

	// 转换标签信息
	if article.Tags != nil && len(article.Tags) > 0 {
		for _, t := range article.Tags {
			resp.Tags = append(resp.Tags, &tag.Response{
				ID:        t.ID,
				Name:      t.Name,
				Slug:      t.Slug,
				Color:     t.Color,
				Count:     t.Count,
				CreatedAt: t.CreatedAt.Format(time.DateTime),
				UpdatedAt: t.UpdatedAt.Format(time.DateTime),
			})
		}
	}

	return resp
}

// AdminList 管理端获取文章列表
func (biz *ArticleBiz) AdminList(ctx context.Context, query *ListQuery) (*ListResponse, error) {
	articles, total, err := biz.articleRepo.List(ctx, query)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取文章列表失败", err)
	}
	var list []*Response
	for _, article := range articles {
		list = append(list, convertToResponse(article))
	}
	return &ListResponse{
		PageInfo: dto.PageInfo{
			Total:    int(total),
			Page:     query.Page,
			PageSize: query.PageSize,
		},
		List: list,
	}, nil
}

// SetTop 置顶/取消置顶文章
func (biz *ArticleBiz) SetTop(ctx context.Context, id uint, isTop bool) error {
	article, err := biz.articleRepo.GetByID(ctx, id)
	if err != nil {
		return errs.New(http.StatusNotFound, "文章不存在")
	}
	article.IsTop = isTop
	return biz.articleRepo.Update(ctx, article)
}

// BatchDelete 批量删除文章
func (biz *ArticleBiz) BatchDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		_ = biz.Delete(ctx, id)
	}
	return nil
}

// BatchUpdateStatus 批量更新文章状态
func (biz *ArticleBiz) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	for _, id := range ids {
		article, err := biz.articleRepo.GetByID(ctx, id)
		if err != nil {
			continue
		}
		article.Status = status
		_ = biz.articleRepo.Update(ctx, article)
	}
	return nil
}

// GetDashboard 获取仪表盘数据
func (biz *ArticleBiz) GetDashboard(ctx context.Context) (*DashboardResponse, error) {
	articleCount, _ := biz.articleRepo.Count(ctx)
	todayViews, _ := biz.articleRepo.GetTodayViews(ctx)
	recentArticles, _ := biz.articleRepo.GetRecent(ctx, 10)
	hotArticles, _ := biz.articleRepo.GetHot(ctx, 10)

	recentList := make([]RecentArticleSimple, 0, len(recentArticles))
	for _, a := range recentArticles {
		recentList = append(recentList, RecentArticleSimple{
			ID:     a.ID,
			Title:  a.Title,
			Status: a.Status,
		})
	}

	hotList := make([]HotArticleSimple, 0, len(hotArticles))
	for _, a := range hotArticles {
		hotList = append(hotList, HotArticleSimple{
			ID:    a.ID,
			Title: a.Title,
			Views: a.Views,
		})
	}

	return &DashboardResponse{
		ArticleCount:   articleCount,
		TodayViews:     todayViews,
		RecentArticles: recentList,
		HotArticles:    hotList,
	}, nil
}
