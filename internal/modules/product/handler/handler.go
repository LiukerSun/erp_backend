package handler

import (
	"errors"
	"net/http"
	"strconv"

	"erp/internal/modules/product/model"
	"erp/internal/modules/product/repository"
	"erp/internal/modules/product/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name          string              `json:"name" binding:"required" example:"iPhone 14"`
	SKU           string              `json:"sku" binding:"required" example:"IPHONE14-128G-BLACK"`
	SourceID      *uint               `json:"source_id" example:"1"` // 货源ID
	Price         float64             `json:"price" binding:"required" example:"6999.00"`
	IsDiscounted  bool                `json:"is_discounted" example:"true"`
	DiscountPrice float64             `json:"discount_price" example:"6799.00"`
	CostPrice     float64             `json:"cost_price" binding:"required" example:"5999.00"`
	IsEnabled     bool                `json:"is_enabled" example:"true"`
	Images        model.ProductImages `json:"images"` // 商品图片列表
	Colors        []string            `json:"colors" example:"['黑色','白色','蓝色']"`
	Tags          []uint              `json:"tags" example:"[1,2,3]"` // 标签ID列表
	ShippingTime  string              `json:"shipping_time" example:"三天"`
}

// UpdateProductRequest 更新商品请求（字段都是可选的）
type UpdateProductRequest struct {
	Name          *string              `json:"name,omitempty" example:"iPhone 14"`
	SKU           *string              `json:"sku,omitempty" example:"IPHONE14-128G-BLACK"`
	SourceID      *uint                `json:"source_id,omitempty" example:"1"` // 货源ID
	Price         *float64             `json:"price,omitempty" example:"6999.00"`
	IsDiscounted  *bool                `json:"is_discounted,omitempty" example:"true"`
	DiscountPrice *float64             `json:"discount_price,omitempty" example:"6799.00"`
	CostPrice     *float64             `json:"cost_price,omitempty" example:"5999.00"`
	IsEnabled     *bool                `json:"is_enabled,omitempty" example:"true"`
	Images        *model.ProductImages `json:"images,omitempty"` // 商品图片列表
	Colors        *[]string            `json:"colors,omitempty" example:"['黑色','白色','蓝色']"`
	Tags          *[]uint              `json:"tags,omitempty" example:"[1,2,3]"` // 标签ID列表
	ShippingTime  *string              `json:"shipping_time,omitempty" example:"三天"`
}

// @Summary 创建商品
// @Description 创建新的商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body CreateProductRequest true "商品信息"
// @Success 200 {object} response.Response{data=model.Product} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /product [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	product := &model.Product{
		Name:          req.Name,
		SKU:           req.SKU,
		SourceID:      req.SourceID,
		Price:         req.Price,
		IsDiscounted:  req.IsDiscounted,
		DiscountPrice: req.DiscountPrice,
		CostPrice:     req.CostPrice,
		IsEnabled:     req.IsEnabled,
		Images:        req.Images,
		ShippingTime:  req.ShippingTime,
	}

	if err := h.svc.CreateProduct(c.Request.Context(), product, req.Colors, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("创建商品成功", product))
}

// @Summary 更新商品
// @Description 更新指定ID的商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Param product body CreateProductRequest true "商品信息"
// @Success 200 {object} response.Response{data=model.Product} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Router /product/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的商品ID"))
		return
	}

	// 先检查商品是否存在
	_, err = h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 更新商品信息，只更新提供的字段
	product := &model.Product{
		ID: uint(id),
	}

	// 只更新提供的字段
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.SourceID != nil {
		product.SourceID = req.SourceID
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.IsDiscounted != nil {
		product.IsDiscounted = *req.IsDiscounted
	}
	if req.DiscountPrice != nil {
		product.DiscountPrice = *req.DiscountPrice
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}
	if req.IsEnabled != nil {
		product.IsEnabled = *req.IsEnabled
	}
	if req.Images != nil {
		product.Images = *req.Images
	}
	if req.ShippingTime != nil {
		product.ShippingTime = *req.ShippingTime
	}

	var colors []string
	if req.Colors != nil {
		colors = *req.Colors
	}

	var tags []uint
	if req.Tags != nil {
		tags = *req.Tags
	}

	if err := h.svc.UpdateProduct(c.Request.Context(), product, colors, tags); err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("更新商品成功", product))
}

// @Summary 删除商品
// @Description 删除指定ID的商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Router /product/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的商品ID"))
		return
	}

	if err := h.svc.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("删除商品成功", "商品已删除"))
}

// @Summary 获取商品详情
// @Description 获取指定ID的商品详细信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=model.Product} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Router /product/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的商品ID"))
		return
	}

	product, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取商品详情成功", product))
}

// @Summary 获取商品列表
// @Description 分页获取商品列表，支持多种筛选条件
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认1" default(1)
// @Param page_size query int false "每页数量，默认10" default(10)
// @Param name query string false "商品名称搜索（模糊匹配）"
// @Param sku query string false "SKU精确搜索"
// @Param product_code query string false "商品编码搜索（模糊匹配）"
// @Param source_id query int false "货源ID筛选"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param is_discounted query boolean false "是否优惠筛选"
// @Param is_enabled query boolean false "是否启用筛选"
// @Param colors query []string false "颜色筛选（可多选）"
// @Param shipping_time query string false "发货时间筛选（模糊匹配）"
// @Param order_by query string false "排序字段: id, name, sku, product_code, price, discount_price, cost_price, is_discounted, is_enabled, shipping_time, created_at, updated_at" Enums(id, name, sku, product_code, price, discount_price, cost_price, is_discounted, is_enabled, shipping_time, created_at, updated_at)
// @Param order_dir query string false "排序方向: asc, desc" Enums(asc, desc)
// @Success 200 {object} response.Response{data=object{items=[]model.Product,total=int64,page=int,page_size=int,total_pages=int}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /product [get]
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 限制每页最大数量
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	// 构建筛选条件，支持排序
	var filter repository.ProductListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("筛选参数错误: "+err.Error()))
		return
	}

	// 调试日志
	println("DEBUG: Received order_by =", filter.OrderBy, "order_dir =", filter.OrderDir)
	println("DEBUG: Raw query params - order_by =", c.Query("order_by"), "order_dir =", c.Query("order_dir"))
	println("DEBUG: All query params:")
	for key, values := range c.Request.URL.Query() {
		println("  ", key, "=", values)
	}

	// 统一使用高级筛选方法，支持排序
	products, total, err := h.svc.ListProductsWithFilter(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, response.Success("获取商品列表成功", gin.H{
		"items":       products,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"filter":      filter, // 返回使用的筛选条件
	}))
}

// CreateColorRequest 创建颜色请求
type CreateColorRequest struct {
	Name     string `json:"name" binding:"required" example:"黑色"`
	Code     string `json:"code" example:"BLACK"`
	HexColor string `json:"hex_color" example:"#000000"`
}

// @Summary 创建颜色
// @Description 创建新的颜色
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param color body CreateColorRequest true "颜色信息"
// @Success 200 {object} response.Response{data=model.Color} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /product/colors [post]
func (h *ProductHandler) CreateColor(c *gin.Context) {
	var req CreateColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	color, err := h.svc.CreateColor(c.Request.Context(), req.Name, req.Code, req.HexColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("创建颜色成功", color))
}

// @Summary 获取颜色列表
// @Description 获取所有颜色列表，支持排序
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_by query string false "排序字段: id, name, code, hex_color, created_at, updated_at" Enums(id, name, code, hex_color, created_at, updated_at)
// @Param order_dir query string false "排序方向: asc, desc" Enums(asc, desc)
// @Success 200 {object} response.Response{data=[]model.Color} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Router /product/colors [get]
func (h *ProductHandler) ListColors(c *gin.Context) {
	orderBy := c.DefaultQuery("order_by", "id")
	orderDir := c.DefaultQuery("order_dir", "asc")

	colors, err := h.svc.ListColors(c.Request.Context(), orderBy, orderDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取颜色列表成功", colors))
}

// @Summary 获取颜色详情
// @Description 获取指定ID的颜色详细信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "颜色ID"
// @Success 200 {object} response.Response{data=model.Color} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "颜色不存在"
// @Router /product/colors/{id} [get]
func (h *ProductHandler) GetColor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的颜色ID"))
		return
	}

	color, err := h.svc.GetColor(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, response.Error("颜色不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取颜色详情成功", color))
}

// @Summary 更新颜色
// @Description 更新指定ID的颜色信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "颜色ID"
// @Param color body CreateColorRequest true "颜色信息"
// @Success 200 {object} response.Response{data=model.Color} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "颜色不存在"
// @Router /product/colors/{id} [put]
func (h *ProductHandler) UpdateColor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的颜色ID"))
		return
	}

	var req CreateColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	color, err := h.svc.UpdateColor(c.Request.Context(), uint(id), req.Name, req.Code, req.HexColor)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, response.Error("颜色不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("更新颜色成功", color))
}

// @Summary 删除颜色
// @Description 删除指定ID的颜色
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "颜色ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "颜色不存在"
// @Router /product/colors/{id} [delete]
func (h *ProductHandler) DeleteColor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的颜色ID"))
		return
	}

	if err := h.svc.DeleteColor(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, response.Error("颜色不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("删除颜色成功", "颜色已删除"))
}

// 注意：OSS相关的接口已移除
// - GetOSSCredentials: 使用 GET /oss/sts/token 统一获取STS凭证
// - ValidateFile: 前端直传模式下不需要后端验证文件

// UpdateImageOrderRequest 更新图片顺序请求
type UpdateImageOrderRequest struct {
	Images model.ProductImages `json:"images" binding:"required"`
}

// @Summary 更新商品图片顺序
// @Description 更新商品的图片排序和主图设置
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Param request body UpdateImageOrderRequest true "图片排序信息"
// @Success 200 {object} response.Response{data=model.Product} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Router /product/{id}/images/order [put]
func (h *ProductHandler) UpdateImageOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的商品ID"))
		return
	}

	var req UpdateImageOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 获取商品信息
	product, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	// 验证只能有一张主图
	mainImageCount := 0
	for _, img := range req.Images {
		if img.IsMain {
			mainImageCount++
		}
	}
	if mainImageCount > 1 {
		c.JSON(http.StatusBadRequest, response.Error("只能设置一张主图"))
		return
	}

	// 更新图片信息
	product.Images = req.Images

	// 更新商品
	if err := h.svc.UpdateProduct(c.Request.Context(), product, []string{}, []uint{}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("图片顺序更新成功", product))
}

// SetMainImageRequest 设置主图请求
type SetMainImageRequest struct {
	ImageURL string `json:"image_url" binding:"required" example:"https://example.com/image1.jpg"`
}

// @Summary 设置商品主图
// @Description 设置商品的主图（会取消其他图片的主图状态）
// @Tags 商品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Param request body SetMainImageRequest true "主图URL"
// @Success 200 {object} response.Response{data=model.Product} "设置成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Router /product/{id}/images/main [put]
func (h *ProductHandler) SetMainImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的商品ID"))
		return
	}

	var req SetMainImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 获取商品信息
	product, err := h.svc.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "商品不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		}
		return
	}

	// 验证图片URL是否存在
	found := false
	for _, img := range product.Images {
		if img.URL == req.ImageURL {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusBadRequest, response.Error("指定的图片不存在"))
		return
	}

	// 设置主图
	product.SetMainImage(req.ImageURL)

	// 更新商品
	if err := h.svc.UpdateProduct(c.Request.Context(), product, []string{}, []uint{}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("主图设置成功", product))
}

// GetByCode 通过SKU获取商品
// @Summary 通过SKU获取商品
// @Description 通过SKU获取商品详情，并记录查询历史
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param code path string true "商品SKU"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.Product} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Failure 500 {object} response.Response "系统错误"
// @Router /product/code/{code} [get]
func (h *ProductHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.Error("商品SKU不能为空"))
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("未获取到用户信息"))
		return
	}

	product, err := h.svc.GetByCode(c.Request.Context(), code, userID.(uint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.Error("商品不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("获取商品信息失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取商品成功", product))
}

// GetBySKU 通过SKU获取商品
// @Summary 通过SKU获取商品
// @Description 通过SKU获取商品详情
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param sku path string true "商品SKU"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.Product} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "商品不存在"
// @Failure 500 {object} response.Response "系统错误"
// @Router /product/sku/{sku} [get]
func (h *ProductHandler) GetBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		c.JSON(http.StatusBadRequest, response.Error("商品SKU不能为空"))
		return
	}

	product, err := h.svc.GetBySKU(c.Request.Context(), sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.Error("商品不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("获取商品信息失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取商品成功", product))
}
