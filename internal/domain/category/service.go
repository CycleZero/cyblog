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
// @Summary 创建分类
// @Description 创建新分类，需要管理员权限
// @Tags 分类
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateRequest true "创建分类请求参数"
// @Success 200 {object} common.Response{data=Response} "创建成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/categories [post]
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
// @Summary 更新分类
// @Description 更新分类信息，需要管理员权限
// @Tags 分类
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateRequest true "更新分类请求参数"
// @Success 200 {object} common.Response{data=Response} "更新成功"
// @Failure 400 {object} common.Response "请求参数错误"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/categories [put]
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
// @Summary 删除分类
// @Description 删除指定分类，需要管理员权限
// @Tags 分类
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "分类ID"
// @Success 200 {object} common.Response "删除成功"
// @Failure 400 {object} common.Response "无效的分类ID"
// @Failure 401 {object} common.Response "未认证"
// @Failure 403 {object} common.Response "无权限"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/categories/{id} [delete]
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
// @Summary 获取分类列表
// @Description 分页获取分类列表
// @Tags 分类
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.Response{data=ListResponse} "获取成功"
// @Failure 500 {object} common.Response "服务器内部错误"
// @Router /api/categories [get]
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
