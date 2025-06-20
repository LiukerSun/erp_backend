# ERP 后端系统

一个基于 Go 语言构建的企业资源规划（ERP）后端系统，采用现代化的架构设计和最佳实践。

## 🚀 特性

- **高性能**: 基于 Gin 框架，支持高并发
- **模块化设计**: 清晰的分层架构，易于扩展
- **JWT 认证**: 安全的用户认证机制
- **API 文档**: 自动生成的 Swagger 文档
- **数据库支持**: PostgreSQL + GORM ORM
- **配置管理**: 环境变量配置，支持开发/生产环境
- **中间件**: CORS、认证、权限控制等中间件

## 📋 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **认证**: JWT
- **文档**: Swagger/OpenAPI
- **配置**: godotenv

## 🛠️ 安装和运行

### 前置要求

- Go 1.21 或更高版本
- PostgreSQL 数据库
- Make (可选，用于使用 Makefile 命令)

### 1. 克隆项目

```bash
git clone <repository-url>
cd erp
```

### 2. 安装依赖

```bash
make install-deps
# 或者
go mod tidy
```

### 3. 配置环境变量

创建 `.env` 文件：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=erp_db

# JWT 配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug
```

### 4. 运行项目

```bash
# 开发模式（推荐）
make dev

# 或者直接运行
make run

# 或者
go run cmd/server/main.go
```

### 5. 访问服务

- **API 服务**: http://localhost:8080
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health

## 📚 API 文档

### 用户管理接口

#### 注册用户
```http
POST /api/user/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 用户登录
```http
POST /api/user/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### 获取用户信息
```http
GET /api/user/profile
Authorization: Bearer <jwt_token>
```

#### 更新用户信息
```http
PUT /api/user/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "full_name": "新姓名"
}
```

#### 修改密码
```http
POST /api/user/change_password
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "old_password": "oldpassword",
  "new_password": "newpassword"
}
```

#### 获取用户列表（管理员）
```http
GET /api/user/list?page=1&limit=10
Authorization: Bearer <jwt_token>
```

## 🏗️ 项目结构

```
erp/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── config/
│   └── config.go            # 配置管理
├── docs/                    # Swagger 文档
├── internal/
│   ├── app/
│   │   └── app.go          # 应用管理器
│   └── modules/
│       └── user/           # 用户模块
│           ├── handler/    # HTTP 处理器
│           ├── model/      # 数据模型
│           ├── repository/ # 数据访问层
│           ├── service/    # 业务逻辑层
│           └── module.go   # 模块定义
├── pkg/                    # 公共包
│   ├── auth/              # JWT 认证
│   ├── database/          # 数据库连接
│   ├── middleware/        # 中间件
│   ├── password/          # 密码处理
│   └── response/          # 响应处理
├── routes/
│   └── routes.go          # 路由配置
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔧 开发命令

```bash
# 安装依赖
make install-deps

# 生成 Swagger 文档
make swagger

# 构建项目
make build

# 运行项目
make run

# 开发模式（热重载）
make dev

# 运行测试
make test

# 清理构建文件
make clean

# 安装 Air（热重载工具）
make install-air
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/modules/user/...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 🔒 安全特性

- JWT 令牌认证
- 密码加密存储
- CORS 中间件
- 角色权限控制
- 输入验证和清理

## 🚀 部署

### Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o erp cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/erp .
EXPOSE 8080
CMD ["./erp"]
```

### 生产环境配置

1. 设置 `SERVER_MODE=release`
2. 使用强密码的 JWT_SECRET
3. 配置生产环境数据库
4. 设置适当的 CORS 策略
5. 启用 HTTPS

## 🤝 贡献

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 支持

如果您有任何问题或建议，请：

1. 查看 [Issues](../../issues)
2. 创建新的 Issue
3. 联系开发团队

---

**注意**: 这是一个开发中的项目，API 可能会发生变化。请查看最新的 Swagger 文档获取最新的 API 信息。 