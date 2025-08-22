# 项目RetCode返回码列表

## 概述

本文档整理了项目中所有使用的RetCode返回码，包括成功码、错误码和状态码。所有返回码都已去重并按类型分类。

**重要更新**: 项目现在使用统一的 `internal/pkg/constants` 包来管理所有RetCode，请使用新的常量系统。

## 新的RetCode管理系统

### 包位置
```
internal/pkg/constants/ret_code.go
```

### 使用方式
```go
import "code-review-go/internal/pkg/constants"

// 成功响应
response.Success(c, data, "操作成功", int(constants.RetCodeSuccess))

// 错误响应
response.Error(c, err, "参数错误", int(constants.RetCodeInvalidParams))
```

### 详细使用说明
请参考: `internal/pkg/constants/README.md`

## 返回码列表

### 成功码 (Success Codes)

| RetCode | 常量名 | 说明 | 使用场景 |
|---------|--------|------|----------|
| 0 | `RetCodeSuccess` | 成功 | 所有成功操作的默认返回码 |

### 系统级错误码 (System Error Codes)

| RetCode | 常量名 | 说明 | 使用场景 |
|---------|--------|------|----------|
| 10000 | `RetCodeNotLogIn` | 未登录 | 用户未登录或Token无效时的默认错误码 |
| 10001 | `RetCodeInvalidParams` | 无效请求参数 | 请求参数验证失败 |
| 10002 | `RetCodeInvalidCredentials` | 无效用户名或密码 | 用户登录认证失败 |
| 10003 | `RetCodeLoginFailed` | 登录失败 | 用户登录过程失败 |
| 10004 | `RetCodeGetUserListFailed` | 获取用户列表失败 | 查询用户列表时发生错误 |
| 10005 | `RetCodeGetUserInfoFailed` | 获取用户信息失败 | 查询单个用户信息时发生错误 |
| 10006 | `RetCodeOperatorMissing` | 操作者信息缺失 | 请求头中缺少操作者信息 |
| 10007 | `RetCodeInvalidUserID` | 无效用户ID格式 | 用户ID格式不正确 |
| 10009 | `RetCodeCheckUserExistenceFailed` | 检查用户存在性失败 | 验证用户是否存在时发生错误 |
| 10010 | `RetCodeUserAlreadyExists` | 用户已存在 | 创建用户时用户已存在 |
| 10011 | `RetCodeUserOperationFailed` | 用户操作失败 | 创建、更新用户时发生错误 |
| 10012 | `RetCodeTokenMissing` | Token缺失 | 请求头中缺少Token |
| 10013 | `RetCodeDeleteUserFailed` | 删除用户失败 | 删除用户时发生错误 |
| 10014 | `RetCodeInvalidToken` | 无效Token | Token格式或内容无效 |
| 10016 | `RetCodeTokenExpired` | Token过期 | Token已过期 |

### HTTP状态码 (HTTP Status Codes)

| RetCode | 常量名 | 说明 | 使用场景 |
|---------|--------|------|----------|
| 400 | `RetCodeBadRequest` | 请求参数错误 | 参数验证失败、业务逻辑错误 |
| 500 | `RetCodeInternalError` | 服务器内部错误 | 数据库操作失败、服务异常 |

### 业务状态码 (Business Status Codes)

| RetCode | 常量名 | 说明 | 使用场景 |
|---------|--------|------|----------|
| 1 | `RetCodeActive` | 启用/成功状态 | AI配置启用、评论类型普通、缓存成功等 |
| 2 | `RetCodeInactive` | 禁用/进行中状态 | AI配置禁用、评论类型行级、缓存中等 |
| 3 | `RetCodeCacheSuccess` | 缓存成功 | Gitlab项目同步成功 |
| 4 | `RetCodeInvalid` | 无效状态 | AI评论人工评分无效 |

## 迁移指南

### 从旧代码迁移到新系统

#### 1. 替换直接数字
```go
// ❌ 旧代码
response.Success(c, data, "操作成功", 0)
response.Error(c, err, "用户未登录", 10000)

// ✅ 新代码
response.Success(c, data, "操作成功", int(constants.RetCodeSuccess))
response.Error(c, err, "用户未登录", int(constants.RetCodeNotLogIn))
```

#### 2. 替换硬编码常量
```go
// ❌ 旧代码
const RetCodeSuccess = 0
const RetCodeNotLogIn = 10000

// ✅ 新代码
import "code-review-go/internal/pkg/constants"
// 直接使用 constants.RetCodeSuccess
```

#### 3. 更新导入
```go
// 在需要使用的文件中添加导入
import "code-review-go/internal/pkg/constants"
```

## 新系统的优势

### 1. 统一管理
- 所有RetCode集中在一个包中
- 避免重复定义和冲突
- 便于维护和更新

### 2. 类型安全
- RetCode是强类型，避免误用
- 编译时检查，减少运行时错误

### 3. 功能丰富
- 内置验证方法（`IsSuccess()`, `IsSystemError()`等）
- 分类信息（`GetCategory()`）
- 字符串表示（`String()`）

### 4. 易于扩展
- 新增RetCode只需在一个地方定义
- 自动继承所有验证和分类方法

## 使用示例

### 基本使用
```go
import (
    "code-review-go/internal/pkg/constants"
    "code-review-go/internal/pkg/response"
)

// 成功响应
response.Success(c, data, "操作成功", int(constants.RetCodeSuccess))

// 错误响应
response.Error(c, err, "参数错误", int(constants.RetCodeInvalidParams))
response.Error(c, err, "服务器错误", int(constants.RetCodeInternalError))
```

### 高级使用
```go
code := constants.RetCodeNotLogIn

// 验证和分类
if code.IsSystemError() {
    log.Error("系统错误:", code.String())
}

category := code.GetCategory()
log.Info("错误分类:", category)
```

## 维护说明

- **新增RetCode**: 在 `internal/pkg/constants/ret_code.go` 中添加
- **更新文档**: 同步更新本文档和 `internal/pkg/constants/README.md`
- **向后兼容**: 保持现有RetCode值不变
- **统一使用**: 新代码必须使用新的常量系统

## 注意事项

1. **避免循环导入**: 不要在constants包中导入其他业务包
2. **类型转换**: 使用RetCode常量时需要转换为int类型
3. **统一管理**: 所有RetCode都应该在constants包中定义
4. **文档同步**: 更新常量时同步更新相关文档

## 统计信息

- **总RetCode数量**: 25个
- **成功码**: 1个
- **系统错误码**: 16个
- **HTTP状态码**: 2个
- **业务状态码**: 6个

## 相关文件

- **常量定义**: `internal/pkg/constants/ret_code.go`
- **使用说明**: `internal/pkg/constants/README.md`
- **包概览**: `internal/pkg/constants/constants.go` 