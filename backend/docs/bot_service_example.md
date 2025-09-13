# 机器人服务使用示例

## 概述

机器人服务提供了多个内置的专业代码审查机器人，每个机器人都有特定的专业领域和提示词。用户可以通过机器人名称和AI模型来调用相应的机器人进行代码审查。

## 内置机器人角色

### 1. 安全审查机器人
- **Key**: security_reviewer
- **名称**: 安全审查机器人
- **分类**: security
- **专业领域**: 代码安全漏洞检测
- **重点关注**: SQL注入、XSS攻击、CSRF攻击、权限绕过、敏感信息泄露等

### 2. 性能优化机器人
- **Key**: performance_optimizer
- **名称**: 性能优化机器人
- **分类**: performance
- **专业领域**: 代码性能分析和优化建议
- **重点关注**: 算法复杂度、数据库查询优化、内存使用、并发处理等

### 3. 代码规范机器人
- **Key**: code_style_checker
- **名称**: 代码规范机器人
- **分类**: style
- **专业领域**: 代码规范和最佳实践检查
- **重点关注**: 命名规范、代码结构、注释完整性、函数设计等

### 4. 架构设计机器人
- **Key**: architecture_reviewer
- **名称**: 架构设计机器人
- **分类**: architecture
- **专业领域**: 代码架构和设计模式审查
- **重点关注**: 设计模式、架构层次、模块耦合、接口设计等

### 5. 测试质量机器人
- **Key**: test_quality_analyzer
- **名称**: 测试质量机器人
- **分类**: testing
- **专业领域**: 测试覆盖率和测试质量审查
- **重点关注**: 测试覆盖率、测试用例设计、边界条件测试等

### 6. 文档质量机器人
- **Key**: documentation_reviewer
- **名称**: 文档质量机器人
- **分类**: documentation
- **专业领域**: 代码文档和注释质量审查
- **重点关注**: API文档、函数注释、复杂逻辑注释、README文档等

## API接口

### 1. 机器人代码审查

**接口地址**: `POST /v1/ai/check/bot`

**请求参数**:
```json
{
    "bot_name": "security_reviewer",
    "type": "OpenAI",
    "model": "gpt-4",
    "api_url": "https://api.openai.com/v1",
    "code_content": "function login(username, password) { return db.query('SELECT * FROM users WHERE username=' + username + ' AND password=' + password); }",
    "additional_prompt": "请特别关注SQL注入问题"
}
```

**响应示例**:
```json
{
    "code": 200,
    "message": "机器人审查完成",
    "data": {
        "bot_name": "安全审查机器人",
        "bot_category": "安全",
        "ai_model": "gpt-4",
        "result": "🔴 **高危问题**：SQL注入漏洞\n- 问题位置：第1行\n- 问题描述：直接拼接用户输入到SQL查询中\n- 潜在风险：攻击者可以通过输入恶意SQL代码获取数据库信息\n- 修复建议：使用参数化查询或预处理语句",
        "create_time": 1703001600
    }
}
```

### 2. 获取所有机器人角色

**接口地址**: `GET /v1/ai/check/bot/roles`

**响应示例**:
```json
{
    "code": 200,
    "message": "获取机器人角色列表成功",
    "data": {
        "roles": [
            {
                "name": "安全审查机器人",
                "description": "专门负责代码安全漏洞检测",
                "category": "安全"
            },
            {
                "name": "性能优化机器人",
                "description": "专门负责代码性能分析和优化建议",
                "category": "性能"
            }
        ]
    }
}
```

### 3. 根据分类获取机器人角色

**接口地址**: `GET /v1/ai/check/bot/roles/{category}`

**示例**: `GET /v1/ai/check/bot/roles/安全`

**响应示例**:
```json
{
    "code": 200,
    "message": "获取机器人角色列表成功",
    "data": {
        "roles": [
            {
                "name": "安全审查机器人",
                "description": "专门负责代码安全漏洞检测",
                "category": "安全"
            }
        ]
    }
}
```

### 4. 获取机器人角色详情

**接口地址**: `GET /v1/ai/check/bot/roles/detail/{bot_name}`

**示例**: `GET /v1/ai/check/bot/roles/detail/安全审查机器人`

**响应示例**:
```json
{
    "code": 200,
    "message": "获取机器人角色详情成功",
    "data": {
        "name": "安全审查机器人",
        "description": "专门负责代码安全漏洞检测",
        "category": "安全",
        "prompt": "你是一位专业的代码安全审查专家。请仔细审查以下代码，重点关注以下安全方面：\n\n1. **SQL注入漏洞**：检查数据库查询语句是否使用了参数化查询\n2. **XSS攻击**：检查用户输入是否进行了适当的转义和过滤\n..."
    }
}
```

## 使用场景

### 1. 多维度代码审查
```javascript
// 使用多个机器人进行全方位审查
const bots = [
    "安全审查机器人",
    "性能优化机器人", 
    "代码规范机器人"
];

for (const bot of bots) {
    const result = await botReview({
        bot_name: bot,
        ai_model_id: 1,
        code_content: codeToReview
    });
    console.log(`${bot} 审查结果:`, result);
}
```

### 2. 特定领域深度审查
```javascript
// 专门进行安全审查
const securityResult = await botReview({
    bot_name: "security_reviewer",
    type: "OpenAI",
    model: "gpt-4",
    api_url: "https://api.openai.com/v1",
    code_content: userInputCode,
    additional_prompt: "请特别关注用户输入验证和权限控制"
});
```

### 3. 分类审查
```javascript
// 获取所有安全相关的机器人
const securityBots = await getBotRolesByCategory("security");

// 使用所有安全机器人进行审查
for (const bot of securityBots) {
    await botReview({
        bot_name: bot.name,
        type: "OpenAI",
        model: "gpt-4",
        api_url: "https://api.openai.com/v1",
        code_content: codeToReview
    });
}
```

## 优势

1. **专业化**: 每个机器人都有特定的专业领域，提供更精准的审查
2. **标准化**: 内置的提示词确保审查结果的一致性和专业性
3. **灵活性**: 支持额外提示词，可以根据具体需求定制审查重点
4. **可扩展**: 可以轻松添加新的机器人角色
5. **易用性**: 通过简单的API调用即可使用专业级的代码审查

## 注意事项

1. 机器人名称必须与内置角色完全匹配
2. AI类型、模型和API地址必须正确配置
3. 代码内容应该包含完整的上下文信息
4. 额外提示词是可选的，用于进一步定制审查重点
5. 审查结果的质量取决于AI模型的能力和代码的复杂度
6. 确保API地址可访问且具有相应的权限