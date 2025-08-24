package handler

import (
	"net/http"

	"erp/internal/modules/excel/service"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler Excel处理器
type Handler struct {
	service *service.Service
}

// NewHandler 创建Excel处理器
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// ParseExcel 解析Excel文件
// @Summary 解析Excel文件
// @Description 上传Excel文件并解析为JSON格式返回
// @Tags Excel
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Param sheet_name formData string false "Sheet名称（可选）"
// @Success 200 {object} response.Response{data=model.ExcelParseResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /excel/parse [post]
func (h *Handler) ParseExcel(c *gin.Context) {
	// TODO 添加一个店铺参数，解析的每一条数据都加上店铺
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("文件上传失败: "+err.Error()))
		return
	}

	// 获取Sheet名称（可选）
	sheetName := c.PostForm("sheet_name")

	// 调用服务层解析文件
	result, err := h.service.ParseExcel(c.Request.Context(), file, sheetName)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.Success("Excel解析成功", result))
}
