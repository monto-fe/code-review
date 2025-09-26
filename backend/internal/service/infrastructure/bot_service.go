package infrastructure

import (
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"fmt"
	"strings"
	"time"
)

// BotRole 机器人角色定义
type BotRole struct {
	Name        string `json:"name"`        // 机器人名称
	Description string `json:"description"` // 机器人描述
	Prompt      string `json:"prompt"`      // 机器人提示词
	Category    string `json:"category"`    // 机器人分类
}

// BotService 机器人服务
type BotService struct {
	roles map[string]BotRole // 内置机器人角色
}

// NewBotService 创建机器人服务实例
func NewBotService() *BotService {
	service := &BotService{
		roles: make(map[string]BotRole),
	}
	service.initBuiltinRoles()
	return service
}

// initBuiltinRoles 初始化内置机器人角色
func (s *BotService) initBuiltinRoles() {
	// 安全审查机器人
	s.roles["security_reviewer"] = BotRole{
		Name:        "安全审查机器人",
		Description: "专门负责代码安全漏洞检测",
		Category:    "security",
		Prompt: `你是一位专业的代码安全审查专家。请仔细审查以下代码，重点关注以下安全方面：

1. **SQL注入漏洞**：检查数据库查询语句是否使用了参数化查询
2. **XSS攻击**：检查用户输入是否进行了适当的转义和过滤
3. **CSRF攻击**：检查是否有适当的CSRF保护机制
4. **权限绕过**：检查访问控制逻辑是否正确
5. **敏感信息泄露**：检查是否有硬编码的密码、密钥等敏感信息
6. **输入验证**：检查用户输入是否进行了充分的验证和过滤
7. **文件上传安全**：检查文件上传功能的安全性
8. **加密使用**：检查敏感数据的加密存储和传输

请按照以下格式提供审查结果：
- 🔴 **高危问题**：[具体问题描述]
- 🟡 **中危问题**：[具体问题描述]
- 🟢 **建议改进**：[具体建议]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 问题描述
3. 潜在风险
4. 修复建议`,
	}

	// 性能优化机器人
	s.roles["performance_optimizer"] = BotRole{
		Name:        "性能优化机器人",
		Description: "专门负责代码性能分析和优化建议",
		Category:    "performance",
		Prompt: `你是一位专业的性能优化专家。请仔细分析以下代码的性能问题，重点关注：

1. **算法复杂度**：检查算法的时间复杂度和空间复杂度
2. **数据库查询优化**：检查SQL查询是否高效，是否有N+1查询问题
3. **内存使用**：检查是否有内存泄漏、过度内存分配等问题
4. **并发处理**：检查并发安全性和性能
5. **缓存策略**：检查是否有适当的缓存机制
6. **资源管理**：检查文件、网络连接等资源的正确释放
7. **循环优化**：检查循环逻辑是否可以优化
8. **数据结构选择**：检查是否选择了合适的数据结构

请按照以下格式提供分析结果：
- 🚀 **性能瓶颈**：[具体问题描述]
- ⚡ **优化建议**：[具体建议]
- 📊 **性能指标**：[相关性能指标分析]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 性能影响分析
3. 优化方案
4. 预期性能提升`,
	}

	// 代码规范机器人
	s.roles["code_style_checker"] = BotRole{
		Name:        "代码规范机器人",
		Description: "专门负责代码规范和最佳实践检查",
		Category:    "style",
		Prompt: `你是一位专业的代码规范审查专家。请仔细检查以下代码是否符合编程规范和最佳实践：

1. **命名规范**：检查变量、函数、类名是否符合命名约定
2. **代码结构**：检查代码组织结构是否清晰合理
3. **注释完整性**：检查关键逻辑是否有适当的注释
4. **函数设计**：检查函数是否单一职责，参数是否合理
5. **错误处理**：检查是否有适当的错误处理机制
6. **代码重复**：检查是否有重复代码可以提取
7. **依赖管理**：检查依赖关系是否合理
8. **可读性**：检查代码是否易于理解和维护

请按照以下格式提供审查结果：
- 📝 **规范问题**：[具体问题描述]
- 💡 **改进建议**：[具体建议]
- ✅ **最佳实践**：[相关最佳实践建议]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 规范要求
3. 改进方案
4. 代码示例（如适用）`,
	}

	// 架构设计机器人
	s.roles["architecture_reviewer"] = BotRole{
		Name:        "架构设计机器人",
		Description: "专门负责代码架构和设计模式审查",
		Category:    "architecture",
		Prompt: `你是一位专业的软件架构师。请仔细审查以下代码的架构设计和设计模式使用：

1. **设计模式**：检查是否合理使用了设计模式
2. **架构层次**：检查代码分层是否清晰合理
3. **模块耦合**：检查模块间的耦合度是否过高
4. **接口设计**：检查接口设计是否合理
5. **依赖注入**：检查依赖管理是否合理
6. **单一职责**：检查类和模块是否遵循单一职责原则
7. **开闭原则**：检查代码是否易于扩展
8. **可测试性**：检查代码是否易于单元测试

请按照以下格式提供审查结果：
- 🏗️ **架构问题**：[具体问题描述]
- 🔧 **设计改进**：[具体建议]
- 📐 **架构建议**：[相关架构建议]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 架构原则分析
3. 改进方案
4. 设计模式建议`,
	}

	// 测试质量机器人
	s.roles["test_quality_analyzer"] = BotRole{
		Name:        "测试质量机器人",
		Description: "专门负责测试覆盖率和测试质量审查",
		Category:    "testing",
		Prompt: `你是一位专业的测试质量专家。请仔细审查以下代码的测试相关方面：

1. **测试覆盖率**：检查关键逻辑是否有测试覆盖
2. **测试用例设计**：检查测试用例是否全面
3. **边界条件测试**：检查是否测试了边界条件
4. **异常情况测试**：检查是否测试了异常情况
5. **集成测试**：检查是否有适当的集成测试
6. **性能测试**：检查是否有性能测试
7. **测试可维护性**：检查测试代码是否易于维护
8. **测试数据管理**：检查测试数据是否合理

请按照以下格式提供审查结果：
- 🧪 **测试问题**：[具体问题描述]
- 📋 **测试建议**：[具体建议]
- ✅ **质量保证**：[相关质量保证建议]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 测试需求分析
3. 测试用例建议
4. 测试策略建议`,
	}

	// 文档质量机器人
	s.roles["documentation_reviewer"] = BotRole{
		Name:        "文档质量机器人",
		Description: "专门负责代码文档和注释质量审查",
		Category:    "documentation",
		Prompt: `你是一位专业的文档质量专家。请仔细审查以下代码的文档和注释质量：

1. **API文档**：检查API是否有完整的文档说明
2. **函数注释**：检查函数是否有清晰的参数和返回值说明
3. **复杂逻辑注释**：检查复杂逻辑是否有适当的注释
4. **README文档**：检查项目是否有完整的README文档
5. **变更日志**：检查是否有变更日志记录
6. **示例代码**：检查是否有使用示例
7. **文档一致性**：检查文档与代码是否一致
8. **文档可读性**：检查文档是否易于理解

请按照以下格式提供审查结果：
- 📚 **文档问题**：[具体问题描述]
- ✍️ **文档建议**：[具体建议]
- 📖 **文档标准**：[相关文档标准建议]

对于每个问题，请提供：
1. 问题位置（文件名和行号）
2. 文档需求分析
3. 文档改进建议
4. 文档模板建议`,
	}
}

// GetBotRole 根据机器人名称获取机器人角色
func (s *BotService) GetBotRole(botName string) (*BotRole, error) {
	role, exists := s.roles[botName]
	if !exists {
		return nil, fmt.Errorf("机器人角色 '%s' 不存在", botName)
	}
	return &role, nil
}

// GetAllBotRoles 获取所有内置机器人角色
func (s *BotService) GetAllBotRoles() []BotRole {
	roles := make([]BotRole, 0, len(s.roles))
	for _, role := range s.roles {
		roles = append(roles, role)
	}
	return roles
}

// GetBotRolesByCategory 根据分类获取机器人角色
func (s *BotService) GetBotRolesByCategory(category string) []BotRole {
	var roles []BotRole
	for _, role := range s.roles {
		if strings.EqualFold(role.Category, category) {
			roles = append(roles, role)
		}
	}
	return roles
}

// ExecuteBotReview 执行机器人审查
func (s *BotService) ExecuteBotReview(botName string, aiType string, aiModel string, apiURL string, codeContent string, additionalPrompt string) (*dto.BotReviewResult, error) {
	// 获取机器人角色
	role, err := s.GetBotRole(botName)
	if err != nil {
		return nil, err
	}

	// 构建AI配置信息
	aiConfig := &model.AIConfig{
		Type:   aiType,
		Model:  aiModel,
		APIURL: apiURL,
	}

	// 构建完整的提示词
	fullPrompt := s.buildFullPrompt(role.Prompt, additionalPrompt, codeContent)

	// 调用AI进行审查
	reviewResult, err := s.callAIForReview(aiConfig, fullPrompt)
	if err != nil {
		return nil, fmt.Errorf("AI审查失败: %w", err)
	}

	return &dto.BotReviewResult{
		BotName:     role.Name,
		BotCategory: role.Category,
		AIModel:     aiModel,
		Result:      reviewResult,
		CreateTime:  time.Now().Unix(),
	}, nil
}

// buildFullPrompt 构建完整的提示词
func (s *BotService) buildFullPrompt(basePrompt, additionalPrompt, codeContent string) string {
	var prompt strings.Builder

	prompt.WriteString(basePrompt)
	prompt.WriteString("\n\n")

	if additionalPrompt != "" {
		prompt.WriteString("## 额外要求\n")
		prompt.WriteString(additionalPrompt)
		prompt.WriteString("\n\n")
	}

	prompt.WriteString("## 待审查代码\n")
	prompt.WriteString("```\n")
	prompt.WriteString(codeContent)
	prompt.WriteString("\n```\n\n")

	prompt.WriteString("请根据上述要求对代码进行详细审查。")

	return prompt.String()
}

// callAIForReview 调用AI进行审查
func (s *BotService) callAIForReview(aiConfig *model.AIConfig, prompt string) (string, error) {
	// 这里应该调用实际的AI服务
	// 暂时返回模拟结果
	return fmt.Sprintf("AI审查结果（使用模型：%s）\n\n%s", aiConfig.Model, prompt), nil
}

// GetBotService 获取机器人服务实例
func GetBotService() *BotService {
	return NewBotService()
}
