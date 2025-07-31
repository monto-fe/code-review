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

# 安装 curl 用于健康检查
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY pyproject.toml ./

# 使用最保守的方式安装依赖
RUN python -m pip install --upgrade "pip<24.0" --no-cache-dir --disable-pip-version-check --quiet

# 逐个安装依赖，允许依赖解析
RUN python -m pip install fastapi --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install uvicorn --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install python-dotenv --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install gitpython --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install langchain --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install langchain-community --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install sentence-transformers --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install faiss-cpu --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install torch==2.1.2 --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install torchvision==0.16.2 --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install pydantic --no-cache-dir --disable-pip-version-check --quiet
RUN python -m pip install huggingface-hub --no-cache-dir --disable-pip-version-check --quiet

# 复制应用代码
COPY . .

# 设置上海时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

# 暴露端口
EXPOSE 8000

# 运行命令
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]