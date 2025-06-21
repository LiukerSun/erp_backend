package handler

import (
	"net/http"
	"strconv"

	"erp/internal/modules/tags/model"
	"erp/internal/modules/tags/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required" example:"热销"`
	Description string `json:"description" example:"热销商品标签"`
	Color       string `json:"color" binding:"required" example:"#FF6B6B"`
	IsEnabled   bool   `json:"is_enabled" example:"true"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        *string `json:"name,omitempty" example:"热销"`
	Description *string `json:"description,omitempty" example:"热销商品标签"`
	Color       *string `json:"color,omitempty" example:"#FF6B6B"`
	IsEnabled   *bool   `json:"is_enabled,omitempty" example:"true"`
}

type TagsHandler struct {
	service *service.TagsService
}

func NewTagsHandler(service *service.TagsService) *TagsHandler {
	return &TagsHandler{service: service}
}

// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建新的标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param tag body CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=internal_modules_tags_model.Tag}
// @Router /api/tags [post]
func (h *TagsHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	tag := &model.Tag{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		IsEnabled:   req.IsEnabled,
	}

	if err := h.service.CreateTag(tag); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("创建标签失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("标签创建成功", tag))
}

// GetTagByID 根据ID获取标签
// @Summary 获取标签详情
// @Description 根据ID获取标签详细信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response{data=internal_modules_tags_model.Tag}
// @Router /api/tags/{id} [get]
func (h *TagsHandler) GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	tag, err := h.service.GetTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("标签不存在"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取标签成功", tag))
}

// GetAllTags 获取所有标签
// @Summary 获取所有标签
// @Description 获取所有标签列表
// @Tags 标签管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]internal_modules_tags_model.Tag}
// @Router /api/tags [get]
func (h *TagsHandler) GetAllTags(c *gin.Context) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取标签列表失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取标签列表成功", tags))
}

// GetEnabledTags 获取启用的标签
// @Summary 获取启用的标签
// @Description 获取所有启用的标签列表
// @Tags 标签管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]internal_modules_tags_model.Tag}
// @Router /api/tags/enabled [get]
func (h *TagsHandler) GetEnabledTags(c *gin.Context) {
	tags, err := h.service.GetEnabledTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取启用标签列表失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取启用标签列表成功", tags))
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param tag body UpdateTagRequest true "标签信息"
// @Success 200 {object} response.Response{data=internal_modules_tags_model.Tag}
// @Router /api/tags/{id} [put]
func (h *TagsHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 获取现有标签
	existingTag, err := h.service.GetTagByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("标签不存在"))
		return
	}

	// 只更新提供的字段
	if req.Name != nil {
		existingTag.Name = *req.Name
	}
	if req.Description != nil {
		existingTag.Description = *req.Description
	}
	if req.Color != nil {
		existingTag.Color = *req.Color
	}
	if req.IsEnabled != nil {
		existingTag.IsEnabled = *req.IsEnabled
	}

	if err := h.service.UpdateTag(existingTag); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新标签失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("标签更新成功", existingTag))
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除指定标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response
// @Router /api/tags/{id} [delete]
func (h *TagsHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	if err := h.service.DeleteTag(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除标签失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("标签删除成功", nil))
}

// GetProductsByTag 获取标签下的所有产品
// @Summary 获取标签下的产品
// @Description 获取指定标签下的所有产品
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} response.Response{data=[]model.Product}
// @Router /api/tags/{id}/products [get]
func (h *TagsHandler) GetProductsByTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	products, err := h.service.GetProductsByTag(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取标签产品失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取标签产品成功", products))
}

// AddProductToTag 为标签添加产品
// @Summary 为标签添加产品
// @Description 为指定标签添加产品
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param product_id query int true "产品ID"
// @Success 200 {object} response.Response
// @Router /api/tags/{id}/products [post]
func (h *TagsHandler) AddProductToTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	if err := h.service.AddProductToTag(uint(id), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("添加产品到标签失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("产品添加成功", nil))
}

// RemoveProductFromTag 从标签移除产品
// @Summary 从标签移除产品
// @Description 从指定标签移除产品
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param product_id query int true "产品ID"
// @Success 200 {object} response.Response
// @Router /api/tags/{id}/products [delete]
func (h *TagsHandler) RemoveProductFromTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的标签ID"))
		return
	}

	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	if err := h.service.RemoveProductFromTag(uint(id), uint(productID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("从标签移除产品失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("产品移除成功", nil))
}

// GetTagsByProduct 获取产品的所有标签
// @Summary 获取产品的标签
// @Description 获取指定产品的所有标签
// @Tags 标签管理
// @Accept json
// @Produce json
// @Param product_id query int true "产品ID"
// @Success 200 {object} response.Response{data=[]internal_modules_tags_model.Tag}
// @Router /api/tags/product [get]
func (h *TagsHandler) GetTagsByProduct(c *gin.Context) {
	productIDStr := c.Query("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的产品ID"))
		return
	}

	tags, err := h.service.GetTagsByProduct(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取产品标签失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取产品标签成功", tags))
}
