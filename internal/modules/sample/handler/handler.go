package handler

import (
	"erp/internal/modules/sample/model"
	"erp/internal/modules/sample/service"
	"erp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 样品处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建样品处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateSample 创建样品
// @Summary 创建样品
// @Description 创建新的样品
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param sample body model.CreateSampleRequest true "样品信息"
// @Success 201 {object} response.Response{data=model.SampleResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /samples [post]
func (h *Handler) CreateSample(c *gin.Context) {
	var req model.CreateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	sample, err := h.service.CreateSample(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("样品创建成功", sample))
}

// GetSample 获取样品详情
// @Summary 获取样品详情
// @Description 根据ID获取样品详细信息
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param id path int true "样品ID"
// @Success 200 {object} response.Response{data=model.SampleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /samples/{id} [get]
func (h *Handler) GetSample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的样品ID"))
		return
	}

	sample, err := h.service.GetSample(uint(id))
	if err != nil {
		if err.Error() == "样品不存在" {
			c.JSON(http.StatusNotFound, response.Error("样品不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("获取样品失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取样品成功", sample))
}

// UpdateSample 更新样品
// @Summary 更新样品
// @Description 更新样品信息
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param id path int true "样品ID"
// @Param sample body model.UpdateSampleRequest true "样品更新信息"
// @Success 200 {object} response.Response{data=model.SampleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /samples/{id} [put]
func (h *Handler) UpdateSample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的样品ID"))
		return
	}

	var req model.UpdateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	sample, err := h.service.UpdateSample(uint(id), &req)
	if err != nil {
		if err.Error() == "样品不存在" {
			c.JSON(http.StatusNotFound, response.Error("样品不存在"))
		} else {
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("样品更新成功", sample))
}

// DeleteSample 删除样品
// @Summary 删除样品
// @Description 删除指定样品（软删除）
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param id path int true "样品ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /samples/{id} [delete]
func (h *Handler) DeleteSample(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的样品ID"))
		return
	}

	err = h.service.DeleteSample(uint(id))
	if err != nil {
		if err.Error() == "样品不存在" {
			c.JSON(http.StatusNotFound, response.Error("样品不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("删除样品失败"))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("样品删除成功", nil))
}

// ListSamples 获取样品列表
// @Summary 获取样品列表
// @Description 获取样品列表，支持分页、搜索、排序、筛选
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param search query string false "搜索关键词（货号）"
// @Param supplier_id query int false "供应商ID筛选"
// @Param store_id query int false "店铺ID筛选"
// @Param has_link query boolean false "是否制作链接筛选"
// @Param is_offline query boolean false "是否下架筛选"
// @Param can_modify_stock query boolean false "是否可修改库存筛选"
// @Param order_by query string false "排序字段" default("created_at DESC")
// @Success 200 {object} response.Response{data=model.SampleListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /samples [get]
func (h *Handler) ListSamples(c *gin.Context) {
	var req model.QuerySamplesRequest
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

	result, err := h.service.ListSamples(req.Page, req.Limit, req.Search, req.SupplierID, req.StoreID, req.HasLink, req.IsOffline, req.CanModifyStock, req.OrderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询样品失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("查询样品成功", result))
}

// BatchUpdateSamples 批量更新样品状态
// @Summary 批量更新样品状态
// @Description 批量更新样品的制作链接、下架状态、库存修改权限
// @Tags 样品管理
// @Accept json
// @Produce json
// @Param request body model.BatchUpdateSamplesRequest true "批量更新请求"
// @Success 200 {object} response.Response{data=model.BatchUpdateSamplesResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /samples/batch-update [patch]
func (h *Handler) BatchUpdateSamples(c *gin.Context) {
	var req model.BatchUpdateSamplesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	// 验证样品ID列表
	if len(req.SampleIDs) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("样品ID列表不能为空"))
		return
	}

	// 限制批量操作的数量
	if len(req.SampleIDs) > 100 {
		c.JSON(http.StatusBadRequest, response.Error("批量操作数量不能超过100个"))
		return
	}

	// 验证是否至少有一个要更新的字段
	if req.HasLink == nil && req.IsOffline == nil && req.CanModifyStock == nil {
		c.JSON(http.StatusBadRequest, response.Error("至少需要指定一个要更新的字段"))
		return
	}

	result, err := h.service.BatchUpdateSamples(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("批量更新样品状态成功", result))
}
