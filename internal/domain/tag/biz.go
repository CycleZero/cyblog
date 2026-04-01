package tag

import (
	"context"
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"cyblog/pkg/model"
	"cyblog/pkg/repo"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type TagBiz struct {
	tagRepo *repo.TagRepo
	logger  *log.Logger
}

func NewTagBiz(
	tagRepo *repo.TagRepo,
	logger *log.Logger,
) *TagBiz {
	return &TagBiz{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

// Create 创建标签
func (biz *TagBiz) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查slug是否已存在
	_, err := biz.tagRepo.GetBySlug(ctx, req.Slug)
	if err == nil {
		return nil, errs.New(http.StatusBadRequest, "标签别名已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.Wrap(http.StatusInternalServerError, "检查标签别名失败", err)
	}

	// 默认颜色处理
	color := req.Color
	if color == "" {
		color = "#165DFF"
	}

	// 创建标签
	tag := &model.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: color,
		Count: 0,
	}

	err = biz.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "创建标签失败", err)
	}

	return convertToResponse(tag), nil
}

// Update 更新标签
func (biz *TagBiz) Update(ctx context.Context, req *UpdateRequest) (*Response, error) {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查标签是否存在
	tag, err := biz.tagRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "标签不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取标签失败", err)
	}

	// 如果更新了slug，检查新slug是否已存在
	if req.Slug != "" && req.Slug != tag.Slug {
		existTag, err := biz.tagRepo.GetBySlug(ctx, req.Slug)
		if err == nil && existTag.ID != req.ID {
			return nil, errs.New(http.StatusBadRequest, "标签别名已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(http.StatusInternalServerError, "检查标签别名失败", err)
		}
		tag.Slug = req.Slug
	}

	// 更新字段
	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	// 保存更新
	err = biz.tagRepo.Update(ctx, tag)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "更新标签失败", err)
	}

	return convertToResponse(tag), nil
}

// Delete 删除标签
func (biz *TagBiz) Delete(ctx context.Context, id uint) error {
	// 检查权限：仅管理员可操作
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return errs.New(http.StatusForbidden, "无权限操作")
	}

	// 检查标签是否存在
	tag, err := biz.tagRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "标签不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取标签失败", err)
	}

	// 检查标签是否有关联文章
	if tag.Count > 0 {
		return errs.New(http.StatusBadRequest, "标签下存在关联文章，无法删除")
	}

	// 删除标签
	err = biz.tagRepo.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "删除标签失败", err)
	}

	return nil
}

// List 获取标签列表
func (biz *TagBiz) List(ctx context.Context, page, pageSize int) (*ListResponse, error) {
	tags, total, err := biz.tagRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取标签列表失败", err)
	}

	// 转换为响应格式
	var list []*Response
	for _, tag := range tags {
		list = append(list, convertToResponse(tag))
	}

	return &ListResponse{
		PageInfo: common.PageInfo{
			Total:    int(total),
			Page:     page,
			PageSize: pageSize,
		},
		List: list,
	}, nil
}

// convertToResponse 转换为响应格式
func convertToResponse(tag *model.Tag) *Response {
	return &Response{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		Color:     tag.Color,
		Count:     tag.Count,
		CreatedAt: tag.CreatedAt.Format(time.DateTime),
		UpdatedAt: tag.UpdatedAt.Format(time.DateTime),
	}
}
