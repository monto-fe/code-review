# 使用 Python 3.9 slim 镜像
FROM python:3.9-slim

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on \
    PIP_QUIET=1 \
    PIP_PROGRESS_BAR=off \
    PYTHONHASHSEED=0

# 配置国内镜像源
ENV PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple/
ENV PIP_TRUSTED_HOST=pypi.tuna.tsinghua.edu.cn

# 安装系统依赖
RUN apt-get update && \
    apt-get install -y \
        curl \
        git \
        build-essential \
        gcc \
        g++ \
        && apt-get clean \
        && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# 复制依赖文件
COPY pyproject.toml ./

# 安装 Poetry
RUN pip install poetry

# 配置 Poetry
RUN poetry config virtualenvs.create false

# 安装 Python 依赖
RUN poetry install --no-dev --no-interaction --no-ansi

# 复制应用代码
COPY . .

# 生成 gRPC 协议文件
RUN python scripts/generate_proto.py

# 创建必要的目录
RUN mkdir -p /app/data/vector_store /tmp/git_repos /app/models

# 设置上海时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

# 创建非 root 用户
RUN useradd -m -u 1000 aiuser && \
    chown -R aiuser:aiuser /app
USER aiuser

# 暴露端口
EXPOSE 50051

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:50051/health || exit 1

# 运行命令
CMD ["python", "main.py"]