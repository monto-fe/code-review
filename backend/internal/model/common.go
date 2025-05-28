package model

// TableName 表名常量
const (
	TableUser           = "t_user"
	TableRole           = "t_role"
	TableToken          = "t_login_token"
	TableUserRole       = "t_user_role"
	TableRolePermission = "t_role_permission"
	TableResource       = "t_resource"
	TableGitlab         = "t_gitlab_info"
	TableAIMessage      = "t_ai_message"
	TableAIConfig       = "t_ai_config"
	TableCommonRule     = "t_common_rule"
	TableNamespace      = "t_namespace"
	TableCustomRule     = "t_custom_rule"
)

// PageModel 分页模型
type PageModel struct {
	Current  int `json:"current" form:"current"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusInit    UserStatus = iota // 初始状态
	UserStatusEnable                    // 启用
	UserStatusDisable                   // 禁用
)

// Language 编程语言枚举
type Language string

const (
	LanguagePython     Language = "Python"
	LanguageJava       Language = "Java"
	LanguageJavaScript Language = "JavaScript"
	LanguageGolang     Language = "Golang"
	LanguageRuby       Language = "Ruby"
	LanguageCpp        Language = "C++"
	LanguageOther      Language = "other"
)
