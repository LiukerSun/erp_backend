-- 创建管理员账户的SQL脚本
-- 这是一个示例脚本，您需要根据实际情况修改

-- 注意：这里的密码 "admin123" 使用 bcrypt 加密
-- 如果您使用的是 golang 的 bcrypt 包，密码"admin123"的哈希值如下：
-- 您可以使用 Go 程序生成正确的哈希值

INSERT INTO users (
    username, 
    email, 
    password, 
    role, 
    is_active, 
    password_version,
    created_at, 
    updated_at
) VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$JNT5F2q7G4.QVjCcGlYrruNMJ0SgXJgJjmJ0zJgLz4JNfJDfGJE8G', -- 这是 "admin123" 的 bcrypt 哈希值
    'admin',
    true,
    1,
    NOW(),
    NOW()
) ON DUPLICATE KEY UPDATE
    email = VALUES(email),
    password = VALUES(password),
    role = VALUES(role),
    is_active = VALUES(is_active),
    updated_at = NOW();

-- 如果您使用 PostgreSQL，请使用以下语句：
-- INSERT INTO users (username, email, password, role, is_active, password_version, created_at, updated_at)
-- VALUES ('admin', 'admin@example.com', '$2a$10$JNT5F2q7G4.QVjCcGlYrruNMJ0SgXJgJjmJ0zJgLz4JNfJDfGJE8G', 'admin', true, 1, NOW(), NOW())
-- ON CONFLICT (username) DO UPDATE SET
--     email = EXCLUDED.email,
--     password = EXCLUDED.password,
--     role = EXCLUDED.role,
--     is_active = EXCLUDED.is_active,
--     updated_at = NOW(); 