# 机器人服务与AI管理器集成示例

## 概述

机器人服务现在支持直接使用AI管理器的参数（Type、Model、APIURL）来调用机器人进行代码审查，无需依赖预配置的AI模型ID。

## 新的请求参数结构

### BotReviewRequest
```go
type BotReviewRequest struct {
    BotName          string `json:"bot_name"`           // 机器人名称，必填
    Type             string `json:"type"`               // AI类型，必填
    Model            string `json:"model"`              // AI模型，必填
    APIURL           string `json:"api_url"`            // API地址，必填
    CodeContent      string `json:"code_content"`       // 代码内容，必填
    AdditionalPrompt string `json:"additional_prompt"`  // 额外提示词，可选
}
```

## 使用示例

### 1. 使用OpenAI GPT-4进行安全审查

```json
{
    "bot_name": "security_reviewer",
    "type": "OpenAI",
    "model": "gpt-4",
    "api_url": "https://api.openai.com/v1",
    "code_content": "function login(username, password) {\n    const query = `SELECT * FROM users WHERE username='${username}' AND password='${password}'`;\n    return db.query(query);\n}",
    "additional_prompt": "请特别关注SQL注入和XSS攻击"
}
```

### 2. 使用DeepSeek进行性能优化审查

```json
{
    "bot_name": "performance_optimizer",
    "type": "DeepSeek",
    "model": "deepseek-coder",
    "api_url": "https://api.deepseek.com/v1",
    "code_content": "function processData(items) {\n    for (let i = 0; i < items.length; i++) {\n        for (let j = 0; j < items.length; j++) {\n            if (items[i].id === items[j].id) {\n                // 处理逻辑\n            }\n        }\n    }\n}",
    "additional_prompt": "请重点关注算法复杂度和循环优化"
}
```

### 3. 使用Qwen进行代码规范审查

```json
{
    "bot_name": "code_style_checker",
    "type": "Qwen",
    "model": "qwen-turbo",
    "api_url": "https://dashscope.aliyuncs.com/api/v1",
    "code_content": "function getUserData(userId) {\n    var data = null;\n    if (userId != null) {\n        data = fetch('/api/users/' + userId);\n    }\n    return data;\n}",
    "additional_prompt": "请检查命名规范和代码风格"
}
```

### 4. 使用UCloud AI进行架构设计审查

```json
{
    "bot_name": "architecture_reviewer",
    "type": "UCloud",
    "model": "ucloud-ai",
    "api_url": "https://api.ucloud.cn/ai",
    "code_content": "class UserService {\n    constructor() {\n        this.db = new Database();\n        this.cache = new Cache();\n        this.email = new EmailService();\n    }\n    \n    async createUser(userData) {\n        // 直接操作数据库\n        const user = await this.db.create(userData);\n        // 发送邮件\n        await this.email.send(user.email);\n        return user;\n    }\n}",
    "additional_prompt": "请分析类的职责和依赖关系"
}
```

## 支持的AI类型

### 1. OpenAI
- **Type**: "OpenAI"
- **常用模型**: "gpt-4", "gpt-3.5-turbo", "gpt-4-turbo"
- **API URL**: "https://api.openai.com/v1"

### 2. DeepSeek
- **Type**: "DeepSeek"
- **常用模型**: "deepseek-coder", "deepseek-chat"
- **API URL**: "https://api.deepseek.com/v1"

### 3. Qwen (通义千问)
- **Type**: "Qwen"
- **常用模型**: "qwen-turbo", "qwen-plus", "qwen-max"
- **API URL**: "https://dashscope.aliyuncs.com/api/v1"

### 4. UCloud AI
- **Type**: "UCloud"
- **常用模型**: "ucloud-ai"
- **API URL**: "https://api.ucloud.cn/ai"

## 响应格式

所有请求都返回统一的响应格式：

```json
{
    "code": 200,
    "message": "机器人审查完成",
    "data": {
        "bot_name": "安全审查机器人",
        "bot_category": "安全",
        "ai_model": "gpt-4",
        "result": "🔴 **高危问题**：SQL注入漏洞\n- 问题位置：第2行\n- 问题描述：使用模板字符串直接拼接用户输入到SQL查询中\n- 潜在风险：攻击者可以通过输入恶意SQL代码获取数据库信息\n- 修复建议：使用参数化查询或预处理语句\n\n🟡 **中危问题**：密码明文传输\n- 问题位置：第1行\n- 问题描述：密码参数直接传递给数据库查询\n- 潜在风险：密码可能以明文形式存储在日志中\n- 修复建议：对密码进行哈希处理后再存储",
        "create_time": 1703001600
    }
}
```

## 最佳实践

### 1. 选择合适的机器人
- **安全审查**: 使用"安全审查机器人"检查安全漏洞
- **性能优化**: 使用"性能优化机器人"分析性能问题
- **代码规范**: 使用"代码规范机器人"检查编码规范
- **架构设计**: 使用"架构设计机器人"分析架构问题
- **测试质量**: 使用"测试质量机器人"检查测试覆盖
- **文档质量**: 使用"文档质量机器人"检查文档完整性

### 2. 配置合适的AI模型
- **复杂代码审查**: 使用GPT-4或Qwen-Max等高级模型
- **简单代码检查**: 使用GPT-3.5-turbo或Qwen-Turbo等快速模型
- **代码生成**: 使用DeepSeek-Coder等专门的代码模型

### 3. 编写有效的额外提示词
```javascript
// 好的额外提示词示例
const additionalPrompts = {
    security: "请特别关注SQL注入、XSS攻击和权限绕过问题",
    performance: "请重点关注算法复杂度和数据库查询优化",
    style: "请检查命名规范、代码格式和注释完整性",
    architecture: "请分析类的职责分离和依赖注入"
};
```

### 4. 错误处理
```javascript
try {
    const result = await botReview({
        bot_name: "安全审查机器人",
        type: "OpenAI",
        model: "gpt-4",
        api_url: "https://api.openai.com/v1",
        code_content: codeToReview
    });
    console.log("审查结果:", result.data.result);
} catch (error) {
    console.error("审查失败:", error.message);
}
```

## 优势

1. **灵活性**: 支持多种AI提供商和模型
2. **实时配置**: 无需预配置AI模型，可以动态指定
3. **成本控制**: 可以根据需求选择不同价位的模型
4. **专业化**: 每个机器人都有专业的提示词和审查重点
5. **可扩展**: 支持添加新的AI提供商和机器人角色