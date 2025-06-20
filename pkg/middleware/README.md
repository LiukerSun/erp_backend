# 中间件包 (pkg/middleware)

这个包提供了可复用的Gin中间件组件，可以被其他项目导入使用。

## 可用的中间件

### 1. 认证中间件 (auth.go)

#### AuthMiddleware
JWT认证中间件，用于验证用户身份。

```go
import "erp/pkg/middleware"

// 使用认证中间件
r.Use(middleware.AuthMiddleware())
```

**功能**：
- 验证Authorization头部
- 解析JWT令牌
- 将用户信息存储到上下文中

**上下文变量**：
- `user_id`: 用户ID
- `username`: 用户名
- `role`: 用户角色

#### RoleMiddleware
角色权限中间件，用于检查用户权限。

```go
// 检查单个角色
r.Use(middleware.RoleMiddleware("admin"))

// 检查多个角色
r.Use(middleware.RoleMiddleware("admin", "manager"))
```

**功能**：
- 检查用户角色
- 支持多角色验证
- 权限不足时返回403错误

### 2. CORS中间件 (cors.go)

#### CORSMiddleware
默认的CORS跨域中间件。

```go
import "erp/pkg/middleware"

// 使用默认CORS中间件
r.Use(middleware.CORSMiddleware())
```

**默认配置**：
- 允许所有源 (`*`)
- 允许方法：GET, POST, PUT, DELETE, OPTIONS
- 允许头部：Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization

#### CORSMiddlewareWithConfig
可配置的CORS中间件。

```go
// 自定义CORS配置
origins := []string{"http://localhost:3000", "https://example.com"}
methods := []string{"GET", "POST", "PUT", "DELETE"}
headers := []string{"Authorization", "Content-Type"}

r.Use(middleware.CORSMiddlewareWithConfig(origins, methods, headers))
```

**参数**：
- `origins`: 允许的源列表
- `methods`: 允许的HTTP方法
- `headers`: 允许的请求头部

## 使用示例

### 基本使用

```go
package main

import (
    "erp/pkg/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 添加CORS中间件
    r.Use(middleware.CORSMiddleware())
    
    // 公开路由
    r.GET("/public", publicHandler)
    
    // 需要认证的路由
    auth := r.Group("/api")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.GET("/profile", profileHandler)
        auth.PUT("/profile", updateProfileHandler)
    }
    
    // 需要管理员权限的路由
    admin := r.Group("/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
    {
        admin.GET("/users", listUsersHandler)
    }
    
    r.Run(":8080")
}
```

### 自定义配置

```go
package main

import (
    "erp/pkg/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 自定义CORS配置
    corsConfig := middleware.CORSMiddlewareWithConfig(
        []string{"http://localhost:3000", "https://myapp.com"},
        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        []string{"Authorization", "Content-Type", "X-Requested-With"},
    )
    r.Use(corsConfig)
    
    // 其他路由配置...
    
    r.Run(":8080")
}
```

## 依赖关系

中间件包依赖以下包：
- `erp/pkg/auth` - JWT认证工具
- `erp/pkg/response` - 统一响应格式
- `github.com/gin-gonic/gin` - Gin框架

## 扩展中间件

要添加新的中间件，只需在 `pkg/middleware/` 目录下创建新的文件：

```go
// pkg/middleware/logger.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
    "time"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```

## 最佳实践

1. **中间件顺序**：CORS中间件应该在其他中间件之前
2. **错误处理**：使用统一的响应格式
3. **性能考虑**：避免在中间件中进行重计算
4. **可配置性**：提供默认配置和自定义配置选项
5. **文档化**：为每个中间件提供清晰的使用说明

## 测试

```go
// 测试中间件
func TestAuthMiddleware(t *testing.T) {
    // 测试代码...
}
```

中间件包设计为可复用组件，可以在不同的项目中导入使用。 