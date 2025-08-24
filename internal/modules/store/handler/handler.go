package handler

import (
	"erp/internal/modules/store/model"
	"erp/internal/modules/store/service"
	"erp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 店铺处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建店铺处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateStore 创建店铺
// @Summary 创建店铺
// @Description 创建新的店铺
// @Tags 店铺管理
// @Accept json
// @Produce json
// @Param store body model.CreateStoreRequest true "店铺信息"
// @Success 201 {object} response.Response{data=model.StoreResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /stores [post]
func (h *Handler) CreateStore(c *gin.Context) {
	var req model.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	store, err := h.service.CreateStore(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("创建店铺失败"))
		return
	}

	c.JSON(http.StatusCreated, response.Success("店铺创建成功", store))
}

// GetStore 获取店铺详情
// @Summary 获取店铺详情
// @Description 根据ID获取店铺详细信息
// @Tags 店铺管理
// @Accept json
// @Produce json
// @Param id path int true "店铺ID"
// @Success 200 {object} response.Response{data=model.StoreResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /stores/{id} [get]
func (h *Handler) GetStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的店铺ID"))
		return
	}

	store, err := h.service.GetStore(uint(id))
	if err != nil {
		if err.Error() == "店铺不存在" {
			c.JSON(http.StatusNotFound, response.Error("店铺不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("获取店铺失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取店铺成功", store))
}

// UpdateStore 更新店铺
// @Summary 更新店铺
// @Description 更新店铺信息
// @Tags 店铺管理
// @Accept json
// @Produce json
// @Param id path int true "店铺ID"
// @Param store body model.UpdateStoreRequest true "店铺更新信息"
// @Success 200 {object} response.Response{data=model.StoreResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /stores/{id} [put]
func (h *Handler) UpdateStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的店铺ID"))
		return
	}

	var req model.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	store, err := h.service.UpdateStore(uint(id), &req)
	if err != nil {
		if err.Error() == "店铺不存在" {
			c.JSON(http.StatusNotFound, response.Error("店铺不存在"))
		} else {
			c.JSON(http.StatusBadRequest, response.Error("更新店铺失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("店铺更新成功", store))
}

// DeleteStore 删除店铺
// @Summary 删除店铺
// @Description 删除指定店铺（软删除）
// @Tags 店铺管理
// @Accept json
// @Produce json
// @Param id path int true "店铺ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /stores/{id} [delete]
func (h *Handler) DeleteStore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的店铺ID"))
		return
	}

	err = h.service.DeleteStore(uint(id))
	if err != nil {
		if err.Error() == "店铺不存在" {
			c.JSON(http.StatusNotFound, response.Error("店铺不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("删除店铺失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("店铺删除成功", nil))
}

// ListStores 获取店铺列表
// @Summary 获取店铺列表
// @Description 获取店铺列表，支持分页、搜索、排序、筛选
// @Tags 店铺管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param search query string false "搜索关键词（名称）"
// @Param supplier_id query int false "供应商ID筛选"
// @Param is_active query bool false "活跃状态筛选"
// @Param is_featured query bool false "精选状态筛选"
// @Param order_by query string false "排序字段" default("created_at DESC")
// @Success 200 {object} response.Response{data=model.StoreListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /stores [get]
func (h *Handler) ListStores(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// 获取搜索参数
	search := c.Query("search")

	// 获取供应商ID筛选
	var supplierID *uint
	if supplierIDStr := c.Query("supplier_id"); supplierIDStr != "" {
		if id, err := strconv.ParseUint(supplierIDStr, 10, 32); err == nil {
			supplierIDUint := uint(id)
			supplierID = &supplierIDUint
		}
	}

	// 获取活跃状态筛选
	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	// 获取精选状态筛选
	var isFeatured *bool
	if isFeaturedStr := c.Query("is_featured"); isFeaturedStr != "" {
		if featured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			isFeatured = &featured
		}
	}

	// 获取排序参数
	orderBy := c.DefaultQuery("order_by", "created_at DESC")

	stores, err := h.service.ListStores(page, limit, search, supplierID, isActive, isFeatured, orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取店铺列表失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取店铺列表成功", stores))
}
