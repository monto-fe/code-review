#!/bin/bash

echo "=== 重新构建并启动RAG服务 ==="

# 停止现有容器
echo "1. 停止现有容器..."
docker-compose down rag-service

# 清理镜像（可选，取消注释以强制重新构建）
# echo "2. 清理旧镜像..."
# docker rmi $(docker images -q code-review-rag-service) 2>/dev/null || true

# 重新构建镜像
echo "2. 重新构建镜像..."
docker-compose build rag-service

# 启动服务
echo "3. 启动服务..."
docker-compose up -d rag-service

# 等待服务启动
echo "4. 等待服务启动..."
sleep 10

# 检查服务状态
echo "5. 检查服务状态..."
docker-compose ps rag-service

# 检查健康状态
echo "6. 检查健康状态..."
curl -f http://localhost:8000/health || echo "服务还未完全启动，请稍等..."

echo "=== 完成 ==="
echo "服务地址: http://localhost:8000"
echo "API文档: http://localhost:8000/docs" 