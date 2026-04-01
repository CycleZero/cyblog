package category

import (
	"context"
	"cyblog/internal/common"
	"cyblog/internal/common/dto"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"cyblog/pkg/repo"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CategoryBiz struct {
	categoryRepo *repo.CategoryRepo
	logger       *log.Logger
}

func NewCategoryBiz(
	categoryRepo *repo.CategoryRepo,
	logger *log.Logger,
) *CategoryBiz {
	return &CategoryBiz{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// Create 创建分类
func (biz *CategoryBiz) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查slug是否已存在
	_, err := biz.categoryRepo.GetBySlug(ctx, req.Slug)
	if err == nil {
		return nil, errs.New(http.StatusBadRequest, "分类别名已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(http.StatusInternalServerError, "检查分类别名失败", err)
	}

	// 创建分类
	category := &model.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
	}

	err = biz.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "创建分类失败", err)
	}

	return convertToResponse(category), nil
}

// Update 更新分类
func (biz *CategoryBiz) Update(ctx context.Context, req *UpdateRequest) (*Response, error) {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查分类是否存在
	category, err := biz.categoryRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "分类不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取分类失败", err)
	}

	// 如果更新了slug，检查新slug是否已存在
	if req.Slug != "" && req.Slug != category.Slug {
		existCategory, err := biz.categoryRepo.GetBySlug(ctx, req.Slug)
		if err == nil && existCategory.ID != req.ID {
			return nil, errs.New(http.StatusBadRequest, "分类别名已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查分类别名失败", err)
		}
		category.Slug = req.Slug
	}

	// 更新字段
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = *req.ParentID
	}
	if req.Sort != nil {
		category.Sort = *req.Sort
	}

	// 保存更新
	err = biz.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "更新分类失败", err)
	}

	return convertToResponse(category), nil
}

// Delete 删除分类
func (biz *CategoryBiz) Delete(ctx context.Context, id uint) error {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查分类是否存在
	_, err := biz.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "分类不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取分类失败", err)
	}

	// 检查分类下是否有文章
	count, err := biz.categoryRepo.GetArticleCount(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "检查分类关联文章失败", err)
	}
	if count > 0 {
		return errs.New(http.StatusBadRequest, "分类下存在文章，无法删除")
	}

	// 删除分类
	err = biz.categoryRepo.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "删除分类失败", err)
	}

	return nil
}

// List 获取分类列表
func (biz *CategoryBiz) List(ctx context.Context, page, pageSize int) (*ListResponse, error) {
	categories, total, err := biz.categoryRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取分类列表失败", err)
	}

	// 转换为响应格式
	var list []*Response
	for _, category := range categories {
		list = append(list, convertToResponse(category))
	}

	return &ListResponse{
		PageInfo: dto.PageInfo{
			Total:    int(total),
			Page:     page,
			PageSize: pageSize,
		},
		List: list,
	}, nil
}

// convertToResponse 转换为响应格式
func convertToResponse(category *model.Category) *Response {
	return &Response{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		Sort:        category.Sort,
		CreatedAt:   category.CreatedAt.Format(time.DateTime),
		UpdatedAt:   category.UpdatedAt.Format(time.DateTime),
	}
}
