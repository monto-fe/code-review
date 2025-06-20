#!/bin/bash

echo "=== 启动RAG服务 ==="

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 启动服务
echo "启动RAG服务..."
docker-compose up -d rag-service

# 等待服务启动
echo "等待服务启动..."
sleep 5

# 检查服务状态
echo "检查服务状态..."
docker-compose ps rag-service

echo "=== 服务启动完成 ==="
echo "服务地址: http://localhost:8000"
echo "API文档: http://localhost:8000/docs"
echo "健康检查: http://localhost:8000/health" 