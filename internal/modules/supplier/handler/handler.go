package handler

import (
	"erp/internal/modules/supplier/model"
	"erp/internal/modules/supplier/service"
	"erp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 供应商处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建供应商处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateSupplier 创建供应商
// @Summary 创建供应商
// @Description 创建新的供应商
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Param supplier body model.CreateSupplierRequest true "供应商信息"
// @Success 201 {object} response.Response{data=model.SupplierResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers [post]
func (h *Handler) CreateSupplier(c *gin.Context) {
	var req model.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	supplier, err := h.service.CreateSupplier(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("创建供应商失败"))
		return
	}

	c.JSON(http.StatusCreated, response.Success("供应商创建成功", supplier))
}

// GetSupplier 获取供应商详情
// @Summary 获取供应商详情
// @Description 根据ID获取供应商详细信息，可选择是否包含关联的店铺
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Param id path int true "供应商ID"
// @Param include_stores query boolean false "是否包含关联店铺" default(true)
// @Success 200 {object} response.Response{data=model.SupplierResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /suppliers/{id} [get]
func (h *Handler) GetSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的供应商ID"))
		return
	}

	// 检查是否包含店铺信息，默认为true
	includeStores := c.DefaultQuery("include_stores", "true") == "true"

	supplier, err := h.service.GetSupplier(uint(id), includeStores)
	if err != nil {
		if err.Error() == "供应商不存在" {
			c.JSON(http.StatusNotFound, response.Error("供应商不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("获取供应商失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取供应商成功", supplier))
}

// UpdateSupplier 更新供应商
// @Summary 更新供应商
// @Description 更新供应商信息
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Param id path int true "供应商ID"
// @Param supplier body model.UpdateSupplierRequest true "供应商更新信息"
// @Success 200 {object} response.Response{data=model.SupplierResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /suppliers/{id} [put]
func (h *Handler) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的供应商ID"))
		return
	}

	var req model.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	supplier, err := h.service.UpdateSupplier(uint(id), &req)
	if err != nil {
		if err.Error() == "供应商不存在" {
			c.JSON(http.StatusNotFound, response.Error("供应商不存在"))
		} else {
			c.JSON(http.StatusBadRequest, response.Error("更新供应商失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("供应商更新成功", supplier))
}

// DeleteSupplier 删除供应商
// @Summary 删除供应商
// @Description 删除指定供应商（软删除）
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Param id path int true "供应商ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /suppliers/{id} [delete]
func (h *Handler) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的供应商ID"))
		return
	}

	err = h.service.DeleteSupplier(uint(id))
	if err != nil {
		if err.Error() == "供应商不存在" {
			c.JSON(http.StatusNotFound, response.Error("供应商不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("删除供应商失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("供应商删除成功", nil))
}

// ListSuppliers 获取供应商列表
// @Summary 获取供应商列表
// @Description 获取供应商列表，支持分页、搜索、排序、筛选
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param search query string false "搜索关键词（名称、备注）"
// @Param include_stores query boolean false "是否包含关联店铺" default(true)
// @Param is_active query boolean false "活跃状态筛选"
// @Param order_by query string false "排序字段" default("created_at DESC")
// @Success 200 {object} response.Response{data=model.SupplierListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /suppliers [get]
func (h *Handler) ListSuppliers(c *gin.Context) {
	var req model.QuerySuppliersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	result, err := h.service.ListSuppliers(req.Page, req.Limit, req.Search, req.IncludeStores, req.IsActive, req.OrderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询供应商失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("查询供应商成功", result))
}
