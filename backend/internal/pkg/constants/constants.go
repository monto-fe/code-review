package constants

// 包概览
//
// constants包用于统一管理项目中的所有常量定义，包括但不限于：
// - 返回码 (RetCode)
// - 业务状态码
// - 系统配置常量
// - 枚举值定义
//
// 使用方式：
// import "code-review-go/internal/pkg/constants"
//
// 示例：
// response.Success(c, data, "操作成功", int(constants.RetCodeSuccess))
// response.Error(c, err, "参数错误", int(constants.RetCodeInvalidParams))
