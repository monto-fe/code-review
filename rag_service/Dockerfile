FROM python:3.9-slim

WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    git \
    curl \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p /app/data/vector_store /tmp/git_repos

# 复制依赖文件并安装
COPY requirements.txt .

# 升级pip并安装依赖
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# 复制应用代码
COPY . .

# 设置环境变量
# ENV PYTHONPATH=/app
# ENV OPENAI_API_KEY=""
# ENV RAG_SERVICE_HOST="0.0.0.0"
# ENV RAG_SERVICE_PORT="8000"
# ENV VECTOR_STORE_TYPE="faiss"
# ENV VECTOR_STORE_PATH="/app/data/vector_store"
# ENV GIT_TEMP_DIR="/tmp/git_repos"
# ENV GITLAB_TOKEN=""

# 创建非root用户
RUN useradd -m -u 1000 raguser && \
    chown -R raguser:raguser /app /tmp/git_repos

# 切换到非root用户
USER raguser

# 暴露端口
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# 启动应用
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "1"]
