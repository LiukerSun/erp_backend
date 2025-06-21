package handler

import (
	"net/http"
	"strconv"

	"erp/internal/modules/source/model"
	"erp/internal/modules/source/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	svc service.SourceService
}

func NewSourceHandler(svc service.SourceService) *SourceHandler {
	return &SourceHandler{svc: svc}
}

// CreateSourceRequest 创建货源请求
type CreateSourceRequest struct {
	Name   string `json:"name" binding:"required" example:"Apple官方旗舰店"`
	Code   string `json:"code" binding:"required" example:"APPLE001"`
	Status int    `json:"status" example:"1"`
	Remark string `json:"remark" example:"优质货源"`
}

// @Summary 创建货源
// @Description 创建新的货源信息
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param source body CreateSourceRequest true "货源信息"
// @Success 200 {object} response.Response{data=internal_modules_source_model.Source} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /source [post]
func (h *SourceHandler) Create(c *gin.Context) {
	var req CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	source := &model.Source{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Remark: req.Remark,
	}

	if err := h.svc.CreateSource(c.Request.Context(), source); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("创建货源成功", source))
}

// @Summary 更新货源
// @Description 更新指定ID的货源信息
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "货源ID"
// @Param source body CreateSourceRequest true "货源信息"
// @Success 200 {object} response.Response{data=internal_modules_source_model.Source} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "货源不存在"
// @Router /source/{id} [put]
func (h *SourceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的货源ID"))
		return
	}

	var req CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	source := &model.Source{
		ID:     uint(id),
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Remark: req.Remark,
	}

	if err := h.svc.UpdateSource(c.Request.Context(), source); err != nil {
		if err.Error() == "货源不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("更新货源成功", source))
}

// @Summary 删除货源
// @Description 删除指定ID的货源
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "货源ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "货源不存在"
// @Router /source/{id} [delete]
func (h *SourceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的货源ID"))
		return
	}

	if err := h.svc.DeleteSource(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "货源不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("删除货源成功", "货源已删除"))
}

// @Summary 获取货源详情
// @Description 获取指定ID的货源详细信息
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "货源ID"
// @Success 200 {object} response.Response{data=internal_modules_source_model.Source} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "货源不存在"
// @Router /source/{id} [get]
func (h *SourceHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的货源ID"))
		return
	}

	source, err := h.svc.GetSource(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "货源不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取货源详情成功", source))
}

// @Summary 获取货源列表
// @Description 分页获取货源列表
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页数量，默认10" default(10)
// @Success 200 {object} response.Response{data=object{items=[]internal_modules_source_model.Source,total=int64}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /source [get]
func (h *SourceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	sources, total, err := h.svc.ListSources(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	result := gin.H{
		"items": sources,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}

	c.JSON(http.StatusOK, response.Success("获取货源列表成功", result))
}

// @Summary 获取启用状态的货源列表
// @Description 获取所有启用状态的货源，用于下拉选择
// @Tags 货源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]internal_modules_source_model.Source} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /source/active [get]
func (h *SourceHandler) ListActive(c *gin.Context) {
	sources, err := h.svc.ListActiveSource(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取启用货源列表成功", sources))
}
