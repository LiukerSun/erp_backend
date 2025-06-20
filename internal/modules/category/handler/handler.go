package handler

import (
	"erp/internal/modules/category/model"
	"erp/internal/modules/category/service"
	"erp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler 分类处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建分类处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// CreateCategory godoc
// @Summary 创建分类
// @Description 创建新的分类
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body model.CreateCategoryRequest true "分类创建信息"
// @Success 200 {object} response.Response{data=model.CategoryResponse} "创建成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	category, err := h.service.CreateCategory(c, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("分类创建成功", category))
}

// GetCategory godoc
// @Summary 获取分类详情
// @Description 根据ID获取分类详细信息
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param with_path query bool false "是否包含路径信息"
// @Success 200 {object} response.Response{data=model.CategoryResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "分类不存在"
// @Router /category/{id} [get]
func (h *Handler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	withPath := c.Query("with_path") == "true"

	if withPath {
		category, err := h.service.GetCategoryWithPath(c, uint(id))
		if err != nil {
			response.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, response.Success("获取成功", category))
	} else {
		category, err := h.service.GetCategory(c, uint(id))
		if err != nil {
			response.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, response.Success("获取成功", category))
	}
}

// GetCategories godoc
// @Summary 获取分类列表
// @Description 获取分类列表（支持分页和筛选）
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "分类名称（模糊搜索）"
// @Param parent_id query int false "父分类ID"
// @Param is_active query bool false "是否启用"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=model.CategoryListResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Router /category [get]
func (h *Handler) GetCategories(c *gin.Context) {
	var query model.CategoryQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	categories, err := h.service.GetCategories(c, query)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", categories))
}

// GetCategoryTree godoc
// @Summary 获取分类树
// @Description 获取完整的分类树结构
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.CategoryTreeListResponse} "获取成功"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category/tree [get]
func (h *Handler) GetCategoryTree(c *gin.Context) {
	tree, err := h.service.GetCategoryTree(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", tree))
}

// GetRootCategories godoc
// @Summary 获取根分类
// @Description 获取所有根分类（一级分类）
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.CategoryTreeListResponse} "获取成功"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category/root [get]
func (h *Handler) GetRootCategories(c *gin.Context) {
	categories, err := h.service.GetRootCategories(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", categories))
}

// GetChildrenCategories godoc
// @Summary 获取子分类
// @Description 获取指定分类的所有直接子分类
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "父分类ID"
// @Success 200 {object} response.Response{data=model.CategoryTreeListResponse} "获取成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Router /category/{id}/children [get]
func (h *Handler) GetChildrenCategories(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	categories, err := h.service.GetChildrenCategories(c, uint(id))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("获取成功", categories))
}

// UpdateCategory godoc
// @Summary 更新分类
// @Description 更新分类信息
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param category body model.UpdateCategoryRequest true "分类更新信息"
// @Success 200 {object} response.Response{data=model.CategoryResponse} "更新成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "分类不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	category, err := h.service.UpdateCategory(c, uint(id), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("分类更新成功", category))
}

// MoveCategory godoc
// @Summary 移动分类
// @Description 移动分类到新的父分类下
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param move body model.MoveCategoryRequest true "移动信息"
// @Success 200 {object} response.Response{data=model.CategoryResponse} "移动成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "分类不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category/{id}/move [post]
func (h *Handler) MoveCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	var req model.MoveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	category, err := h.service.MoveCategory(c, uint(id), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("分类移动成功", category))
}

// DeleteCategory godoc
// @Summary 删除分类
// @Description 删除指定分类（软删除）
// @Tags Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response{data=string} "删除成功"
// @Failure 400 {object} response.Response{error=string} "请求参数错误"
// @Failure 401 {object} response.Response{error=string} "未授权"
// @Failure 404 {object} response.Response{error=string} "分类不存在"
// @Failure 500 {object} response.Response{error=string} "服务器内部错误"
// @Router /category/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("无效的分类ID"))
		return
	}

	if err := h.service.DeleteCategory(c, uint(id)); err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("分类删除成功", "分类已删除"))
}
