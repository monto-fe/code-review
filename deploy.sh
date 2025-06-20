#!/bin/bash

# 代码审查系统部署脚本
# 包含前端、后端和RAG服务的完整部署

set -e

echo "=== 代码审查系统部署脚本 ==="

# 检查Docker和Docker Compose是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 检查.env文件是否存在
if [ ! -f .env ]; then
    echo "警告: .env文件不存在，将创建默认配置"
    cat > .env << EOF
# 基础配置
IP=localhost
API_PORT=9000
CONSOLE_PORT=3000
NODE_ENV=development

# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_DATABASE=code_review

# GitLab配置
GITLAB_TOKEN=your_gitlab_token_here

# RAG服务配置
RAG_SERVICE_URL=http://localhost:8000
RAG_SERVICE_HOST=0.0.0.0
RAG_SERVICE_PORT=8000
VECTOR_STORE_TYPE=faiss
VECTOR_STORE_PATH=/app/data/vector_store
GIT_TEMP_DIR=/tmp/git_repos

# AI服务配置
OPENAI_API_KEY=your_openai_api_key_here
DEEPSEEK_API_KEY=your_deepseek_api_key_here
UCLOUD_API_KEY=your_ucloud_api_key_here
EOF
    echo "已创建.env文件，请根据实际情况修改配置"
    echo "特别注意修改以下配置："
    echo "1. GITLAB_TOKEN: 你的GitLab访问令牌"
    echo "2. 数据库配置: 根据你的数据库设置修改"
    echo "3. AI服务配置: 根据需要配置相应的API密钥"
    echo ""
    read -p "请修改.env文件后按回车继续..."
fi

# 检查必要的目录是否存在
echo "检查项目结构..."
if [ ! -d "frontend" ]; then
    echo "错误: frontend目录不存在"
    exit 1
fi

if [ ! -d "backend" ]; then
    echo "错误: backend目录不存在"
    exit 1
fi

if [ ! -d "rag_service" ]; then
    echo "错误: rag_service目录不存在"
    exit 1
fi

# 构建和启动服务
echo "开始构建和启动服务..."

# 停止现有服务
echo "停止现有服务..."
docker-compose down

# 清理旧镜像（可选）
if [ "$1" = "--clean" ]; then
    echo "清理旧镜像..."
    docker-compose down --rmi all --volumes --remove-orphans
fi

# 构建镜像
echo "构建Docker镜像..."
docker-compose build --no-cache

# 启动服务
echo "启动服务..."
docker-compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 30

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

# 检查RAG服务健康状态
echo "检查RAG服务健康状态..."
if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo "✅ RAG服务运行正常"
else
    echo "❌ RAG服务可能未正常启动，请检查日志"
    docker-compose logs rag-service
fi

# 检查后端服务
echo "检查后端服务..."
if curl -f http://localhost:9000/health > /dev/null 2>&1; then
    echo "✅ 后端服务运行正常"
else
    echo "❌ 后端服务可能未正常启动，请检查日志"
    docker-compose logs backend
fi

# 检查前端服务
echo "检查前端服务..."
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ 前端服务运行正常"
else
    echo "❌ 前端服务可能未正常启动，请检查日志"
    docker-compose logs frontend
fi

echo ""
echo "=== 部署完成 ==="
echo "服务访问地址："
echo "前端界面: http://localhost:3000"
echo "后端API: http://localhost:9000"
echo "RAG服务: http://localhost:8000"
echo ""
echo "查看服务日志："
echo "docker-compose logs -f [service_name]"
echo ""
echo "停止服务："
echo "docker-compose down"
echo ""
echo "重启服务："
echo "docker-compose restart [service_name]" 