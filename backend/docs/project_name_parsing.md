# 项目名称解析功能优化

## 概述

在原有的GitLab URL解析功能基础上，新增了项目名称解析功能，使系统能够从GitLab合并请求URL中提取并生成有意义的项目名称。

## 功能增强

### 1. 函数签名更新

**原始函数**:
```go
func (s *ManualCheckService) parseMergeURL(mergeURL string) (int, int, error)
```

**优化后函数**:
```go
func (s *ManualCheckService) parseMergeURL(mergeURL string) (int, int, string, error)
```

### 2. 新增返回值

- **项目ID**: 基于项目路径哈希值生成的唯一标识符
- **合并请求ID**: 从URL中提取的数字ID
- **项目名称**: 基于项目路径生成的可读项目名称
- **错误信息**: 解析过程中的错误详情

## 核心实现

### 1. generateProjectName 函数

```go
func (s *ManualCheckService) generateProjectName(projectPath string) string {
    // 如果项目路径为空，返回默认名称
    if projectPath == "" {
        return "unknown-project"
    }
    
    // 直接使用项目路径作为项目名称
    return projectPath
}
```

**设计理念**:
- 保持项目名称与项目路径的一致性
- 提供默认值处理边界情况
- 简单直接，易于理解和维护

### 2. 项目名称生成规则

| 项目路径 | 生成的项目名称 | 说明 |
|----------|----------------|------|
| `usms/api` | `usms/api` | 直接使用路径 |
| `group/project` | `group/project` | 包含组的路径 |
| `username` | `username` | 单层路径 |
| `org/group/subgroup/project` | `org/group/subgroup/project` | 多层路径 |
| `""` (空) | `unknown-project` | 默认名称 |

## 使用场景

### 1. 任务创建时

```go
// 解析URL获取项目信息
projectID, mergeID, projectName, err := s.parseMergeURL(req.MergeURL)

// 创建任务时直接使用解析出的项目名称
task := &model.ManualCheckTask{
    ProjectID:   projectID,
    ProjectName: projectName, // 使用解析出的项目名称
    MergeID:     mergeID,
    // ... 其他字段
}
```

### 2. 异步更新时

```go
// 异步更新任务信息
go func() {
    s.updateTaskWithCachedInfo(task.ID, projectID, projectName, req.AIModelID)
}()
```

### 3. 任务信息更新

```go
func (s *ManualCheckService) updateTaskWithCachedInfo(taskID uint, projectID int, projectName string, aiModelID uint) {
    // 使用解析出的项目名称，如果没有则生成默认名称
    if projectName == "" {
        projectName = fmt.Sprintf("project-%d", projectID)
    }
    
    updateData["project_name"] = projectName
    // ... 更新数据库
}
```

## 解析示例

### 示例1: 企业GitLab实例
```
输入URL: http://165.154.112.72:9980/usms/api/-/merge_requests/3
解析结果:
- 项目ID: 1531043042
- 合并请求ID: 3
- 项目名称: usms/api
```

### 示例2: 标准GitLab.com
```
输入URL: https://gitlab.com/username/project/-/merge_requests/789
解析结果:
- 项目ID: 1566607395
- 合并请求ID: 789
- 项目名称: username/project
```

### 示例3: 复杂组织架构
```
输入URL: https://gitlab.example.com/org/group/subgroup/project/-/merge_requests/999
解析结果:
- 项目ID: 1649659397
- 合并请求ID: 999
- 项目名称: org/group/subgroup/project
```

## 优势

### 1. 用户体验提升
- **可读性**: 项目名称比数字ID更直观
- **识别性**: 用户可以通过项目名称快速识别项目
- **一致性**: 项目名称与GitLab中的实际路径保持一致

### 2. 系统功能增强
- **数据完整性**: 任务记录包含完整的项目信息
- **查询便利**: 支持按项目名称进行查询和筛选
- **报告生成**: 便于生成包含项目名称的报告

### 3. 开发维护
- **代码清晰**: 函数职责明确，易于理解
- **扩展性**: 可以轻松扩展项目名称的生成规则
- **测试友好**: 每个函数都可以独立测试

## 测试覆盖

### 1. 单元测试
- 测试各种URL格式的解析
- 测试项目名称生成逻辑
- 测试边界情况处理

### 2. 集成测试
- 测试完整的任务创建流程
- 测试异步更新功能
- 测试数据库存储和读取

### 3. 边界测试
- 空项目路径处理
- 特殊字符处理
- 超长路径处理

## 兼容性

### 1. 向后兼容
- 保持原有API接口不变
- 现有功能不受影响
- 数据库结构兼容

### 2. 数据迁移
- 新任务自动包含项目名称
- 历史任务保持原有数据
- 可选择性地更新历史数据

## 未来扩展

### 1. 项目名称优化
- 支持自定义项目名称映射
- 支持项目名称的本地化
- 支持项目名称的格式化

### 2. 智能识别
- 根据项目路径自动识别项目类型
- 支持项目标签和分类
- 支持项目描述信息

### 3. 缓存优化
- 缓存项目名称映射关系
- 支持项目信息的批量更新
- 优化数据库查询性能

## 总结

项目名称解析功能的添加，显著提升了系统的用户体验和功能完整性：

1. **功能完整**: 从URL中提取完整的项目信息
2. **用户友好**: 提供可读的项目名称
3. **系统健壮**: 处理各种边界情况
4. **易于维护**: 代码结构清晰，便于扩展

这个优化使得手动审核任务能够提供更丰富的项目信息，为用户提供更好的使用体验。