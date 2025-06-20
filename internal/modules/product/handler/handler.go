package handler

import (
	"net/http"
	"strconv"

	"erp/internal/modules/product/model"
	"erp/internal/modules/product/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler 产品处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建产品处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateProduct godoc
// @Summary 创建产品
// @Description 创建新产品
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body model.CreateProductRequest true "产品创建信息"
// @Success 200 {object} response.Response{data=model.ProductResponse} "创建成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /product [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	product, err := h.service.CreateProduct(c, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("产品创建成功", product))
}

// GetProduct godoc
// @Summary 获取产品详情
// @Description 根据ID获取产品详细信息
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=model.ProductResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "产品不存在"
// @Router /product/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	product, err := h.service.GetProduct(c, uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", product))
}

// UpdateProduct godoc
// @Summary 更新产品
// @Description 更新产品信息
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param product body model.UpdateProductRequest true "产品更新信息"
// @Success 200 {object} response.Response{data=model.ProductResponse} "更新成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "产品不存在"
// @Router /product/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	product, err := h.service.UpdateProduct(c, uint(id), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("产品更新成功", product))
}

// DeleteProduct godoc
// @Summary 删除产品
// @Description 删除产品（软删除）
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} response.Response{data=string} "删除成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "产品不存在"
// @Router /product/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	err = h.service.DeleteProduct(c, uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("产品删除成功", "产品已删除"))
}

// GetProducts godoc
// @Summary 获取产品列表
// @Description 获取产品列表（支持分页和筛选）
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "产品名称（模糊搜索）"
// @Param category_id query int false "产品分类ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=model.ProductListResponse} "获取成功"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Router /product [get]
func (h *Handler) GetProducts(c *gin.Context) {
	var req model.ProductQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	products, err := h.service.SearchProducts(c, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", products))
}

// GetProductsByCategory godoc
// @Summary 根据分类获取产品
// @Description 获取指定分类下的所有产品
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response{data=[]model.ProductResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Router /product/category/{category_id} [get]
func (h *Handler) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	products, err := h.service.GetProductsByCategory(c, uint(categoryID))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", products))
}
