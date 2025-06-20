-- 添加password_version字段到users表
-- 这个字段用于跟踪密码版本，当密码修改时递增，使旧token失效

ALTER TABLE users ADD COLUMN password_version INT DEFAULT 1;

-- 为现有用户设置默认密码版本
UPDATE users SET password_version = 1 WHERE password_version IS NULL;

-- 添加索引以提高查询性能
CREATE INDEX idx_users_password_version ON users(password_version); 