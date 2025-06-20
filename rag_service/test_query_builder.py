#!/usr/bin/env python3
"""
测试组合查询构建功能
"""

from rag_engine import extract_query_from_diff, build_optimal_query

def test_extract_query_from_diff():
    """测试从diff中提取查询"""
    print("=== 测试从diff中提取查询 ===")
    
    # 测试用例1: Python函数
    diff1 = """
diff --git a/src/main.py b/src/main.py
index 1234567..abcdefg 100644
--- a/src/main.py
+++ b/src/main.py
@@ -1,3 +1,3 @@
-def old_function():
+def new_function():
     return True
"""
    
    query1 = extract_query_from_diff(diff1)
    print(f"测试1 - Python函数:")
    print(f"Diff: {diff1.strip()}")
    print(f"提取的查询: {query1}")
    print()
    
    # 测试用例2: JavaScript API
    diff2 = """
diff --git a/api/user.js b/api/user.js
index 1111111..2222222 100644
--- a/api/user.js
+++ b/api/user.js
@@ -1,5 +1,5 @@
-export function getUserInfo() {
+export function getUserInfo(userId) {
     const response = await fetch('/api/users/' + userId);
     return response.json();
 }
"""
    
    query2 = extract_query_from_diff(diff2)
    print(f"测试2 - JavaScript API:")
    print(f"Diff: {diff2.strip()}")
    print(f"提取的查询: {query2}")
    print()
    
    # 测试用例3: 数据库操作
    diff3 = """
diff --git a/models/user.py b/models/user.py
index 3333333..4444444 100644
--- a/models/user.py
+++ b/models/user.py
@@ -1,3 +1,3 @@
-class User:
+class User:
     def __init__(self):
         self.db = Database()
     
+    def save_user(self, user_data):
+        sql = "INSERT INTO users (name, email) VALUES (?, ?)"
+        return self.db.execute(sql, user_data)
"""
    
    query3 = extract_query_from_diff(diff3)
    print(f"测试3 - 数据库操作:")
    print(f"Diff: {diff3.strip()}")
    print(f"提取的查询: {query3}")
    print()

def test_build_optimal_query():
    """测试构建最优查询"""
    print("=== 测试构建最优查询 ===")
    
    # 测试用例
    diff_content = """
diff --git a/api/auth.py b/api/auth.py
index 5555555..6666666 100644
--- a/api/auth.py
+++ b/api/auth.py
@@ -1,3 +1,3 @@
-def login():
+def login(username, password):
     # 验证用户登录
     token = generate_jwt_token(username)
     return {"token": token}
"""
    
    original_query = "请检查代码是否符合以下要求：\n1. 代码风格是否规范\n2. 是否有潜在的安全问题"
    
    merge_request_info = {
        "title": "修复用户登录功能",
        "description": "添加用户名和密码参数，优化JWT token生成逻辑"
    }
    
    optimal_query = build_optimal_query(diff_content, original_query, merge_request_info)
    
    print(f"Diff内容: {diff_content.strip()}")
    print(f"原始查询: {original_query}")
    print(f"MR信息: {merge_request_info}")
    print(f"最优查询: {optimal_query}")
    print()

def test_semantic_keywords():
    """测试语义关键词提取"""
    print("=== 测试语义关键词提取 ===")
    
    test_cases = [
        {
            "name": "API接口",
            "diff": "diff --git a/api/users.py b/api/users.py\n+def get_user_api():\n+    response = requests.get('/api/users')\n+    return response.json()"
        },
        {
            "name": "数据库操作", 
            "diff": "diff --git a/db/user.py b/db/user.py\n+def save_user():\n+    sql = 'INSERT INTO users VALUES (?, ?)'\n+    return db.execute(sql)"
        },
        {
            "name": "认证授权",
            "diff": "diff --git a/auth/login.py b/auth/login.py\n+def authenticate():\n+    token = jwt.encode(payload, secret)\n+    return token"
        },
        {
            "name": "错误处理",
            "diff": "diff --git a/utils/error.py b/utils/error.py\n+try:\n+    result = process_data()\n+except Exception as e:\n+    log_error(e)"
        }
    ]
    
    for case in test_cases:
        query = extract_query_from_diff(case["diff"])
        print(f"{case['name']}: {query}")

if __name__ == "__main__":
    print("开始测试组合查询构建功能...")
    print("=" * 60)
    
    test_extract_query_from_diff()
    test_build_optimal_query()
    test_semantic_keywords()
    
    print("=" * 60)
    print("测试完成") 