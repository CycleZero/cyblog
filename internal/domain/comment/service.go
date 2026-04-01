package comment

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentService struct {
	common.BaseService
	biz    *CommentBiz
	logger *log.Logger
}

func NewCommentService(biz *CommentBiz, logger *log.Logger) *CommentService {
	return &CommentService{
		biz:    biz,
		logger: logger,
	}
}

// Create 创建评论
// @Summary 创建评论
// @Description 对文章发表评论，需要登录
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateRequest true "创建评论请求参数"
// @Success 200 {object} common.Response{data=Response} "创建成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments [post]
func (s *CommentService) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Create(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// Update 更新评论
// @Summary 更新评论
// @Description 更新评论内容，仅作者或管理员可操作
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param request body UpdateRequest true "更新评论请求参数"
// @Success 200 {object} common.Response{data=Response} "更新成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 404 {object} common.Response "评论不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments/{id} [put]
func (s *CommentService) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}
	req.ID = uint(id)

	resp, err := s.biz.Update(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// Delete 删除评论
// @Summary 删除评论
// @Description 删除指定评论，仅作者或管理员可操作
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "无效的评论ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 404 {object} common.Response "评论不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments/{id} [delete]
func (s *CommentService) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	err = s.biz.Delete(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, nil)
}

// List 获取文章评论列表
// @Summary 获取文章评论列表
// @Description 分页获取指定文章的评论列表
// @Tags 评论
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param parent_id query int false "父评论ID，不传则获取一级评论"
// @Param sort_by query string false "排序字段：created_at, likes" default(created_at)
// @Param sort_order query string false "排序方式：asc, desc" default(desc)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments [get]
func (s *CommentService) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.List(c, &query)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// GetReplies 获取评论的回复列表
// @Summary 获取评论的回复列表
// @Description 分页获取指定评论的回复列表
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 400 {object} common.Response "无效的评论ID"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments/{id}/replies [get]
func (s *CommentService) GetReplies(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := s.biz.GetReplies(c, uint(id), page, pageSize)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// GetCommentCount 获取文章评论数
// @Summary 获取文章评论数
// @Description 获取指定文章的评论数量
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} common.Response{data=ArticleCommentCountResponse} "获取成功"
// @Failure 400 {object} common.Response "无效的文章ID"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/{id}/comment-count [get]
func (s *CommentService) GetCommentCount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的文章ID"))
		return
	}

	resp, err := s.biz.GetCommentCount(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// Like 点赞评论
// @Summary 点赞评论
// @Description 对指定评论点赞，需要登录
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.Response "点赞成功"
// @Failure 400 {object} common.Response "无效的评论ID或已点赞"
// @Failure 401 {object} common.Response "未认证"
// @Failure 404 {object} common.Response "评论不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments/{id}/like [post]
func (s *CommentService) Like(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	err = s.biz.Like(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, nil)
}

// Unlike 取消点赞
// @Summary 取消点赞
// @Description 取消对指定评论的点赞，需要登录
// @Tags 评论
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.Response "取消点赞成功"
// @Failure 400 {object} common.Response "无效的评论ID或未点赞"
// @Failure 401 {object} common.Response "未认证"
// @Failure 404 {object} common.Response "评论不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/comments/{id}/like [delete]
func (s *CommentService) Unlike(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	err = s.biz.Unlike(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, nil)
}

// AdminList 管理端获取评论列表
// @Summary 管理端获取评论列表
// @Description 分页获取所有评论列表，需要管理员权限
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param keyword query string false "搜索关键词"
// @Param article_id query int false "文章ID"
// @Param user_id query int false "用户ID"
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方式" default(desc)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/admin/comments [get]
func (s *CommentService) AdminList(c *gin.Context) {
	var query AdminListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.AdminList(c, &query)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, resp)
}

// AdminDelete 管理端删除评论
// @Summary 管理端删除评论
// @Description 删除指定评论，需要管理员权限
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "无效的评论ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 404 {object} common.Response "评论不存在"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/admin/comments/{id} [delete]
func (s *CommentService) AdminDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的评论ID"))
		return
	}

	err = s.biz.Delete(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "服务器内部错误", err))
		return
	}

	common.Success(c, nil)
}
