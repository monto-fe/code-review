# Merge Review 任务管理

## 功能描述
这个模块用于管理AI Code Review的异步任务，支持对特定的merge进行效果评估，包括任务列表、创建新任务和查看任务详情。

## 页面结构

### 1. 任务列表页面 (`/aicodecheck/merge-review`)
- **功能**：展示所有创建的异步任务
- **特性**：
  - 任务列表表格展示
  - 任务状态管理（待处理、处理中、已完成、失败）
  - 支持查看详情、删除任务等操作
  - 新建任务按钮跳转

### 2. 新建任务页面 (`/aicodecheck/merge-review-create`)
- **功能**：创建新的Merge Review检测任务
- **布局**：左右分栏布局
  - **左侧**：任务配置表单
    - Merge链接输入（支持自动解析项目ID和Merge ID）
    - 项目ID输入
    - Merge ID输入
    - AI模型选择
  - **右侧**：机器人选择和输出
    - 机器人选择区域（支持多选）
    - 输出内容展示区域
    - 最终总结输出

### 3. 任务详情页面 (`/aicodecheck/merge-review-detail/:id`)
- **功能**：查看已完成的检测任务结果
- **特性**：
  - 任务基本信息展示（只读）
  - 检测结果详细展示
  - 支持刷新数据
  - 返回列表功能

## 机器人类型
- **代码规范检测机器人**：检测代码风格、命名规范、注释等
- **逻辑错误检测机器人**：检测代码逻辑错误、边界条件等
- **安全漏洞检测机器人**：检测SQL注入、XSS等安全漏洞
- **性能优化检测机器人**：检测性能瓶颈、内存泄漏等

## 技术实现
- 使用项目通用的CommonTable组件
- 支持响应式布局和滚动
- 集成国际化支持（中英文）
- 使用MobX进行状态管理
- 模块化CSS样式

## 路由配置
- 列表页：`/aicodecheck/merge-review`
- 新建页：`/aicodecheck/merge-review-create`
- 详情页：`/aicodecheck/merge-review-detail/:id`

## API接口
- `GET /ai/check/history` - 获取任务列表
- `POST /ai/check/create` - 创建新任务
- `GET /ai/check/history/:id` - 获取任务详情
- `DELETE /ai/check/history/:id` - 删除任务

## 权限要求
- 需要admin角色访问
- 所有页面都受权限控制保护 