# AI 服务 Makefile

.PHONY: help install dev-install test lint format clean build run docker-build docker-run docker-stop proto

# 默认目标
help:
	@echo "AI 服务管理命令:"
	@echo "  install      - 安装生产依赖"
	@echo "  dev-install  - 安装开发依赖"
	@echo "  test         - 运行测试"
	@echo "  lint         - 代码检查"
	@echo "  format       - 代码格式化"
	@echo "  clean        - 清理临时文件"
	@echo "  build        - 构建项目"
	@echo "  run          - 运行服务"
	@echo "  proto        - 生成 gRPC 协议文件"
	@echo "  docker-build - 构建 Docker 镜像"
	@echo "  docker-run   - 运行 Docker 容器"
	@echo "  docker-stop  - 停止 Docker 容器"

# 安装生产依赖
install:
	poetry install --no-dev

# 安装开发依赖
dev-install:
	poetry install

# 运行测试
test:
	poetry run pytest tests/ -v

# 代码检查
lint:
	poetry run flake8 .
	poetry run mypy .

# 代码格式化
format:
	poetry run black .
	poetry run isort .

# 清理临时文件
clean:
	find . -type f -name "*.pyc" -delete
	find . -type d -name "__pycache__" -delete
	find . -type d -name "*.egg-info" -exec rm -rf {} +
	rm -rf build/
	rm -rf dist/
	rm -rf .pytest_cache/
	rm -rf .mypy_cache/

# 构建项目
build:
	poetry build

# 运行服务
run:
	poetry run python main.py

# 生成 gRPC 协议文件
proto:
	poetry run python scripts/generate_proto.py

# 构建 Docker 镜像
docker-build:
	docker build -t ai-service:latest .

# 运行 Docker 容器
docker-run:
	docker-compose up -d

# 停止 Docker 容器
docker-stop:
	docker-compose down

# 查看日志
logs:
	docker-compose logs -f ai-service

# 进入容器
shell:
	docker-compose exec ai-service /bin/bash

# 数据库迁移
migrate:
	poetry run alembic upgrade head

# 创建数据库迁移
migration:
	poetry run alembic revision --autogenerate -m "$(message)"

# 初始化数据库
init-db:
	poetry run python -c "from services.database.connection import db_manager; import asyncio; asyncio.run(db_manager.create_tables())"