package comment

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

	"gorm.io/gorm"
)

type CommentBiz struct {
	commentRepo     *repo.CommentRepo
	commentLikeRepo *repo.CommentLikeRepo
	logger          *log.Logger
}

func NewCommentBiz(
	commentRepo *repo.CommentRepo,
	commentLikeRepo *repo.CommentLikeRepo,
	logger *log.Logger,
) *CommentBiz {
	return &CommentBiz{
		commentRepo:     commentRepo,
		commentLikeRepo: commentLikeRepo,
		logger:          logger,
	}
}

// Create 创建评论
func (biz *CommentBiz) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	meta := common.GetRequestMetadata(ctx)
	if meta.UserID == 0 {
		return nil, errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 创建评论
	comment := &model.Comment{
		ArticleID: req.ArticleID,
		UserID:    meta.UserID,
		ParentID:  req.ParentID,
		Content:   req.Content,
		IP:        meta.RealIp,
		UserAgent: meta.UserAgent,
	}

	err := biz.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "创建评论失败", err)
	}

	return biz.GetByID(ctx, comment.ID)
}

// Update 更新评论（仅作者或管理员）
func (biz *CommentBiz) Update(ctx context.Context, req *UpdateRequest) (*Response, error) {
	meta := common.GetRequestMetadata(ctx)
	if meta.UserID == 0 {
		return nil, errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查评论是否存在
	comment, err := biz.commentRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "评论不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取评论失败", err)
	}

	// 检查权限：仅作者或管理员可操作
	if comment.UserID != meta.UserID && meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 更新内容
	comment.Content = req.Content
	err = biz.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "更新评论失败", err)
	}

	return biz.GetByID(ctx, comment.ID)
}

// Delete 删除评论（仅作者或管理员）
func (biz *CommentBiz) Delete(ctx context.Context, id uint) error {
	meta := common.GetRequestMetadata(ctx)
	if meta.UserID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查评论是否存在
	comment, err := biz.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "评论不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取评论失败", err)
	}

	// 检查权限：仅作者或管理员可操作
	if comment.UserID != meta.UserID && meta.User.Role != model.RoleAdmin {
		return errs.New(http.StatusForbidden, "无权限操作")
	}

	// 删除评论
	err = biz.commentRepo.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "删除评论失败", err)
	}

	return nil
}

// GetByID 获取评论详情
func (biz *CommentBiz) GetByID(ctx context.Context, id uint) (*Response, error) {
	comment, err := biz.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, "评论不存在")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, "获取评论失败", err)
	}

	return biz.convertToResponse(ctx, comment, true), nil
}

// List 获取文章评论列表
func (biz *CommentBiz) List(ctx context.Context, query *ListQuery) (*ListResponse, error) {
	comments, total, err := biz.commentRepo.List(ctx, query)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取评论列表失败", err)
	}

	// 批量获取点赞状态
	commentIDs := make([]uint, 0, len(comments))
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}

	var list []*Response
	for _, comment := range comments {
		list = append(list, biz.convertToResponse(ctx, comment, false))
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

// GetReplies 获取评论的回复列表
func (biz *CommentBiz) GetReplies(ctx context.Context, parentID uint, page, pageSize int) (*ListResponse, error) {
	comments, total, err := biz.commentRepo.GetReplies(ctx, parentID, page, pageSize)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取回复列表失败", err)
	}

	var list []*Response
	for _, comment := range comments {
		list = append(list, biz.convertToResponse(ctx, comment, true))
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

// AdminList 管理端获取评论列表
func (biz *CommentBiz) AdminList(ctx context.Context, query *AdminListQuery) (*ListResponse, error) {
	meta := common.GetRequestMetadata(ctx)
	if meta.User.Role != model.RoleAdmin {
		return nil, errs.New(http.StatusForbidden, "无权限操作")
	}

	// 构造查询条件
	type Query struct {
		dto.PageParam
		Keyword   string
		ArticleID uint
		UserID    uint
		SortBy    string
		SortOrder string
	}

	q := &Query{
		PageParam: query.PageParam,
		Keyword:   query.Keyword,
		ArticleID: query.ArticleID,
		UserID:    query.UserID,
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	comments, total, err := biz.commentRepo.List(ctx, q)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取评论列表失败", err)
	}

	var list []*Response
	for _, comment := range comments {
		list = append(list, biz.convertToResponse(ctx, comment, false))
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

// Like 评论点赞
func (biz *CommentBiz) Like(ctx context.Context, commentID uint) error {
	meta := common.GetRequestMetadata(ctx)
	if meta.UserID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查评论是否存在
	_, err := biz.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "评论不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取评论失败", err)
	}

	// 检查是否已点赞
	exists, err := biz.commentLikeRepo.Exists(ctx, commentID, meta.UserID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "检查点赞状态失败", err)
	}
	if exists {
		return errs.New(http.StatusBadRequest, "已点赞该评论")
	}

	// 创建点赞记录
	like := &model.CommentLike{
		CommentID: commentID,
		UserID:    meta.UserID,
		IP:        meta.RealIp,
		UserAgent: meta.UserAgent,
	}

	err = biz.commentLikeRepo.Create(ctx, like)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "点赞失败", err)
	}

	// 增加评论点赞数
	err = biz.commentRepo.IncrementLikes(ctx, commentID)
	if err != nil {
		// 回滚点赞记录
		_ = biz.commentLikeRepo.Delete(ctx, commentID, meta.UserID)
		return errs.Wrap(http.StatusInternalServerError, "更新点赞数失败", err)
	}

	return nil
}

// Unlike 取消点赞
func (biz *CommentBiz) Unlike(ctx context.Context, commentID uint) error {
	meta := common.GetRequestMetadata(ctx)
	if meta.UserID == 0 {
		return errs.New(http.StatusUnauthorized, "请先登录")
	}

	// 检查评论是否存在
	_, err := biz.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, "评论不存在")
		}
		return errs.Wrap(http.StatusInternalServerError, "获取评论失败", err)
	}

	// 检查是否已点赞
	exists, err := biz.commentLikeRepo.Exists(ctx, commentID, meta.UserID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "检查点赞状态失败", err)
	}
	if !exists {
		return errs.New(http.StatusBadRequest, "未点赞该评论")
	}

	// 删除点赞记录
	err = biz.commentLikeRepo.Delete(ctx, commentID, meta.UserID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "取消点赞失败", err)
	}

	// 减少评论点赞数
	err = biz.commentRepo.DecrementLikes(ctx, commentID)
	if err != nil {
		return errs.Wrap(http.StatusInternalServerError, "更新点赞数失败", err)
	}

	return nil
}

// GetCommentCount 获取文章评论数
func (biz *CommentBiz) GetCommentCount(ctx context.Context, articleID uint) (*ArticleCommentCountResponse, error) {
	count, err := biz.commentRepo.GetCommentCountByArticleID(ctx, articleID)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "获取评论数失败", err)
	}

	return &ArticleCommentCountResponse{
		ArticleID:    articleID,
		CommentCount: count,
	}, nil
}

// BatchGetCommentCount 批量获取文章评论数
func (biz *CommentBiz) BatchGetCommentCount(ctx context.Context, articleIDs []uint) (map[uint]int, error) {
	countMap, err := biz.commentRepo.GetCommentCountByArticleIDs(ctx, articleIDs)
	if err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, "批量获取评论数失败", err)
	}

	return countMap, nil
}

// convertToResponse 转换为响应格式
func (biz *CommentBiz) convertToResponse(ctx context.Context, comment *model.Comment, loadParent bool) *Response {
	meta := common.GetRequestMetadata(ctx)

	resp := &Response{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Likes:     comment.Likes,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	// 设置用户信息
	if comment.User != nil {
		resp.User = &dto.User{
			ID:     comment.User.ID,
			Name:   comment.User.Name,
			Avatar: comment.User.Avatar,
		}
	}

	// 设置被回复的用户信息
	if loadParent && comment.Parent != nil && comment.Parent.User != nil {
		resp.ReplyTo = &dto.User{
			ID:   comment.Parent.User.ID,
			Name: comment.Parent.User.Name,
		}
	}

	// 检查当前用户是否已点赞
	if meta.UserID > 0 {
		isLiked, _ := biz.commentLikeRepo.Exists(ctx, comment.ID, meta.UserID)
		resp.IsLiked = isLiked
	}

	return resp
}
