# ERP 后端系统

## 📖 项目简介

这是一个基于 Go 语言构建的企业资源规划（ERP）后端系统，采用模块化架构设计，提供完整的 RESTful API 接口。系统使用 JWT 进行用户认证，PostgreSQL 作为数据库，支持 Docker 容器化部署。

## 🚀 技术栈

- **后端框架**: Go 1.23 + Gin Framework
- **数据库**: PostgreSQL + GORM ORM
- **认证**: JWT (JSON Web Token)
- **API 文档**: Swagger/OpenAPI 3.0
- **微服务**: gRPC + Protocol Buffers
- **对象存储**: 阿里云 OSS
- **容器化**: Docker + Docker Compose
- **密码加密**: bcrypt
- **配置管理**: Environment Variables + dotenv

## 🎯 功能模块

### 用户管理模块
- ✅ 用户注册/登录
- ✅ JWT 认证和刷新
- ✅ 用户资料管理
- ✅ 密码修改
- ✅ 角色权限管理
- ✅ 管理员功能

### 供应商管理模块
- ✅ 供应商信息管理
- ✅ 供应商状态管理
- ✅ 供应商查询和筛选

### 店铺管理模块
- ✅ 店铺信息管理
- ✅ 店铺状态管理
- ✅ 供应商-店铺关联管理

### Excel 文件处理模块
- ✅ Excel 文件上传和解析
- ✅ 商品信息提取
- ✅ gRPC 服务支持

## 📁 项目结构

```
erp_backend/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── config/
│   └── config.go                # 配置管理
├── internal/
│   ├── app/
│   │   └── app.go               # 应用管理器
│   └── modules/
│       ├── user/                # 用户模块
│       ├── supplier/            # 供应商模块
│       ├── store/               # 店铺模块
│       └── excel/               # Excel处理模块
├── pkg/
│   ├── auth/                    # JWT认证
│   ├── database/                # 数据库连接
│   ├── middleware/              # 中间件
│   ├── password/                # 密码加密
│   ├── proto/                   # gRPC协议
│   └── response/                # 响应格式
├── routes/
│   └── routes.go                # 路由配置
├── docs/                        # Swagger文档
├── docker-compose.yml           # Docker编排
├── Dockerfile                   # Docker构建
├── env.example                  # 环境变量示例
├── go.mod                       # Go模块
└── README.md                    # 项目说明
```

## 🔧 环境要求

- Go 1.23 或更高版本
- PostgreSQL 12 或更高版本
- Docker 和 Docker Compose (可选)

## 📦 安装和运行

### 1. 克隆项目
```bash
git clone git@github.com:LiukerSun/erp_backend.git
cd erp_backend
```

### 2. 配置环境变量
```bash
cp env.example .env
# 编辑 .env 文件，配置数据库和其他参数
```

### 3. 安装依赖
```bash
go mod download
```

### 4. 运行应用
```bash
go run cmd/server/main.go
```

### 5. 使用 Docker 运行
```bash
docker-compose up -d
```

## 🗄️ 数据库配置

确保 PostgreSQL 已安装并运行，然后创建数据库：

```sql
CREATE DATABASE erp_db;
CREATE USER erp_dev WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE erp_db TO erp_dev;
```

## 🔑 环境变量配置

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=erp_db

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=release

# JWT配置
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_DAYS=30

# gRPC配置
GRPC_SERVER_ADDR=localhost:50051

# OSS配置（可选）
OSS_ACCESS_KEY_ID=your_access_key_id
OSS_ACCESS_KEY_SECRET=your_access_key_secret
OSS_BUCKET_NAME=your_bucket_name
OSS_REGION=cn-beijing
OSS_ROLE_ARN=your_role_arn
OSS_ROLE_SESSION_NAME=erp-frontend-upload
```

## 📚 API 文档

启动服务后，访问以下地址查看 API 文档：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health

### 主要 API 端点

#### 用户管理
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录
- `POST /api/user/refresh` - 刷新令牌
- `GET /api/user/profile` - 获取用户信息
- `POST /api/user/change_password` - 修改密码

#### 供应商管理
- `GET /api/suppliers` - 获取供应商列表
- `POST /api/suppliers` - 创建供应商
- `GET /api/suppliers/:id` - 获取供应商详情
- `PUT /api/suppliers/:id` - 更新供应商
- `DELETE /api/suppliers/:id` - 删除供应商

#### 店铺管理
- `GET /api/stores` - 获取店铺列表
- `POST /api/stores` - 创建店铺
- `GET /api/stores/:id` - 获取店铺详情
- `PUT /api/stores/:id` - 更新店铺
- `DELETE /api/stores/:id` - 删除店铺

#### Excel 处理
- `POST /api/excel/parse` - 解析 Excel 文件

## 🔐 身份认证

系统使用 JWT 进行身份认证，请在请求头中添加：

```
Authorization: Bearer <your-jwt-token>
```

## 🐳 Docker 部署

### 构建镜像
```bash
docker build -t erp-backend .
```

### 使用 Docker Compose
```bash
docker-compose up -d
```

### 健康检查
```bash
curl http://localhost:8080/health
```

## 🔨 开发

### 生成 Swagger 文档
```bash
swag init -g cmd/server/main.go
```

### 生成 gRPC 代码
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/proto/excel.proto
```
