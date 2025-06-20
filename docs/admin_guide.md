# 管理员功能使用指南

## 概述

本文档介绍了ERP系统中管理员用户管理功能的使用方法。管理员可以通过这些API来管理系统中的用户。

## 管理员权限

管理员具有以下权限：
- 查看所有用户列表
- 创建新用户（包括其他管理员）
- 更新用户信息（邮箱、角色、激活状态）
- 重置用户密码
- 删除用户（软删除）
- 禁用/启用用户

## 创建第一个管理员账户

### 方法1：使用Go脚本（推荐）

```bash
# 创建默认管理员账户（用户名：admin，邮箱：admin@example.com）
go run scripts/create_admin.go -password=your_secure_password

# 创建自定义管理员账户
go run scripts/create_admin.go \
  -username=youradmin \
  -email=youradmin@company.com \
  -password=your_secure_password
```

### 方法2：使用SQL脚本

1. 修改 `scripts/create_admin.sql` 中的密码哈希值
2. 在数据库中执行SQL脚本

## API 接口

### 基础URL
所有管理员API的基础URL为：`/api/user/admin`

### 认证
所有管理员API都需要：
1. Bearer Token认证
2. 管理员角色权限

### 接口列表

#### 1. 获取用户列表
```http
GET /api/user/admin/users?page=1&limit=10
Authorization: Bearer <your_admin_token>
```

#### 2. 创建用户
```http
POST /api/user/admin/users
Content-Type: application/json
Authorization: Bearer <your_admin_token>

{
  "username": "newuser",
  "email": "newuser@company.com",
  "password": "securepassword",
  "role": "user"  // 或 "admin"
}
```

#### 3. 更新用户信息
```http
PUT /api/user/admin/users/{user_id}
Content-Type: application/json
Authorization: Bearer <your_admin_token>

{
  "email": "newemail@company.com",     // 可选
  "role": "admin",                     // 可选：user 或 admin
  "is_active": false                   // 可选：禁用用户
}
```

#### 4. 重置用户密码
```http
POST /api/user/admin/users/{user_id}/reset_password
Content-Type: application/json
Authorization: Bearer <your_admin_token>

{
  "new_password": "newsecurepassword"
}
```

#### 5. 删除用户
```http
DELETE /api/user/admin/users/{user_id}
Authorization: Bearer <your_admin_token>
```

**注意：** 管理员不能删除自己的账户。

## 使用示例

### 1. 管理员登录
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_admin_password"
  }'
```

### 2. 创建新的管理员
```bash
curl -X POST http://localhost:8080/api/user/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "username": "manager",
    "email": "manager@company.com",
    "password": "managerpassword",
    "role": "admin"
  }'
```

### 3. 禁用用户
```bash
curl -X PUT http://localhost:8080/api/user/admin/users/123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "is_active": false
  }'
```

### 4. 重置用户密码
```bash
curl -X POST http://localhost:8080/api/user/admin/users/123/reset_password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "new_password": "newpassword123"
  }'
```

## 安全注意事项

1. **管理员账户安全**：
   - 使用强密码
   - 定期更换密码
   - 不要共享管理员凭据

2. **权限控制**：
   - 只给需要的用户分配管理员权限
   - 定期审查管理员账户

3. **操作审计**：
   - 所有管理员操作都应该记录日志
   - 定期检查管理员操作记录

4. **密码重置**：
   - 重置密码会使用户的所有现有Token失效
   - 通知用户密码已被重置

## 故障排除

### 常见错误

1. **权限不足 (403)**
   - 确保使用管理员账户的Token
   - 检查Token是否过期

2. **用户不存在 (404)**
   - 检查用户ID是否正确
   - 确保用户未被删除

3. **邮箱已存在 (409)**
   - 检查邮箱是否已被其他用户使用

4. **不能删除自己 (400)**
   - 管理员不能删除自己的账户
   - 需要其他管理员来删除

## 支持

如果遇到问题，请检查：
1. 服务器日志文件
2. 数据库连接状态
3. Token是否有效
4. 用户权限是否正确 