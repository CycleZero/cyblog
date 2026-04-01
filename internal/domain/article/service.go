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
