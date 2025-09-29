# Push 事件直接评论功能测试

## 功能说明

现在 Push 事件会直接在 commit 的评论框中添加评论，而不是创建 Issue。

## 修改内容

### 1. 通知服务更新

**文件**: `backend/internal/service/gitlab_service/notification_service.go`

**修改**: `SendPushEventComment` 方法
- **之前**: 根据 `commentType` 选择创建 Issue 或 Discussion
- **现在**: 统一使用 commit 评论，不创建 Issue

```go
// 修改前
switch commentType {
case utils.CommentTypeCommon:
    // 创建 Issue
case utils.CommentTypeInline:
    // 创建 commit 评论
}

// 修改后
// Push 事件统一使用 commit 评论，不创建 Issue
if commitSHA != "" {
    // 为 commit 创建评论
    CreateCommitComment(api, projectID, commitSHA, token, title, comments)
} else {
    // 如果没有 commit SHA，创建项目 Discussion
    CreateDiscussionForPush(api, projectID, token, title, comments)
}
```

### 2. 评论方式

**Push 事件现在会**:
1. **有 commit SHA**: 直接在 commit 的评论框中添加评论
2. **没有 commit SHA**: 创建项目 Discussion（作为备选方案）

## 测试场景

### 场景1: 单 commit push
- **输入**: 单个 commit 的 push 事件
- **预期**: 在该 commit 的评论框中添加 AI 审查结果
- **API**: `POST /v4/projects/{project_id}/repository/commits/{commit_sha}/comments`

### 场景2: 多 commit push
- **输入**: 多个 commit 的 push 事件
- **预期**: 为每个 commit 分别添加评论
- **API**: 为每个 commit 调用 commit 评论 API

### 场景3: 无 commit 信息
- **输入**: 没有 commit 信息的 push 事件
- **预期**: 创建项目 Discussion
- **API**: `POST /v4/projects/{project_id}/discussions`

## 评论内容格式

```
## AI代码审查 - Push to main

**提交信息**: 修复用户登录bug

**审查结果**:
- src/auth.go:25: 建议添加输入验证
- src/user.go:10: 密码加密方式需要更新

---
**请评价此评论的质量：**
- [ ] 评论内容准确且有用
- [ ] 精准定位并提供建议
- [ ] 建议不够具体
- [ ] 完全误导性建议
```

## 验证方法

1. **触发 Push 事件**: 向 GitLab 仓库推送代码
2. **检查评论**: 在 GitLab 的 commit 页面查看是否有 AI 评论
3. **确认位置**: 评论应该出现在 commit 的评论区域，而不是 Issues 页面

## 注意事项

- **不再创建 Issue**: Push 事件不会在 Issues 页面创建新的 Issue
- **直接评论**: 所有 AI 审查结果都会直接添加到 commit 的评论中
- **保持 webhook**: 仍然会发送 webhook 通知到企业微信等外部系统
- **错误处理**: 如果 commit 评论失败，会记录错误但不中断流程

## 配置要求

- GitLab Token 需要有 `api` 权限
- 需要能够访问 `POST /v4/projects/{project_id}/repository/commits/{commit_sha}/comments` API
- 项目需要允许 commit 评论功能
