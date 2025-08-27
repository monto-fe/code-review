package utils

import (
	"fmt"
	"regexp"
)

// 正则表达式常量
var (
	// 邮箱格式验证
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// 手机号格式验证（中国大陆）
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

	// 用户名格式验证（字母开头，允许字母数字下划线，长度3-20）
	usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,19}$`)

	// 密码强度验证（至少8位，包含字母和数字）
	passwordRegex = regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$`)

	// 命名空间格式验证（字母开头，允许字母数字下划线连字符，长度2-50）
	namespaceRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{1,49}$`)
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if email == "" {
		return nil // 允许空值
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	if len(email) > 100 {
		return fmt.Errorf("邮箱长度不能超过100个字符")
	}

	return nil
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // 允许空值
	}

	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("手机号格式不正确，请输入11位中国大陆手机号")
	}

	return nil
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("用户名格式不正确，必须以字母开头，只能包含字母、数字、下划线，长度3-20位")
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于8位")
	}

	if len(password) > 50 {
		return fmt.Errorf("密码长度不能超过50位")
	}

	if !passwordRegex.MatchString(password) {
		return fmt.Errorf("密码必须包含字母和数字")
	}

	return nil
}

// ValidateNamespace 验证命名空间格式
func ValidateNamespace(namespace string) error {
	if namespace == "" {
		return fmt.Errorf("命名空间不能为空")
	}

	if !namespaceRegex.MatchString(namespace) {
		return fmt.Errorf("命名空间格式不正确，必须以字母开头，只能包含字母、数字、下划线、连字符，长度2-50位")
	}

	return nil
}

// ValidateName 验证姓名格式
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("姓名不能为空")
	}

	if len(name) < 2 {
		return fmt.Errorf("姓名长度不能少于2个字符")
	}

	if len(name) > 50 {
		return fmt.Errorf("姓名长度不能超过50个字符")
	}

	// 检查是否包含特殊字符（除了中文、英文、数字、空格）
	specialCharRegex := regexp.MustCompile(`[^\p{Han}a-zA-Z0-9\s]`)
	if specialCharRegex.MatchString(name) {
		return fmt.Errorf("姓名不能包含特殊字符")
	}

	return nil
}

// ValidateJob 验证职位格式
func ValidateJob(job string) error {
	if job == "" {
		return nil // 允许空值
	}

	if len(job) > 100 {
		return fmt.Errorf("职位长度不能超过100个字符")
	}

	return nil
}

// ValidateRoleIDs 验证角色ID数组
func ValidateRoleIDs(roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil // 允许空数组
	}

	// 检查是否有重复的ID
	seen := make(map[uint]bool)
	for _, id := range roleIDs {
		if id == 0 {
			return fmt.Errorf("角色ID不能为0")
		}
		if seen[id] {
			return fmt.Errorf("角色ID不能重复")
		}
		seen[id] = true
	}

	return nil
}
