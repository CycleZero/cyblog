package article

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleService struct {
	common.BaseService
	biz    *ArticleBiz
	logger *log.Logger
}

func NewArticleService(biz *ArticleBiz, logger *log.Logger) *ArticleService {
	return &ArticleService{
		biz:    biz,
		logger: logger,
	}
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建新文章，需要认证
// @Tags 文章
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateRequest true "创建文章请求参数"
// @Success 200 {object} common.Response{data=Response} "创建成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles [post]
func (s *ArticleService) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Create(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "创建文章失败", err))
		return
	}

	common.Success(c, resp)
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新文章信息，需要认证
// @Tags 文章
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateRequest true "更新文章请求参数"
// @Success 200 {object} common.Response{data=Response} "更新成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles [put]
func (s *ArticleService) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Update(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新文章失败", err))
		return
	}

	common.Success(c, resp)
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除指定文章，需要认证
// @Tags 文章
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "无效的文章ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/{id} [delete]
func (s *ArticleService) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的文章ID"))
		return
	}

	err = s.biz.Delete(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "删除文章失败", err))
		return
	}

	common.Success(c, nil)
}

// GetByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据文章ID获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param increment_views query bool false "是否增加浏览量" default(true)
// @Success 200 {object} common.Response{data=Response} "获取成功"
// @Failure 400 {object} common.Response "无效的文章ID"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/{id} [get]
func (s *ArticleService) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的文章ID"))
		return
	}

	// 是否增加浏览量（默认增加，后台编辑的时候可以传increment_views=false不增加）
	incrementViews := true
	if inc := c.Query("increment_views"); inc == "false" {
		incrementViews = false
	}

	resp, err := s.biz.GetByID(c, uint(id), incrementViews)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取文章详情失败", err))
		return
	}

	common.Success(c, resp)
}

// GetBySlug 根据别名获取文章
// @Summary 根据别名获取文章
// @Description 根据文章别名获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param slug path string true "文章别名"
// @Param increment_views query bool false "是否增加浏览量" default(true)
// @Success 200 {object} common.Response{data=Response} "获取成功"
// @Failure 400 {object} common.Response "文章别名不能为空"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/slug/{slug} [get]
func (s *ArticleService) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		common.Error(c, errs.New(http.StatusBadRequest, "文章别名不能为空"))
		return
	}

	// 是否增加浏览量
	incrementViews := true
	if inc := c.Query("increment_views"); inc == "false" {
		incrementViews = false
	}

	resp, err := s.biz.GetBySlug(c, slug, incrementViews)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取文章失败", err))
		return
	}

	common.Success(c, resp)
}

// List 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取文章列表，支持多种筛选条件
// @Tags 文章
// @Accept json
// @Produce json
// @Param keyword query string false "关键词搜索"
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param status query int false "状态"
// @Param author_id query int false "作者ID"
// @Param is_top query bool false "是否置顶"
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方式" default(desc)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles [get]
func (s *ArticleService) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.List(c, &query)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取文章列表失败", err))
		return
	}

	common.Success(c, resp)
}

// Like 点赞文章
// @Summary 点赞文章
// @Description 为文章点赞，需要认证
// @Tags 文章
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} common.Response "点赞成功"
// @Failure 400 {object} common.Response "无效的文章ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/{id}/like [post]
func (s *ArticleService) Like(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的文章ID"))
		return
	}

	err = s.biz.Like(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "点赞文章失败", err))
		return
	}

	common.Success(c, nil)
}

// Unlike 取消点赞
// @Summary 取消点赞
// @Description 取消对文章的点赞，需要认证
// @Tags 文章
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "文章ID"
// @Success 200 {object} common.Response "取消点赞成功"
// @Failure 400 {object} common.Response "无效的文章ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/articles/{id}/like [delete]
func (s *ArticleService) Unlike(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的文章ID"))
		return
	}

	err = s.biz.Unlike(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "取消点赞失败", err))
		return
	}

	common.Success(c, nil)
}
