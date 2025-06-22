# RAG服务集成说明

## 概述

本项目已集成RAG（Retrieval-Augmented Generation）服务，用于代码审查。RAG服务能够基于代码库的上下文信息提供更准确的代码审查建议。

## RAG服务配置

### 环境变量

设置RAG服务URL：

```bash
export RAG_SERVICE_URL="http://localhost:8000"
```

如果不设置，系统将使用默认URL：`http://localhost:8000`

### RAG服务接口

RAG服务提供以下接口：

**POST** `/api/code-analysis`

请求体：
```json
{
    "git_url": "http://localhost:9980/testuser/test.git",
    "branch": "main", 
    "diff_content": "File: index.js\n@@ -2,14 +2,6 @@ function a(){wess",
    "query": "检查",
    "gitlab_token": "glpat-xxxxxxxxxxxx"
}
```

响应：
```json
{
    "review": "代码变更:\nFile: index.js\n@@ -2,14 +2,6 @@ function a(){wess\n\n相关上下文:\n## e built-in continuous integration in GitLab.\n\n最优查询:\n检查"
}
```

## 代码集成

### 主要函数

1. **CheckMergeRequestWithRAG** - 使用RAG服务检查合并请求
2. **NewRAGClient** - 创建RAG客户端
3. **AnalyzeCodeWithRequest** - 调用RAG服务进行分析

### 使用流程

1. **Webhook触发** - GitLab合并请求触发webhook
2. **RAG检查** - 优先使用RAG服务进行代码审查
3. **回退机制** - 如果RAG服务失败，回退到传统AI检查
4. **结果处理** - 保存审查结果并发送到GitLab

### 代码示例

```go
// 创建RAG客户端
ragClient := NewRAGClient("http://localhost:8000")

// 准备请求
req := &CodeReviewRequest{
    GitURL:      "http://gitlab.com/project/repo.git",
    Branch:      "main",
    DiffContent: "diff --git a/file.js b/file.js\n@@ -1,1 +1,1 @@\n-old code\n+new code",
    Query:       "检查代码质量和安全性",
    GitlabToken: "your-gitlab-token",
}

// 调用RAG服务
result, err := ragClient.AnalyzeCodeWithRequest(req)
if err != nil {
    return "", fmt.Errorf("RAG服务调用失败: %v", err)
}

// 获取审查结果
review := result.Review
```

## 测试

运行RAG服务集成测试：

```bash
cd backend
go test -v ./internal/service -run TestRAGServiceIntegration
```

## 错误处理

RAG服务包含以下错误处理机制：

1. **超时处理** - 3分钟超时
2. **网络错误** - 自动重试和错误报告
3. **服务不可用** - 回退到传统AI检查
4. **响应解析错误** - 详细的错误信息

## 监控和日志

RAG服务调用会记录以下信息：

- 请求开始时间
- 响应时间
- 错误信息
- 审查结果摘要

## 性能优化

1. **连接池** - HTTP客户端复用连接
2. **超时控制** - 避免长时间等待
3. **异步处理** - 不阻塞webhook响应
4. **缓存机制** - 避免重复请求

## 故障排除

### 常见问题

1. **RAG服务不可达**
   - 检查网络连接
   - 验证RAG服务URL
   - 确认防火墙设置

2. **超时错误**
   - 检查RAG服务性能
   - 调整超时时间
   - 优化请求内容

3. **认证失败**
   - 验证GitLab Token
   - 检查Token权限
   - 确认项目访问权限

### 调试模式

启用详细日志：

```bash
export DEBUG_RAG=true
```

## 更新日志

- **v1.0.0** - 初始RAG服务集成
- **v1.1.0** - 添加回退机制和错误处理
- **v1.2.0** - 优化性能和监控 