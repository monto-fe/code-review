# 手动代码审核功能

## 功能概述

手动代码审核功能允许用户通过输入项目ID和合并请求ID，手动触发AI代码审核，无需依赖GitLab webhook。

## 功能特性

- ✅ 手动触发代码审核
- ✅ 异步执行审核任务
- ✅ 审核历史查询
- ✅ 审核结果详情查看
- ✅ 支持多种筛选条件
- ✅ 用户权限控制

## API接口

### 1. 手动触发代码审核

**接口**: `POST /v1/ai/check/manual`

**请求头**:
```
jwt_token: <JWT_TOKEN>
Content-Type: application/json
```

**请求体**:
```json
{
    "project_id": 123,
    "merge_id": 456,
    "ai_model": "gpt-4",  // 可选，默认使用配置的模型
    "rule_id": 1          // 可选，默认使用项目配置的规则
}
```

**响应**:
```json
{
    "ret_code": 0,
    "message": "审核任务已创建",
    "data": {
        "task_id": 789,
        "status": "processing",
        "estimated_time": 30
    }
}
```

### 2. 查询审核历史

**接口**: `GET /v1/ai/check/history`

**请求头**:
```
jwt_token: <JWT_TOKEN>
```

**查询参数**:
- `current`: 当前页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20)
- `status`: 状态筛选 (0-全部, 1-进行中, 2-完成, 3-失败)
- `start_date`: 开始日期 (时间戳)
- `end_date`: 结束日期 (时间戳)

**响应**:
```json
{
    "ret_code": 0,
    "message": "查询成功",
    "data": {
        "list": [
            {
                "id": 789,
                "project_id": 123,
                "project_name": "frontend-app",
                "merge_id": 456,
                "merge_title": "Add new feature",
                "status": 2,
                "status_text": "完成",
                "create_time": 1703001600,
                "update_time": 1703001630
            }
        ],
        "total": 100
    }
}
```

### 3. 获取审核结果详情

**接口**: `GET /v1/ai/check/result/{id}`

**请求头**:
```
jwt_token: <JWT_TOKEN>
```

**路径参数**:
- `id`: 任务ID

**响应**:
```json
{
    "ret_code": 0,
    "message": "查询成功",
    "data": {
        "id": 789,
        "project_id": 123,
        "project_name": "frontend-app",
        "merge_id": 456,
        "merge_title": "Add new feature",
        "merge_url": "https://gitlab.com/project/merge/456",
        "status": 2,
        "result": "代码审核结果...",
        "ai_model": "gpt-4",
        "rule_name": "代码规范检查",
        "create_time": 1703001600,
        "update_time": 1703001630
    }
}
```

## 数据库表结构

### t_manual_check_task 表

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键 |
| user_id | BIGINT | 触发用户ID |
| project_id | INT | 项目ID |
| merge_id | INT | 合并请求ID |
| project_name | VARCHAR(200) | 项目名称 |
| merge_title | VARCHAR(500) | 合并请求标题 |
| merge_url | VARCHAR(500) | 合并请求URL |
| status | TINYINT | 状态: 1-进行中, 2-完成, 3-失败 |
| result | TEXT | 审核结果 |
| error_message | VARCHAR(500) | 错误信息 |
| ai_model | VARCHAR(50) | AI模型 |
| rule_id | BIGINT | 规则ID |
| rule_name | VARCHAR(100) | 规则名称 |
| create_time | BIGINT | 创建时间 |
| update_time | BIGINT | 更新时间 |

## 使用示例

### 使用curl测试

```bash
# 1. 手动触发代码审核
curl -X POST "http://localhost:9000/v1/ai/check/manual" \
  -H "Content-Type: application/json" \
  -H "jwt_token: YOUR_JWT_TOKEN" \
  -d '{
    "project_id": 123,
    "merge_id": 456,
    "ai_model": "gpt-4"
  }'

# 2. 查询审核历史
curl -X GET "http://localhost:9000/v1/ai/check/history?current=1&page_size=10" \
  -H "jwt_token: YOUR_JWT_TOKEN"

# 3. 获取审核结果详情
curl -X GET "http://localhost:9000/v1/ai/check/result/789" \
  -H "jwt_token: YOUR_JWT_TOKEN"
```

### 使用测试脚本

```bash
# 给脚本执行权限
chmod +x test_manual_check.sh

# 运行测试
./test_manual_check.sh
```

## 业务流程

1. **任务创建**: 用户提交审核请求，系统验证参数并创建任务记录
2. **异步执行**: 系统异步执行AI审核，不阻塞用户请求
3. **状态更新**: 审核完成后更新任务状态和结果
4. **结果查询**: 用户可查询审核历史和结果详情

## 权限控制

- 用户只能查看自己创建的审核任务
- 需要有效的JWT Token进行身份验证
- 支持项目级别的权限控制（待实现）

## 错误处理

- 参数验证错误
- 项目不存在或无权限访问
- 合并请求不存在或状态不正确
- AI服务异常
- 数据库操作失败

## 后续优化

- [ ] 支持批量审核
- [ ] 支持定时审核
- [ ] 支持审核模板
- [ ] 支持自定义审核规则
- [ ] 实时进度显示
- [ ] 审核结果推送
- [ ] 移动端适配 