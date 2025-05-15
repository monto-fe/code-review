-- 创建数据库实例
CREATE DATABASE IF NOT EXISTS ucode_review;
-- 设置数据库字符集和排序规则
ALTER DATABASE ucode_review CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;  
-- 为用户赋予权限
GRANT ALL PRIVILEGES ON ucode_review.* TO 'root'@'%'; 
FLUSH PRIVILEGES;