# Constants 包使用说明

## 概述

`constants` 包用于统一管理项目中的所有常量定义，包括返回码、业务状态码、系统配置常量等。

## 包结构

```
internal/pkg/constants/
├── constants.go      # 包概览和说明
├── ret_code.go      # 返回码常量定义
└── README.md        # 本文件
```

## RetCode 使用说明

### 1. 导入包

```go
import "code-review-go/internal/pkg/constants"
```

### 2. 成功响应

```go
import (
    "code-review-go/internal/pkg/constants"
    "code-review-go/internal/pkg/response"
)

// 成功响应
response.Success(c, data, "操作成功", int(constants.RetCodeSuccess))
```

### 3. 错误响应

```go
// 系统级错误
response.Error(c, err, "用户未登录", int(constants.RetCodeNotLogIn))
response.Error(c, err, "参数错误", int(constants.RetCodeInvalidParams))

// HTTP状态码错误
response.Error(c, err, "请求参数错误", int(constants.RetCodeBadRequest))
response.Error(c, err, "服务器内部错误", int(constants.RetCodeInternalError))

// 业务逻辑错误
response.Error(c, err, "用户已存在", int(constants.RetCodeUserAlreadyExists))
```

### 4. 业务状态码

```go
// 在业务逻辑中使用
if status == int(constants.RetCodeActive) {
    // 处理启用状态
}

if status == int(constants.RetCodeInactive) {
    // 处理禁用状态
}
```

### 5. RetCode 验证方法

```go
code := constants.RetCodeNotLogIn

// 检查是否为成功码
if code.IsSuccess() {
    // 处理成功情况
}

// 检查是否为系统错误
if code.IsSystemError() {
    // 处理系统错误
}

// 检查是否为HTTP错误
if code.IsHTTPError() {
    // 处理HTTP错误
}

// 检查是否为业务状态码
if code.IsBusinessStatus() {
    // 处理业务状态
}

// 获取分类信息
category := code.GetCategory()

// 获取字符串表示
description := code.String()
```

## 返回码分类

### 成功码
- `RetCodeSuccess (0)`: 成功

### 系统级错误码 (10000+)
- `RetCodeNotLogIn (10000)`: 未登录
- `RetCodeInvalidParams (10001)`: 无效请求参数
- `RetCodeInvalidCredentials (10002)`: 无效用户名或密码
- `RetCodeLoginFailed (10003)`: 登录失败
- `RetCodeGetUserListFailed (10004)`: 获取用户列表失败
- `RetCodeGetUserInfoFailed (10005)`: 获取用户信息失败
- `RetCodeOperatorMissing (10006)`: 操作者信息缺失
- `RetCodeInvalidUserID (10007)`: 无效用户ID格式
- `RetCodeCheckUserExistenceFailed (10009)`: 检查用户存在性失败
- `RetCodeUserAlreadyExists (10010)`: 用户已存在
- `RetCodeUserOperationFailed (10011)`: 用户操作失败
- `RetCodeTokenMissing (10012)`: Token缺失
- `RetCodeDeleteUserFailed (10013)`: 删除用户失败
- `RetCodeInvalidToken (10014)`: 无效Token
- `RetCodeTokenExpired (10016)`: Token过期

### HTTP状态码
- `RetCodeBadRequest (400)`: 请求参数错误
- `RetCodeInternalError (500)`: 服务器内部错误

### 业务状态码
- `RetCodeActive (1)`: 启用/成功状态
- `RetCodeInactive (2)`: 禁用/进行中状态
- `RetCodeCacheSuccess (3)`: 缓存成功
- `RetCodeInvalid (4)`: 无效状态

## 最佳实践

### 1. 统一使用常量
```go
// ✅ 推荐：使用常量
response.Error(c, err, "参数错误", int(constants.RetCodeInvalidParams))

// ❌ 不推荐：直接使用数字
response.Error(c, err, "参数错误", 10001)
```

### 2. 类型安全
```go
// RetCode 是强类型，可以避免误用
var code constants.RetCode = constants.RetCodeSuccess
```

### 3. 验证和分类
```go
// 使用内置方法验证和分类
if code.IsSystemError() {
    // 记录系统错误日志
    log.Error("系统错误:", code.String())
}
```

### 4. 扩展性
```go
// 新增返回码时，只需要在 ret_code.go 中添加常量定义
// 并更新 String() 方法和其他相关方法
```

## 迁移指南

### 从旧代码迁移

#### 1. 替换直接数字
```go
// 旧代码
response.Success(c, data, "成功", 0)
response.Error(c, err, "错误", 10000)

// 新代码
response.Success(c, data, "成功", int(constants.RetCodeSuccess))
response.Error(c, err, "错误", int(constants.RetCodeNotLogIn))
```

#### 2. 替换硬编码常量
```go
// 旧代码
const RetCodeSuccess = 0
const RetCodeNotLogIn = 10000

// 新代码
import "code-review-go/internal/pkg/constants"
// 直接使用 constants.RetCodeSuccess
```

#### 3. 更新导入
```go
// 在需要使用的文件中添加导入
import "code-review-go/internal/pkg/constants"
```

## 注意事项

1. **避免循环导入**: 不要在 constants 包中导入其他业务包
2. **类型转换**: 使用 RetCode 常量时需要转换为 int 类型
3. **向后兼容**: 新增常量时保持现有常量的值不变
4. **文档同步**: 更新常量时同步更新相关文档

## 维护说明

- 新增 RetCode 时，请在 `ret_code.go` 中添加常量定义
- 同时更新 `String()` 方法和其他相关方法
- 更新本 README 文档
- 考虑是否需要添加新的验证方法 