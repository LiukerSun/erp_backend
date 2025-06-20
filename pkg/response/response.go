package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 成功响应
func Success(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}

// HandleError 处理错误响应
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 根据错误类型返回不同的状态码
	switch err.Error() {
	case "用户名已存在", "邮箱已存在", "邮箱已被其他用户使用":
		c.JSON(http.StatusConflict, Error(err.Error()))
	case "用户名或密码错误", "原密码错误", "账户已被禁用":
		c.JSON(http.StatusUnauthorized, Error(err.Error()))
	case "用户不存在":
		c.JSON(http.StatusNotFound, Error(err.Error()))
	case "权限不足":
		c.JSON(http.StatusForbidden, Error(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, Error(err.Error()))
	}
}
