# 手动代码审核接口使用示例

## 接口地址
`POST /v1/ai/check/manual`

## 请求参数

### 新的请求格式（支持多个机器人配置）

```json
{
    "merge_url": "https://gitlab.example.com/project/repo/-/merge_requests/123",
    "ai_model_id": 1,
    "bot_configs": [
        {
            "bot_name": "安全审查机器人",
            "bot_prompt": "请重点检查代码中的安全漏洞，包括SQL注入、XSS攻击、权限绕过等问题"
        },
        {
            "bot_name": "性能优化机器人", 
            "bot_prompt": "请分析代码性能问题，关注算法复杂度、内存使用、数据库查询优化等"
        },
        {
            "bot_name": "代码规范机器人",
            "bot_prompt": "请检查代码是否符合编程规范，包括命名规范、注释完整性、代码结构等"
        }
    ]
}
```

### 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| merge_url | string | 是 | 合并请求的完整链接 |
| ai_model_id | uint | 是 | AI模型的ID |
| bot_configs | array | 否 | 机器人配置列表 |

### BotConfig 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bot_name | string | 是 | 检测机器人名称 |
| bot_prompt | string | 是 | 机器人提示词 |

## 响应示例

```json
{
    "code": 200,
    "message": "审核任务已创建",
    "data": {
        "task_id": 123,
        "status": "processing",
        "estimated_time": 30
    }
}
```

## 获取审核结果

### 接口地址
`GET /v1/ai/check/result/{task_id}`

### 响应示例

```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "id": 123,
        "project_id": 456,
        "project_name": "my-project",
        "merge_id": 789,
        "merge_title": "Feature: Add new functionality",
        "merge_url": "https://gitlab.example.com/project/repo/-/merge_requests/123",
        "status": 2,
        "result": "审核完成，发现3个问题...",
        "ai_model": "gpt-4",
        "bot_configs": [
            {
                "bot_name": "安全审查机器人",
                "bot_prompt": "请重点检查代码中的安全漏洞..."
            },
            {
                "bot_name": "性能优化机器人",
                "bot_prompt": "请分析代码性能问题..."
            }
        ],
        "create_time": 1703001600,
        "update_time": 1703001630
    }
}
```

## 处理流程

1. **快速创建任务**: 系统根据merge URL和AI模型ID快速创建审核任务
2. **异步信息补全**: 后台异步从缓存中获取项目信息和AI模型信息
3. **信息更新**: 将获取到的项目名称、合并请求标题、AI模型名称等信息更新到任务中
4. **执行审核**: 使用配置的机器人和提示词执行代码审核

## 使用场景

1. **多维度代码审查**: 可以配置多个专门的机器人，每个机器人负责不同方面的审查
2. **自定义审查重点**: 根据项目特点，为每个机器人设置不同的提示词
3. **灵活配置**: 可以根据需要添加或减少机器人数量
4. **一对一关系**: 每个机器人都有对应的专用提示词，确保审查的针对性
5. **快速响应**: 先创建任务返回响应，再异步补全详细信息，提高用户体验

## 注意事项

1. `bot_configs` 为可选参数，如果不提供，将使用默认的审查配置
2. 每个机器人配置中的 `bot_name` 和 `bot_prompt` 都是必填的
3. 机器人配置会以JSON格式存储在数据库中
4. 建议为每个机器人设置明确的职责和提示词，避免重复审查
5. 项目信息和AI模型信息会在任务创建后异步更新，初始响应中可能不包含完整信息