package handler

import (
	"net/http"
	"strconv"

	"erp/internal/modules/attribute/model"
	"erp/internal/modules/attribute/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler 属性处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建属性处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// 属性管理接口

// CreateAttribute 创建属性
// @Summary 创建属性
// @Description 创建一个新的属性
// @Tags Attribute
// @Accept json
// @Produce json
// @Param request body model.CreateAttributeRequest true "创建属性请求"
// @Success 200 {object} response.Response{data=model.AttributeResponse} "创建成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attributes [post]
func (h *Handler) CreateAttribute(c *gin.Context) {
	var req model.CreateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	attr, err := h.service.CreateAttribute(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("创建属性失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("创建属性成功", attr))
}

// GetAttributes 获取属性列表
// @Summary 获取属性列表
// @Description 获取属性列表，支持分页和筛选
// @Tags Attribute
// @Accept json
// @Produce json
// @Param name query string false "属性名称（模糊搜索）"
// @Param type query string false "属性类型"
// @Param is_active query boolean false "是否启用"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=model.AttributeListResponse} "获取成功"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attributes [get]
func (h *Handler) GetAttributes(c *gin.Context) {
	var req model.AttributeQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	resp, err := h.service.GetAttributesList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取属性列表失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取属性列表成功", resp))
}

// GetAttributeTypes 获取属性类型列表
// @Summary 获取属性类型列表
// @Description 获取所有支持的属性类型
// @Tags Attribute
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string} "获取成功"
// @Security BearerAuth
// @Router /attributes/types [get]
func (h *Handler) GetAttributeTypes(c *gin.Context) {
	types := []string{
		string(model.AttributeTypeText),
		string(model.AttributeTypeNumber),
		string(model.AttributeTypeSelect),
		string(model.AttributeTypeMultiSelect),
		string(model.AttributeTypeBoolean),
		string(model.AttributeTypeDate),
		string(model.AttributeTypeDateTime),
		string(model.AttributeTypeURL),
		string(model.AttributeTypeEmail),
		string(model.AttributeTypeColor),
		string(model.AttributeTypeCurrency),
	}

	c.JSON(http.StatusOK, response.Success("获取属性类型成功", types))
}

// GetAttribute 获取属性详情
// @Summary 获取属性详情
// @Description 根据ID获取属性详情
// @Tags Attribute
// @Accept json
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response{data=model.AttributeResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 404 {object} response.Response{error=string} "属性不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attributes/{id} [get]
func (h *Handler) GetAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性ID"))
		return
	}

	attr, err := h.service.GetAttributeByID(uint(id))
	if err != nil {
		if err.Error() == "属性不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("获取属性失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("获取属性成功", attr))
}

// UpdateAttribute 更新属性
// @Summary 更新属性
// @Description 根据ID更新属性信息
// @Tags Attribute
// @Accept json
// @Produce json
// @Param id path int true "属性ID"
// @Param request body model.UpdateAttributeRequest true "更新属性请求"
// @Success 200 {object} response.Response{data=model.AttributeResponse} "更新成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 404 {object} response.Response{error=string} "属性不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attributes/{id} [put]
func (h *Handler) UpdateAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性ID"))
		return
	}

	var req model.UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	attr, err := h.service.UpdateAttribute(uint(id), &req)
	if err != nil {
		if err.Error() == "属性不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("更新属性失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("更新属性成功", attr))
}

// DeleteAttribute 删除属性
// @Summary 删除属性
// @Description 根据ID删除属性（软删除）
// @Tags Attribute
// @Accept json
// @Produce json
// @Param id path int true "属性ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 404 {object} response.Response{error=string} "属性不存在"
// @Failure 409 {object} response.Response{error=string} "属性已被使用，无法删除"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attributes/{id} [delete]
func (h *Handler) DeleteAttribute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性ID"))
		return
	}

	err = h.service.DeleteAttribute(uint(id))
	if err != nil {
		if err.Error() == "属性不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else if err.Error() == "属性已被使用，无法删除" {
			c.JSON(http.StatusConflict, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("删除属性失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("删除属性成功", nil))
}

// 分类属性管理接口

// GetCategoryAttributes 获取分类的属性列表
// @Summary 获取分类的属性列表
// @Description 根据分类ID获取绑定的属性列表
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response{data=model.CategoryAttributesResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes [get]
func (h *Handler) GetCategoryAttributes(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	attrs, err := h.service.GetCategoryAttributes(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取分类属性失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取分类属性成功", attrs))
}

// GetCategoryAttributesWithInheritance 获取分类的属性列表（包括继承）
// @Summary 获取分类的属性列表（包括继承）
// @Description 根据分类ID获取绑定的属性列表，包括从父分类继承的属性
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response{data=model.CategoryAttributesWithInheritanceResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes/inheritance [get]
func (h *Handler) GetCategoryAttributesWithInheritance(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	attrs, err := h.service.GetCategoryAttributesWithInheritance(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取分类继承属性失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取分类继承属性成功", attrs))
}

// GetAttributeInheritancePath 获取属性的继承路径
// @Summary 获取属性的继承路径
// @Description 获取指定属性在分类层级中的继承路径信息
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Param attribute_id path int true "属性ID"
// @Success 200 {object} response.Response{data=model.AttributeInheritancePathResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes/{attribute_id}/inheritance [get]
func (h *Handler) GetAttributeInheritancePath(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	attributeIDStr := c.Param("attribute_id")
	attributeID, err := strconv.ParseUint(attributeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性ID"))
		return
	}

	path, err := h.service.GetAttributeInheritancePath(uint(categoryID), uint(attributeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取属性继承路径失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取属性继承路径成功", path))
}

// BindAttributeToCategory 绑定属性到分类
// @Summary 绑定属性到分类
// @Description 将指定属性绑定到分类
// @Tags Attribute
// @Accept json
// @Produce json
// @Param request body model.BindAttributeToCategoryRequest true "绑定属性到分类请求"
// @Success 200 {object} response.Response{data=model.CategoryAttributeResponse} "绑定成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/attributes/bind [post]
func (h *Handler) BindAttributeToCategory(c *gin.Context) {
	var req model.BindAttributeToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	categoryAttr, err := h.service.BindAttributeToCategory(req.CategoryID, req.AttributeID, req.IsRequired, req.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("绑定属性到分类失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("绑定属性到分类成功", categoryAttr))
}

// UnbindAttributeFromCategory 从分类解绑属性
// @Summary 从分类解绑属性
// @Description 将指定属性从分类中解绑
// @Tags Attribute
// @Accept json
// @Produce json
// @Param request body model.UnbindAttributeFromCategoryRequest true "从分类解绑属性请求"
// @Success 200 {object} response.Response "解绑成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/attributes/unbind [post]
func (h *Handler) UnbindAttributeFromCategory(c *gin.Context) {
	var req model.UnbindAttributeFromCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	err := h.service.UnbindAttributeFromCategory(req.CategoryID, req.AttributeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("从分类解绑属性失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("从分类解绑属性成功", nil))
}

// UpdateCategoryAttribute 更新分类属性关联
// @Summary 更新分类属性关联
// @Description 更新分类与属性的关联信息
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Param attribute_id path int true "属性ID"
// @Param request body model.UpdateCategoryAttributeRequest true "更新分类属性关联请求"
// @Success 200 {object} response.Response{data=model.CategoryAttributeResponse} "更新成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 404 {object} response.Response{error=string} "分类属性关联不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes/{attribute_id} [put]
func (h *Handler) UpdateCategoryAttribute(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	attributeIDStr := c.Param("attribute_id")
	attributeID, err := strconv.ParseUint(attributeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性ID"))
		return
	}

	var req model.UpdateCategoryAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	categoryAttr, err := h.service.UpdateCategoryAttribute(uint(categoryID), uint(attributeID), req.IsRequired, req.Sort)
	if err != nil {
		if err.Error() == "分类属性关联不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("更新分类属性关联失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("更新分类属性关联成功", categoryAttr))
}

// BatchBindAttributesToCategory 批量绑定属性到分类
// @Summary 批量绑定属性到分类
// @Description 将多个属性批量绑定到指定分类
// @Tags Attribute
// @Accept json
// @Produce json
// @Param request body model.BatchBindAttributesToCategoryRequest true "批量绑定属性到分类请求"
// @Success 200 {object} response.Response "绑定成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/attributes/batch-bind [post]
func (h *Handler) BatchBindAttributesToCategory(c *gin.Context) {
	var req model.BatchBindAttributesToCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	err := h.service.BatchBindAttributesToCategory(req.CategoryID, req.Attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("批量绑定属性到分类失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("批量绑定属性到分类成功", nil))
}

// RebuildCategoryInheritance 重建分类属性继承关系
// @Summary 重建分类属性继承关系
// @Description 重建指定分类的属性继承关系，用于修复不一致的情况
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response "重建成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes/rebuild-inheritance [post]
func (h *Handler) RebuildCategoryInheritance(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	err = h.service.RebuildCategoryInheritance(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("重建分类继承关系失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("重建分类继承关系成功", nil))
}

// ValidateInheritanceConsistency 验证继承关系一致性
// @Summary 验证继承关系一致性
// @Description 验证指定分类的属性继承关系是否一致
// @Tags Attribute
// @Accept json
// @Produce json
// @Param category_id path int true "分类ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "验证成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /categories/{category_id}/attributes/validate-inheritance [get]
func (h *Handler) ValidateInheritanceConsistency(c *gin.Context) {
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	isConsistent, issues, err := h.service.ValidateInheritanceConsistency(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("验证继承关系一致性失败: "+err.Error()))
		return
	}

	result := map[string]interface{}{
		"is_consistent": isConsistent,
		"issues":        issues,
		"category_id":   categoryID,
	}

	c.JSON(http.StatusOK, response.Success("验证继承关系一致性成功", result))
}

// 属性值管理接口

// SetAttributeValue 设置属性值
// @Summary 设置属性值
// @Description 为实体设置属性值
// @Tags AttributeValue
// @Accept json
// @Produce json
// @Param request body model.SetAttributeValueRequest true "设置属性值请求"
// @Success 200 {object} response.Response{data=model.AttributeValueResponse} "设置成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attribute-values [post]
func (h *Handler) SetAttributeValue(c *gin.Context) {
	var req model.SetAttributeValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	value, err := h.service.SetAttributeValue(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("设置属性值失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("设置属性值成功", value))
}

// GetAttributeValues 获取实体属性值
// @Summary 获取实体属性值
// @Description 根据实体类型和ID获取属性值列表
// @Tags AttributeValue
// @Accept json
// @Produce json
// @Param entity_type query string true "实体类型"
// @Param entity_id query int true "实体ID"
// @Success 200 {object} response.Response{data=model.EntityAttributeValuesResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attribute-values [get]
func (h *Handler) GetAttributeValues(c *gin.Context) {
	var req model.GetAttributeValuesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	values, err := h.service.GetAttributeValuesByEntity(req.EntityType, req.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取属性值失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("获取属性值成功", values))
}

// DeleteAttributeValue 删除属性值
// @Summary 删除属性值
// @Description 根据ID删除属性值
// @Tags AttributeValue
// @Accept json
// @Produce json
// @Param id path int true "属性值ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 404 {object} response.Response{error=string} "属性值不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attribute-values/{id} [delete]
func (h *Handler) DeleteAttributeValue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的属性值ID"))
		return
	}

	err = h.service.DeleteAttributeValue(uint(id))
	if err != nil {
		if err.Error() == "属性值不存在" {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError, response.Error("删除属性值失败: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("删除属性值成功", nil))
}

// BatchSetAttributeValues 批量设置属性值
// @Summary 批量设置属性值
// @Description 批量为实体设置属性值
// @Tags AttributeValue
// @Accept json
// @Produce json
// @Param request body []model.SetAttributeValueRequest true "批量设置属性值请求"
// @Success 200 {object} response.Response "设置成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Security BearerAuth
// @Router /attribute-values/batch [post]
func (h *Handler) BatchSetAttributeValues(c *gin.Context) {
	var req []model.SetAttributeValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("属性值列表不能为空"))
		return
	}

	err := h.service.BatchSetAttributeValues(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("批量设置属性值失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("批量设置属性值成功", nil))
}
