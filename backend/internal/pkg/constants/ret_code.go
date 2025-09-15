package constants

// RetCode 返回码类型定义
type RetCode int

// 成功码 (Success Codes)
const (
	RetCodeSuccess RetCode = 0 // 成功
)

// 系统级错误码 (System Error Codes) - 10000+
const (
	RetCodeNotLogIn                 RetCode = 10000 // 未登录
	RetCodeInvalidParams            RetCode = 10001 // 无效请求参数
	RetCodeInvalidCredentials       RetCode = 10002 // 无效用户名或密码
	RetCodeLoginFailed              RetCode = 10003 // 登录失败
	RetCodeGetUserListFailed        RetCode = 10004 // 获取用户列表失败
	RetCodeGetUserInfoFailed        RetCode = 10005 // 获取用户信息失败
	RetCodeOperatorMissing          RetCode = 10006 // 操作者信息缺失
	RetCodeInvalidUserID            RetCode = 10007 // 无效用户ID格式
	RetCodeCheckUserExistenceFailed RetCode = 10009 // 检查用户存在性失败
	RetCodeUserAlreadyExists        RetCode = 10010 // 用户已存在
	RetCodeUserOperationFailed      RetCode = 10011 // 用户操作失败
	RetCodeTokenMissing             RetCode = 10012 // Token缺失
	RetCodeDeleteUserFailed         RetCode = 10013 // 删除用户失败
	RetCodeInvalidToken             RetCode = 10014 // 无效Token
	RetCodeTokenExpired             RetCode = 10016 // Token过期
)

// HTTP状态码 (HTTP Status Codes)
const (
	RetCodeBadRequest    RetCode = 400 // 请求参数错误
	RetCodeInternalError RetCode = 500 // 服务器内部错误
)

// 业务状态码 (Business Status Codes)
const (
	RetCodeActive       RetCode = 1 // 启用/成功状态
	RetCodeInactive     RetCode = 2 // 禁用/进行中状态
	RetCodeCacheSuccess RetCode = 3 // 缓存成功
	RetCodeInvalid      RetCode = 4 // 无效状态
)

// String 返回RetCode的字符串表示
func (rc RetCode) String() string {
	switch rc {
	case RetCodeSuccess:
		return "成功"
	case RetCodeNotLogIn:
		return "未登录"
	case RetCodeInvalidParams:
		return "无效请求参数"
	case RetCodeInvalidCredentials:
		return "无效用户名或密码"
	case RetCodeLoginFailed:
		return "登录失败"
	case RetCodeGetUserListFailed:
		return "获取用户列表失败"
	case RetCodeGetUserInfoFailed:
		return "获取用户信息失败"
	case RetCodeOperatorMissing:
		return "操作者信息缺失"
	case RetCodeInvalidUserID:
		return "无效用户ID格式"
	case RetCodeCheckUserExistenceFailed:
		return "检查用户存在性失败"
	case RetCodeUserAlreadyExists:
		return "用户已存在"
	case RetCodeUserOperationFailed:
		return "用户操作失败"
	case RetCodeTokenMissing:
		return "Token缺失"
	case RetCodeDeleteUserFailed:
		return "删除用户失败"
	case RetCodeInvalidToken:
		return "无效Token"
	case RetCodeTokenExpired:
		return "Token过期"
	case RetCodeBadRequest:
		return "请求参数错误"
	case RetCodeInternalError:
		return "服务器内部错误"
	case RetCodeActive:
		return "启用/成功状态"
	case RetCodeInactive:
		return "禁用/进行中状态"
	case RetCodeCacheSuccess:
		return "缓存成功"
	case RetCodeInvalid:
		return "无效状态"
	default:
		return "未知返回码"
	}
}

// IsSuccess 判断是否为成功码
func (rc RetCode) IsSuccess() bool {
	return rc == RetCodeSuccess
}

// IsError 判断是否为错误码
func (rc RetCode) IsError() bool {
	return rc != RetCodeSuccess
}

// IsSystemError 判断是否为系统级错误码
func (rc RetCode) IsSystemError() bool {
	return rc >= 10000
}

// IsHTTPError 判断是否为HTTP错误码
func (rc RetCode) IsHTTPError() bool {
	return rc == RetCodeBadRequest || rc == RetCodeInternalError
}

// IsBusinessStatus 判断是否为业务状态码
func (rc RetCode) IsBusinessStatus() bool {
	return rc >= 1 && rc <= 4
}

// GetCategory 获取RetCode的分类
func (rc RetCode) GetCategory() string {
	switch {
	case rc == RetCodeSuccess:
		return "成功码"
	case rc.IsSystemError():
		return "系统级错误码"
	case rc.IsHTTPError():
		return "HTTP状态码"
	case rc.IsBusinessStatus():
		return "业务状态码"
	default:
		return "未知分类"
	}
}
