# `/webhook/merge` 调用函数链路分析

## 📋 整体调用链路

```
GitLab Webhook → /webhook/merge → 事件驱动架构 → 具体业务处理
```

## 🔄 详细调用链路

### 1. 初始化阶段 (应用启动时)

```
应用启动
    ↓
init() 函数执行
    ↓
service.NewEventRouter() 创建事件路由器
    ↓
eventRouter.RegisterHandler(service.NewMergeRequestHandler()) 注册 Merge Request 处理器
    ↓
eventRouter.RegisterHandler(service.NewPushHandler()) 注册 Push 事件处理器
    ↓
事件路由器初始化完成
```

### 2. Webhook 请求处理阶段

```
HTTP POST /webhook/merge
    ↓
AICheck(c *gin.Context) 函数
    ↓
【步骤1】解析请求体
    c.ShouldBindJSON(&body) → dto.WebhookBody
    ↓
【步骤2】基础验证
    body.Project.ID == 0 检查
    ↓
【步骤3】检测事件类型
    detectEventType(body) → "merge_request" | "push"
    ↓
【步骤4】立即响应 (不阻塞)
    response.Success(c, responseData, ...)
    ↓
【步骤5】事件路由
    eventRouter.Route(body)
```

### 3. 事件路由阶段

```
eventRouter.Route(body)
    ↓
【步骤1】检测事件类型
    detectEventType(body) → 事件类型字符串
    ↓
【步骤2】获取处理器
    r.handlers[eventType] → EventHandler 接口
    ↓
【步骤3】异步处理
    go func() { handler.Handle(body) }()
```

### 4. Merge Request 事件处理链路

```
MergeRequestHandler.Handle(body)
    ↓
【步骤1】防重复处理
    utils.IsDuplicateWebhook(body) → 缓存检查
    ↓
【步骤2】状态验证
    ShouldProcessState(body) → 检查 MR 状态
    ↓
【步骤3】分支过滤
    shouldProcessBranch(body) → 分支白名单/黑名单检查
    ↓
【步骤4】调用服务层
    service.ProcessMergeRequest(body)
```

### 5. Merge Request 服务处理链路

```
MergeRequestService.ProcessMergeRequest(body)
    ↓
【步骤1】获取 MR 信息
    getMergeRequestInfo(body) → MergeRequestInfo
    ↓
【步骤2】获取项目配置
    getProjectConfig(projectID) → ProjectConfig
    ↓
【步骤3】检查缓存
    getCachedResult(cacheKey) → 缓存结果检查
    ↓
【步骤4】执行 AI 检查
    performAICheck(body, mrInfo, config) → CheckMergeRequestWithAI(body)
    ↓
【步骤5】缓存结果
    cacheResult(cacheKey, result, ttl)
    ↓
【步骤6】发送结果
    sendMergeRequestResult(result, body)
```

### 6. Push 事件处理链路

```
PushHandler.Handle(body)
    ↓
【步骤1】验证 Push 事件
    shouldProcessPush(body) → 分支、提交数、文件变更检查
    ↓
【步骤2】调用服务层
    service.ProcessPush(body)
```

### 7. Push 服务处理链路

```
PushService.ProcessPush(body)
    ↓
【步骤1】分析提交信息
    analyzeCommits(body.Commits) → CommitAnalysis
    ↓
【步骤2】获取文件内容
    getFileContents(body) → map[string]string
    ↓
【步骤3】文件级分析
    analyzeFiles(fileContents) → FileAnalysis
    ↓
【步骤4】生成报告
    generateReport(commitAnalysis, fileAnalysis, body) → PushReport
    ↓
【步骤5】发送结果
    sendPushResults(report, body)
```

## 🏗️ 关键函数调用关系

### 主要入口函数
- `AICheck(c *gin.Context)` - Webhook 入口点
- `eventRouter.Route(body)` - 事件路由器
- `handler.Handle(body)` - 事件处理器接口

### 事件检测函数
- `detectEventType(body)` - 检测事件类型 (在 webhook.go 和 event_router.go 中都有实现)

### 业务处理函数
- `MergeRequestHandler.Handle()` - Merge Request 处理器
- `PushHandler.Handle()` - Push 事件处理器
- `MergeRequestService.ProcessMergeRequest()` - Merge Request 业务逻辑
- `PushService.ProcessPush()` - Push 事件业务逻辑

### 辅助函数
- `extractBranchFromRef(ref)` - 提取分支名称
- `shouldProcessBranch()` - 分支过滤逻辑
- `shouldProcessPush()` - Push 事件验证逻辑

## 📊 调用时序图

```
时间轴: 0ms → 10ms → 50ms → 100ms → 500ms → 1000ms

GitLab Webhook
    ↓ (0ms)
AICheck() 开始
    ↓ (5ms)
解析请求体 + 基础验证
    ↓ (10ms)
检测事件类型
    ↓ (15ms)
立即响应 (HTTP 200)
    ↓ (20ms)
eventRouter.Route() 异步启动
    ↓ (25ms)
handler.Handle() 开始处理
    ↓ (50ms)
业务逻辑处理 (AI检查、文件分析等)
    ↓ (500ms)
发送结果到 GitLab
    ↓ (1000ms)
处理完成
```

## 🔍 关键决策点

### 1. 事件类型检测
```go
func detectEventType(body dto.WebhookBody) string {
    // 优先级1: object_kind 字段
    if body.ObjectKind == "merge_request" { return "merge_request" }
    if body.ObjectKind == "push" { return "push" }
    
    // 优先级2: 字段特征检测
    if body.ObjectAttributes.IID > 0 && body.ObjectAttributes.Title != "" {
        return "merge_request"
    }
    if body.Ref != "" && body.After != "" {
        return "push"
    }
    
    // 优先级3: 默认处理
    return "merge_request"
}
```

### 2. 分支过滤逻辑
```go
func shouldProcessBranch(body dto.WebhookBody) bool {
    targetBranch := body.ObjectAttributes.TargetBranch
    
    // 白名单检查
    for _, branch := range supportedBranches {
        if targetBranch == branch { return true }
    }
    
    // 黑名单检查
    for _, pattern := range excludedBranches {
        if matchesPattern(targetBranch, pattern) { return false }
    }
    
    return true // 默认处理
}
```

### 3. 异步处理机制
```go
// 在 eventRouter.Route() 中
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("事件处理发生panic: %v\n", r)
        }
    }()
    
    err := handler.Handle(body)
    // 记录处理结果
}()
```

## 🛡️ 错误处理链路

### 1. Webhook 层错误处理
- JSON 解析错误 → HTTP 400
- 项目ID验证失败 → HTTP 400
- 事件路由失败 → 记录日志

### 2. 事件处理器层错误处理
- 重复请求 → 跳过处理
- 状态验证失败 → 跳过处理
- 分支过滤失败 → 跳过处理

### 3. 服务层错误处理
- 获取信息失败 → 返回错误
- AI检查失败 → 返回错误
- 发送结果失败 → 返回错误

### 4. 异步处理错误处理
- Panic 恢复 → 记录日志
- 处理失败 → 记录错误日志
- 处理成功 → 记录成功日志

## 📈 性能优化点

### 1. 非阻塞响应
- Webhook 立即返回 HTTP 200
- 业务处理在后台异步执行
- 避免 GitLab 超时

### 2. 缓存机制
- 防重复处理缓存
- 结果缓存减少重复计算
- 可配置的缓存 TTL

### 3. 并发处理
- 多个事件可以并发处理
- 每个事件独立处理，互不影响
- 支持高并发场景

## 🎯 扩展点

### 1. 新增事件类型
1. 创建新的 EventHandler 实现
2. 在 init() 中注册新处理器
3. 在 detectEventType() 中添加检测逻辑

### 2. 新增业务逻辑
1. 在对应的 Service 中添加新方法
2. 在 Handler 中调用新方法
3. 保持接口一致性

### 3. 新增配置选项
1. 在 Config 结构体中添加新字段
2. 在 Handler 中使用新配置
3. 支持动态配置更新

这个调用链路展示了完整的事件驱动架构，从 HTTP 请求到具体的业务处理，每个环节都有清晰的职责分工和错误处理机制。
