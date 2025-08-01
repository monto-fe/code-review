# 使用 Python 3.9 slim 镜像
FROM python:3.9-slim

# 设置工作目录
WORKDIR /app

# 设置环境变量（禁用所有可能导致线程的功能）
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

# 安装系统依赖：curl 用于健康检查，git 用于 GitPython
RUN apt-get update && apt-get install -y curl git && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY pyproject.toml ./

# 使用最保守的方式安装依赖
RUN python -m pip install --upgrade "pip<24.0" --no-cache-dir --disable-pip-version-check --quiet

# 安装基础依赖
RUN python -m pip install fastapi uvicorn python-dotenv pydantic requests numpy --no-cache-dir --disable-pip-version-check --quiet

# 安装 GitPython
RUN python -m pip install gitpython --no-cache-dir --disable-pip-version-check --quiet

# 安装 HuggingFace 相关依赖
RUN python -m pip install huggingface-hub transformers --no-cache-dir --disable-pip-version-check --quiet

# 安装 PyTorch
RUN python -m pip install torch==2.1.2 torchvision==0.16.2 --no-cache-dir --disable-pip-version-check --quiet

# 安装 sentence-transformers
RUN python -m pip install sentence-transformers --no-cache-dir --disable-pip-version-check --quiet

# 安装 FAISS
RUN python -m pip install faiss-cpu --no-cache-dir --disable-pip-version-check --quiet

# 安装 LangChain 相关依赖
RUN python -m pip install langchain==0.0.350 langchain-community==0.0.10 --no-cache-dir --disable-pip-version-check --quiet

# 复制应用代码
COPY . .

# 创建必要的目录
RUN mkdir -p /app/data /tmp/git_repos

# 设置上海时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

# 暴露端口
EXPOSE 8000

# 运行命令
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]