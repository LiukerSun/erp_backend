package oss

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSTSTokenHandler 获取STS临时凭证处理器 (用于前端直传)
// @Summary 获取STS临时凭证
// @Description 为前端直传获取阿里云OSS STS临时访问凭证
// @Tags OSS
// @Produce json
// @Success 200 {object} STSResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /oss/sts/token [get]
func GetSTSTokenHandler(c *gin.Context) {
	// 获取STS临时凭证
	credentials, err := GetSTSCredentials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取STS凭证失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取STS凭证成功",
		"data":    credentials,
	})
}
