# AI 服务

基于 RAG 的 AI 服务，支持多模型、多角色智能体和 gRPC 接口。

## 功能特性

### 🤖 多 AI 模型支持
- **OpenAI**: GPT-3.5/GPT-4 系列模型
- **通义千问**: 阿里云 DashScope 模型
- **本地模型**: 支持 Hugging Face Transformers 模型

### 🎭 多角色智能体
- **代码审查专家**: 专业的代码质量分析和安全检测
- **系统架构师**: 架构设计和技术选型建议
- **安全专家**: 安全漏洞检测和防护建议
- **测试工程师**: 测试用例设计和质量评估
- **通用助手**: 各种编程和技术支持

### 🔍 RAG 增强检索
- **向量存储**: 支持 FAISS 和 ChromaDB
- **智能检索**: 基于语义相似度的文档搜索
- **上下文增强**: 自动构建相关代码上下文

### 🌐 gRPC 服务
- **高性能**: 基于 gRPC 的高性能服务接口
- **多语言支持**: 支持多种编程语言客户端
- **反射支持**: 支持 gRPC 反射，便于调试

### 💾 数据持久化
- **MySQL**: 存储服务数据、会话记录等
- **向量数据库**: 存储文档嵌入向量
- **缓存**: Redis 支持（可选）

## 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │    │   Web Client    │    │  Mobile Client  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      gRPC Server          │
                    │   (AI Service Gateway)    │
                    └─────────────┬─────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
┌─────────▼─────────┐    ┌────────▼────────┐    ┌────────▼────────┐
│  Agent Manager    │    │  Model Manager  │    │  RAG Service    │
│  (多角色智能体)     │    │  (多模型管理)    │    │  (增强检索)      │
└─────────┬─────────┘    └────────┬────────┘    └────────┬────────┘
          │                       │                       │
          └───────────────────────┼───────────────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │     Vector Store         │
                    │  (FAISS/ChromaDB)        │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │        MySQL             │
                    │     (数据持久化)          │
                    └───────────────────────────┘
```

## 快速开始

### 环境要求

- Python 3.9+
- MySQL 8.0+
- Docker & Docker Compose（可选）

### 1. 克隆项目

```bash
git clone <repository-url>
cd ai-service
```

### 2. 安装依赖

```bash
# 使用 Poetry（推荐）
poetry install

# 或使用 pip（如果没有 Poetry）
pip install -e .
```

### 3. 配置服务

#### 方式一：使用环境变量配置（推荐）

```bash
# 复制环境变量模板
cp env.example ../.env

# 编辑环境变量文件
vim ../.env
```

**注意**: 系统现在完全使用环境变量配置，不再依赖 JSON 配置文件。

### 4. 初始化数据库

```bash
# 创建数据库表
python -c "from services.database.connection import db_manager; import asyncio; asyncio.run(db_manager.create_tables())"
```

### 5. 启动服务

```bash
# 使用启动脚本（推荐）
python run.py testing      # 测试环境
python run.py production   # 生产环境

# 或直接运行
python main.py

# 或使用 Poetry
poetry run python main.py
```

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f ai-service

# 停止服务
docker-compose down
```

## API 接口

### gRPC 服务

服务地址: `localhost:50051`

#### 代码审查

```protobuf
rpc CodeReview(CodeReviewRequest) returns (CodeReviewResponse);
```

#### 智能聊天

```protobuf
rpc Chat(ChatRequest) returns (ChatResponse);
```

#### 文档搜索

```protobuf
rpc DocumentSearch(DocumentSearchRequest) returns (DocumentSearchResponse);
```

#### 健康检查

```protobuf
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
```

### 客户端示例

#### Python 客户端

```python
import grpc
from proto import ai_service_pb2
from proto import ai_service_pb2_grpc

# 创建连接
channel = grpc.insecure_channel('localhost:50051')
stub = ai_service_pb2_grpc.AIServiceStub(channel)

# 代码审查
request = ai_service_pb2.CodeReviewRequest(
    git_url="https://github.com/user/repo.git",
    branch="main",
    diff_content="diff --git a/file.py b/file.py...",
    agent_role="code_reviewer"
)

response = stub.CodeReview(request)
print(response.review)
```

#### Go 客户端

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    pb "your-project/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewAIServiceClient(conn)
    
    req := &pb.CodeReviewRequest{
        GitUrl: "https://github.com/user/repo.git",
        Branch: "main",
        DiffContent: "diff --git a/file.py b/file.py...",
        AgentRole: "code_reviewer",
    }
    
    resp, err := client.CodeReview(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println(resp.Review)
}
```

## 配置说明

### 环境变量配置

系统支持两种配置方式，环境变量配置优先级更高：

#### 环境变量列表

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `SERVICE_HOST` | 服务监听地址 | 0.0.0.0 |
| `SERVICE_PORT` | 服务端口 | 50051 |
| `DATABASE_HOST` | 数据库主机 | 127.0.0.1 |
| `DATABASE_PORT` | 数据库端口 | 3306 |
| `DATABASE_USERNAME` | 数据库用户名 | root |
| `DATABASE_PASSWORD` | 数据库密码 | 123456 |
| `DATABASE_NAME` | 数据库名称 | xxx_review |
| `AI_PROVIDER` | AI 提供商 | dashscope |
| `OPENAI_API_KEY` | OpenAI API 密钥 | - |
| `DASHSCOPE_API_KEY` | 通义千问 API 密钥 | - |
| `VECTOR_STORE_TYPE` | 向量存储类型 | faiss |
| `VECTOR_STORE_PATH` | 向量存储路径 | ./data/vector_store |
| `GITLAB_TOKEN` | GitLab 访问令牌 | - |
| `LOG_LEVEL` | 日志级别 | INFO |

#### 配置说明

系统完全使用环境变量配置，所有配置项都从项目根目录的 `.env` 文件中读取。

### 智能体配置

可以在 `services/ai/agent_manager.py` 中自定义智能体角色：

```python
AgentRole(
    role_key="custom_expert",
    name="自定义专家",
    description="特定领域的专家",
    system_prompt="你是一位...",
    model_provider="openai",
    temperature=0.7
)
```

## 开发指南

### 项目结构

```
ai-service/
├── proto/                    # gRPC 协议定义
├── services/                 # 核心服务层
│   ├── ai/                  # AI 模型服务
│   ├── database/            # 数据库服务
│   ├── rag/                 # RAG 检索服务
│   └── grpc/                # gRPC 服务实现
├── models/                  # 数据模型
├── config/                  # 配置管理
├── utils/                   # 工具函数
├── scripts/                 # 脚本文件
├── main.py                  # 服务入口
├── pyproject.toml           # 依赖管理
├── Dockerfile               # Docker 配置
└── docker-compose.yml       # 容器编排
```

### 添加新的 AI 模型

1. 在 `services/ai/model_manager.py` 中继承 `BaseAIModel`
2. 实现 `generate` 和 `chat` 方法
3. 在 `AIModelManager` 中注册新模型

### 添加新的智能体角色

1. 在 `services/ai/agent_manager.py` 中定义新的 `AgentRole`
2. 配置系统提示词和能力描述
3. 在 `AgentManager` 中注册新角色

### 添加新的向量存储后端

1. 在 `services/rag/vector_store.py` 中继承 `BaseVectorStore`
2. 实现所有抽象方法
3. 在 `VectorStoreManager` 中支持新后端

## 监控和日志

### 健康检查

```bash
# gRPC 健康检查
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

### 日志查看

```bash
# Docker 环境
docker-compose logs -f ai-service

# 本地环境
tail -f ai_service.log
```

## 性能优化

### 1. 模型缓存
- 本地模型会缓存在内存中
- 向量嵌入会缓存到磁盘

### 2. 数据库优化
- 使用连接池管理数据库连接
- 合理设置索引

### 3. 向量搜索优化
- 调整 chunk_size 和 chunk_overlap
- 使用合适的相似度阈值

## 故障排除

### 常见问题

1. **gRPC 连接失败**
   - 检查端口是否被占用
   - 确认防火墙设置

2. **AI 模型加载失败**
   - 检查 API 密钥配置
   - 确认网络连接

3. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接参数

4. **向量存储错误**
   - 检查存储路径权限
   - 确认磁盘空间充足

## 许可证

MIT License

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 联系方式

- 作者: richLpf
- 邮箱: 1045465391@qq.com