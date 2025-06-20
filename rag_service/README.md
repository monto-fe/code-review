# RAG 代码审查服务

基于 RAG (Retrieval-Augmented Generation) 技术的智能代码审查服务，支持 GitLab 代码分析和审查建议生成。

## 功能特性

- 🔍 智能代码分析：基于向量相似度搜索相关代码上下文
- 📝 自动审查建议：生成结构化的代码审查报告
- 🔗 GitLab 集成：支持 GitLab 仓库和合并请求
- ⚡ 高性能：使用 FAISS 向量数据库和轻量级嵌入模型
- 🐳 Docker 部署：容器化部署，开箱即用

## 快速开始

### 使用 Docker Compose 启动

```bash
# 启动服务
./start.sh

# 或者手动启动
docker-compose up -d rag-service
```

### 重新构建服务

```bash
# 重新构建并启动
./rebuild_and_start.sh
```

## API 接口

### 代码审查接口

**POST** `/api/code-analysis`

请求参数：
```json
{
  "git_url": "https://gitlab.com/user/repo.git",
  "branch": "main",
  "diff_content": "diff --git a/file.py b/file.py\n@@ -1,3 +1,3 @@\n-def old_function():\n+def new_function():\n     return True",
  "query": "查找数据库相关代码",
  "gitlab_token": "glpat-xxxxxxxx"
}
```

响应：
```json
{
  "review": "代码变更:\n...\n\n相关上下文:\n...\n\n最优查询:\n..."
}
```

### 健康检查

**GET** `/health`

## 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `GITLAB_TOKEN` | - | GitLab 访问令牌 |
| `RAG_SERVICE_HOST` | `0.0.0.0` | 服务监听地址 |
| `RAG_SERVICE_PORT` | `8000` | 服务端口 |
| `VECTOR_STORE_TYPE` | `faiss` | 向量存储类型 |
| `VECTOR_STORE_PATH` | `/app/data/vector_store` | 向量存储路径 |
| `GIT_TEMP_DIR` | `/tmp/git_repos` | Git 临时目录 |
| `SKIP_DEPS_CHECK` | `true` | 跳过依赖检查 |

## 技术栈

- **Web框架**: FastAPI
- **向量数据库**: FAISS
- **嵌入模型**: sentence-transformers/all-MiniLM-L6-v2
- **Git操作**: GitPython
- **容器化**: Docker

## 服务地址

- 服务地址: http://localhost:8000
- API文档: http://localhost:8000/docs
- 健康检查: http://localhost:8000/health 