package middleware

import (
	"net/http"
	"strings"

	"erp/internal/modules/user/repository"
	"erp/pkg/auth"
	"erp/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件（基础版本，不验证密码版本）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error("Authorization header is required"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid or expired token"))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AuthMiddlewareWithPasswordValidation JWT认证中间件（包含密码版本验证）
func AuthMiddlewareWithPasswordValidation(userRepo *repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error("Authorization header is required"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid or expired token"))
			c.Abort()
			return
		}

		// 验证密码版本
		currentPasswordVersion, err := userRepo.GetPasswordVersion(c, claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("用户不存在"))
			c.Abort()
			return
		}

		if !auth.ValidateTokenPasswordVersion(claims, currentPasswordVersion) {
			c.JSON(http.StatusUnauthorized, response.Error("Token已失效，请重新登录"))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, response.Error("User role not found"))
			c.Abort()
			return
		}

		role := userRole.(string)
		hasRole := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, response.Error("权限不足"))
			c.Abort()
			return
		}

		c.Next()
	}
}
