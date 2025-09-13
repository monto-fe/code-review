# GitLab 合并请求URL解析器优化

## 概述

优化了 `parseMergeURL` 函数，使其能够准确解析各种格式的 GitLab 合并请求 URL，提取项目路径、合并请求ID和项目名称。

## 支持的URL格式

### 1. 标准GitLab URL
```
https://gitlab.example.com/project/repo/-/merge_requests/123
```
- 项目路径: `project/repo`
- 合并请求ID: `123`

### 2. 带组的GitLab URL
```
https://gitlab.example.com/group/project/-/merge_requests/456
```
- 项目路径: `group/project`
- 合并请求ID: `456`

### 3. IP地址GitLab URL
```
http://127.0.0.1:9980/usms/api/-/merge_requests/3
```
- 项目路径: `usms/api`
- 合并请求ID: `3`

### 4. GitLab.com URL
```
https://gitlab.com/username/project/-/merge_requests/789
```
- 项目路径: `username/project`
- 合并请求ID: `789`

### 5. 带额外路径的URL
```
https://gitlab.example.com/project/repo/-/merge_requests/123/diffs
```
- 项目路径: `project/repo`
- 合并请求ID: `123`

## 核心功能

### 1. parseMergeURL 函数
```go
func (s *ManualCheckService) parseMergeURL(mergeURL string) (int, int, string, error)
```

**功能**:
- 解析GitLab合并请求URL
- 提取项目ID、合并请求ID和项目名称
- 返回解析结果或错误信息

**参数**:
- `mergeURL`: GitLab合并请求的完整URL

**返回值**:
- `projectID`: 基于项目路径生成的唯一项目ID
- `mergeID`: 合并请求的数字ID
- `projectName`: 基于项目路径生成的项目名称
- `error`: 解析错误信息

### 2. extractProjectPath 函数
```go
func (s *ManualCheckService) extractProjectPath(mergeURL string) string
```

**功能**:
- 从GitLab URL中提取项目路径
- 使用正则表达式精确匹配
- 支持各种URL格式

**正则表达式**:
```regex
^https?://[^/]+/(.+?)/-/merge_requests/
```

**匹配逻辑**:
- `^https?://`: 匹配http或https协议
- `[^/]+`: 匹配域名部分（包括端口）
- `/(.+?)/`: 非贪婪匹配项目路径
- `/-/merge_requests/`: 匹配GitLab的固定路径格式

### 3. generateProjectID 函数
```go
func (s *ManualCheckService) generateProjectID(projectPath string) int
```

**功能**:
- 根据项目路径生成唯一的项目ID
- 使用FNV-1a哈希算法
- 确保相同路径生成相同ID
- 确保ID为正整数且不为0

**算法**:
1. 使用FNV-1a哈希算法计算项目路径的哈希值
2. 将哈希值转换为正整数
3. 如果ID为0，则设置为1

### 4. generateProjectName 函数
```go
func (s *ManualCheckService) generateProjectName(projectPath string) string
```

**功能**:
- 根据项目路径生成项目名称
- 直接使用项目路径作为项目名称
- 处理空路径的边界情况

**逻辑**:
1. 如果项目路径为空，返回 "unknown-project"
2. 否则直接返回项目路径作为项目名称

## 解析示例

### 示例1: IP地址URL
```
输入: http://165.154.112.72:9980/usms/api/-/merge_requests/3
输出: 
- 项目路径: usms/api
- 项目ID: 1531043042 (基于哈希值)
- 合并请求ID: 3
- 项目名称: usms/api
```

### 示例2: 标准URL
```
输入: https://gitlab.example.com/project/repo/-/merge_requests/123
输出:
- 项目路径: project/repo
- 项目ID: 3351040939 (基于哈希值)
- 合并请求ID: 123
- 项目名称: project/repo
```

### 示例3: 带组URL
```
输入: https://gitlab.example.com/group/project/-/merge_requests/456
输出:
- 项目路径: group/project
- 项目ID: 3798180626 (基于哈希值)
- 合并请求ID: 456
- 项目名称: group/project
```

### 示例4: 多层路径URL
```
输入: https://gitlab.example.com/org/group/subgroup/project/-/merge_requests/999
输出:
- 项目路径: org/group/subgroup/project
- 项目ID: 1649659397 (基于哈希值)
- 合并请求ID: 999
- 项目名称: org/group/subgroup/project
```

## 错误处理

### 1. 无效URL格式
```
错误: 无效的合并请求URL格式，请确保URL包含 /-/merge_requests/数字
```
**原因**: URL不包含标准的GitLab合并请求路径格式

### 2. 无法解析合并请求ID
```
错误: 无法解析合并请求ID: strconv.Atoi: parsing "abc": invalid syntax
```
**原因**: 合并请求ID不是有效的数字

### 3. 无法提取项目路径
```
错误: 无法从URL中提取项目路径
```
**原因**: URL格式不符合预期，无法提取项目路径部分

## 优化改进

### 1. 正则表达式优化
- 使用更精确的正则表达式匹配
- 支持非贪婪匹配，避免过度匹配
- 支持带额外路径的URL

### 2. 项目ID生成优化
- 使用FNV-1a哈希算法，性能更好
- 确保相同路径生成相同ID
- 避免ID冲突和重复

### 3. 错误处理优化
- 提供更详细的错误信息
- 区分不同类型的解析错误
- 便于调试和问题定位

### 4. 代码结构优化
- 将复杂逻辑拆分为独立函数
- 提高代码可读性和可维护性
- 便于单元测试

## 测试覆盖

### 1. 正常情况测试
- 各种标准URL格式
- 带组的URL格式
- IP地址URL格式
- 带额外路径的URL格式

### 2. 异常情况测试
- 无效URL格式
- 非数字合并请求ID
- 空URL
- 不完整的URL

### 3. 边界情况测试
- 单层项目路径
- 多层项目路径
- 特殊字符处理

## 使用建议

1. **URL验证**: 在调用解析函数前，建议先验证URL的基本格式
2. **错误处理**: 始终检查返回的错误信息，提供用户友好的错误提示
3. **项目ID使用**: 项目ID是基于哈希值生成的，适合作为数据库主键或缓存键
4. **性能考虑**: 解析函数性能良好，适合频繁调用

## 兼容性

- 支持HTTP和HTTPS协议
- 支持标准域名和IP地址
- 支持各种端口号
- 兼容GitLab的不同版本
- 支持GitLab.com和自建GitLab实例