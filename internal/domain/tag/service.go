package tag

import (
	"cyblog/internal/common"
	"cyblog/pkg/errs"
	"cyblog/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagService struct {
	common.BaseService
	biz    *TagBiz
	logger *log.Logger
}

func NewTagService(biz *TagBiz, logger *log.Logger) *TagService {
	return &TagService{
		biz:    biz,
		logger: logger,
	}
}

// Create 创建标签
func (s *TagService) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Create(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "创建标签失败", err))
		return
	}

	common.Success(c, resp)
}

// Update 更新标签
func (s *TagService) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusBadRequest, "参数错误", err))
		return
	}

	resp, err := s.biz.Update(c, &req)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "更新标签失败", err))
		return
	}

	common.Success(c, resp)
}

// Delete 删除标签
func (s *TagService) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		common.Error(c, errs.New(http.StatusBadRequest, "无效的标签ID"))
		return
	}

	err = s.biz.Delete(c, uint(id))
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "删除标签失败", err))
		return
	}

	common.Success(c, nil)
}

// List 获取标签列表
func (s *TagService) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	resp, err := s.biz.List(c, page, pageSize)
	if err != nil {
		common.Error(c, errs.WrapWithMsg(http.StatusInternalServerError, "获取标签列表失败", err))
		return
	}

	common.Success(c, resp)
}
