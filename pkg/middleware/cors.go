package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// CORSMiddlewareWithConfig 可配置的CORS中间件
func CORSMiddlewareWithConfig(origins []string, methods []string, headers []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		origin := "*"
		if len(origins) > 0 {
			origin = origins[0]
			for _, o := range origins {
				if c.Request.Header.Get("Origin") == o {
					origin = o
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Origin", origin)

		// 设置允许的方法
		methodsStr := "GET, POST, PUT, DELETE, OPTIONS"
		if len(methods) > 0 {
			methodsStr = ""
			for i, m := range methods {
				if i > 0 {
					methodsStr += ", "
				}
				methodsStr += m
			}
		}
		c.Header("Access-Control-Allow-Methods", methodsStr)

		// 设置允许的头部
		headersStr := "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
		if len(headers) > 0 {
			headersStr = ""
			for i, h := range headers {
				if i > 0 {
					headersStr += ", "
				}
				headersStr += h
			}
		}
		c.Header("Access-Control-Allow-Headers", headersStr)

		// 允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
