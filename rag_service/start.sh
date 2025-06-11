#!/bin/bash

# 创建并激活虚拟环境
python3.9 -m venv venv
source venv/bin/activate

# 安装依赖
pip install -r requirements.txt

# 创建必要的目录
mkdir -p data/vector_store
mkdir -p /tmp/git_repos

# 启动服务
uvicorn app:app --host 0.0.0.0 --port 8000 