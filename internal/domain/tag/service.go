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
// @Summary 创建标签
// @Description 创建新标签，需要管理员权限
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateRequest true "创建标签请求参数"
// @Success 200 {object} common.Response{data=Response} "创建成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/tags [post]
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
// @Summary 更新标签
// @Description 更新标签信息，需要管理员权限
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateRequest true "更新标签请求参数"
// @Success 200 {object} common.Response{data=Response} "更新成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/tags [put]
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
// @Summary 删除标签
// @Description 删除指定标签，需要管理员权限
// @Tags 标签
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "标签ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "无效的标签ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/tags/{id} [delete]
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
// @Summary 获取标签列表
// @Description 分页获取标签列表
// @Tags 标签
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/tags [get]
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
