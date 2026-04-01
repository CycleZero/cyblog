package category

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryService struct {
	common.BaseService
	biz    *CategoryBiz
	logger *log.Logger
}

func NewCategoryService(biz *CategoryBiz, logger *log.Logger) *CategoryService {
	return &CategoryService{
		biz:    biz,
		logger: logger,
	}
}

// Create 创建分类
func (s *CategoryService) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Create(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "创建分类失败", err))
		return
	}

	common.Success(c, resp)
}

// Update 更新分类
func (s *CategoryService) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Update(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新分类失败", err))
		return
	}

	common.Success(c, resp)
}

// Delete 删除分类
func (s *CategoryService) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的分类ID"))
		return
	}

	err = s.biz.Delete(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "删除分类失败", err))
		return
	}

	common.Success(c, nil)
}

// List 获取分类列表
func (s *CategoryService) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := s.biz.List(c, page, pageSize)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取分类列表失败", err))
		return
	}

	common.Success(c, resp)
}
