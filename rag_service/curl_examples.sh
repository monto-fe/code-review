#!/bin/bash

# query参数使用示例 - curl命令

echo "=== query参数使用示例 ==="
echo

# 示例1: 数据库相关查询
echo "示例1: 数据库相关查询"
echo "查询: 数据库操作 用户模型 数据持久化"
echo

curl -X POST "http://localhost:8000/api/code-analysis" \
  -H "Content-Type: application/json" \
  -d '{
    "git_url": "https://gitlab.com/example/project.git",
    "branch": "feature/database-update",
    "diff_content": "diff --git a/models/user.py b/models/user.py\n@@ -1,5 +1,8 @@\n class User:\n-    def __init__(self, name, email):\n+    def __init__(self, name, email, created_at=None):\n         self.name = name\n         self.email = email\n+        self.created_at = created_at or datetime.now()\n+\n+    def save(self):\n+        # 保存到数据库\n+        pass",
    "query": "数据库操作 用户模型 数据持久化",
    "gitlab_token": "glpat-xxxxxxxx"
  }'

echo
echo "---"
echo

# 示例2: 安全相关查询
echo "示例2: 安全相关查询"
echo "查询: 身份验证 密码验证 安全防护"
echo

curl -X POST "http://localhost:8000/api/code-analysis" \
  -H "Content-Type: application/json" \
  -d '{
    "git_url": "https://gitlab.com/example/project.git",
    "branch": "feature/auth-improvement",
    "diff_content": "diff --git a/auth/login.py b/auth/login.py\n@@ -10,7 +10,10 @@\n     def authenticate(self, username, password):\n-        return self.db.query(\"SELECT * FROM users WHERE username = %s\", username)\n+        # 使用参数化查询防止SQL注入\n+        user = self.db.query(\"SELECT * FROM users WHERE username = %s AND active = 1\", username)\n+        if user and self.verify_password(password, user.password_hash):\n+            return user\n+        return None",
    "query": "身份验证 密码验证 安全防护 SQL注入",
    "gitlab_token": "glpat-xxxxxxxx"
  }'

echo
echo "---"
echo

# 示例3: 性能相关查询
echo "示例3: 性能相关查询"
echo "查询: 缓存机制 性能优化 数据库查询优化"
echo

curl -X POST "http://localhost:8000/api/code-analysis" \
  -H "Content-Type: application/json" \
  -d '{
    "git_url": "https://gitlab.com/example/project.git",
    "branch": "feature/cache-optimization",
    "diff_content": "diff --git a/services/data_service.py b/services/data_service.py\n@@ -5,8 +5,12 @@\n     def get_user_data(self, user_id):\n-        return self.db.query(\"SELECT * FROM users WHERE id = %s\", user_id)\n+        # 添加缓存机制\n+        cache_key = f\"user_data_{user_id}\"\n+        cached_data = self.cache.get(cache_key)\n+        if cached_data:\n+            return cached_data\n+\n+        data = self.db.query(\"SELECT * FROM users WHERE id = %s\", user_id)\n+        self.cache.set(cache_key, data, ttl=3600)\n+        return data",
    "query": "缓存机制 性能优化 数据库查询优化",
    "gitlab_token": "glpat-xxxxxxxx"
  }'

echo
echo "---"
echo

# 示例4: 不使用query参数（对比）
echo "示例4: 不使用query参数（对比）"
echo "注意: 不使用query参数时，只使用变更文件的完整内容作为上下文"
echo

curl -X POST "http://localhost:8000/api/code-analysis" \
  -H "Content-Type: application/json" \
  -d '{
    "git_url": "https://gitlab.com/example/project.git",
    "branch": "feature/database-update",
    "diff_content": "diff --git a/models/user.py b/models/user.py\n@@ -1,5 +1,8 @@\n class User:\n-    def __init__(self, name, email):\n+    def __init__(self, name, email, created_at=None):\n         self.name = name\n         self.email = email\n+        self.created_at = created_at or datetime.now()\n+\n+    def save(self):\n+        # 保存到数据库\n+        pass",
    "gitlab_token": "glpat-xxxxxxxx"
  }'

echo
echo "=== 使用说明 ==="
echo "1. query参数是可选的，用于在代码库中进行语义搜索"
echo "2. 搜索结果会作为上下文提供给LLM，提高审查质量"
echo "3. 使用中文查询词效果更好，模型支持中文语义理解"
echo "4. 查询词要具体且相关，避免过于宽泛"
echo "5. 可以组合多个相关概念，提高搜索精度" 