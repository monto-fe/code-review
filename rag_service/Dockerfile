FROM python:3.9-slim as builder

WORKDIR /app

# 设置pip镜像源和优化参数
ENV PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple/
ENV PIP_TRUSTED_HOST=pypi.tuna.tsinghua.edu.cn
ENV PIP_NO_CACHE_DIR=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1


# 安装构建依赖 - 兼容CentOS和Ubuntu
RUN if command -v apt-get >/dev/null 2>&1; then \
        # Ubuntu/Debian系统 - 使用国内镜像源
        sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
        sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
        apt-get update && apt-get install -y \
            build-essential \
            && apt-get clean && rm -rf /var/lib/apt/lists/* || true; \
    elif command -v yum >/dev/null 2>&1; then \
        # CentOS/RHEL系统
        yum update -y && yum install -y \
            gcc \
            gcc-c++ \
            make \
            && yum clean all; \
    elif command -v dnf >/dev/null 2>&1; then \
        # Fedora/CentOS 8+系统
        dnf update -y && dnf install -y \
            gcc \
            gcc-c++ \
            make \
            && dnf clean all; \
    else \
        echo "Unsupported package manager"; exit 1; \
    fi

# 先复制pyproject.toml以利用Docker缓存
COPY pyproject.toml ./

# 安装Python依赖（限制线程数）
RUN pip install --upgrade pip setuptools wheel --no-cache-dir --disable-pip-version-check && \
    pip install -e . --no-cache-dir --disable-pip-version-check

# 生产阶段
FROM python:3.9-slim

WORKDIR /app

# 安装运行时依赖 - 兼容CentOS和Ubuntu
RUN if command -v apt-get >/dev/null 2>&1; then \
        # Ubuntu/Debian系统 - 使用国内镜像源
        sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
        sed -i 's/security.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
        apt-get update && apt-get install -y \
            git \
            curl \
            && apt-get clean && rm -rf /var/lib/apt/lists/* || true; \
    elif command -v yum >/dev/null 2>&1; then \
        # CentOS/RHEL系统
        yum update -y && yum install -y \
            git \
            curl \
            && yum clean all; \
    elif command -v dnf >/dev/null 2>&1; then \
        # Fedora/CentOS 8+系统
        dnf update -y && dnf install -y \
            git \
            curl \
            && dnf clean all; \
    else \
        echo "Unsupported package manager"; exit 1; \
    fi

# 从构建阶段复制Python包到系统Python目录
COPY --from=builder /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

# 创建必要的目录
RUN mkdir -p /app/data/vector_store /tmp/git_repos

# 复制应用代码
COPY . .

# 设置环境变量
ENV PYTHONPATH=/app
ENV OPENAI_API_KEY=""
ENV RAG_SERVICE_HOST="0.0.0.0"
ENV RAG_SERVICE_PORT="8000"
ENV VECTOR_STORE_TYPE="faiss"
ENV VECTOR_STORE_PATH="/app/data/vector_store"
ENV GIT_TEMP_DIR="/tmp/git_repos"
ENV GITLAB_TOKEN=""
ENV SKIP_DEPS_CHECK="true"

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

# uvicorn app:app --host 0.0.0.0 --port 8000 --reload