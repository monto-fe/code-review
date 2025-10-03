"""
配置管理器 - 统一管理配置传递
"""
import os
from typing import Dict, Any, Optional

# 全局配置存储
_global_config = None

def load_env_config():
    """从环境变量加载配置"""
    return {
        "service": {
            "host": os.getenv("SERVICE_HOST", "0.0.0.0"),
            "port": int(os.getenv("SERVICE_PORT", "50051")),
            "max_workers": int(os.getenv("SERVICE_MAX_WORKERS", "10")),
            "max_message_length": int(os.getenv("SERVICE_MAX_MESSAGE_LENGTH", "4194304")),
            "enable_reflection": os.getenv("SERVICE_ENABLE_REFLECTION", "true").lower() == "true"
        },
        "database": {
            "host": os.getenv("DATABASE_HOST", "127.0.0.1"),
            "port": int(os.getenv("DATABASE_PORT", "3306")),
            "username": os.getenv("DATABASE_USERNAME", "root"),
            "password": os.getenv("DATABASE_PASSWORD", "123456"),
            "database": os.getenv("DATABASE_NAME", "xxx_review"),
            "charset": os.getenv("DATABASE_CHARSET", "utf8mb4"),
            "pool_size": int(os.getenv("DATABASE_POOL_SIZE", "10")),
            "max_overflow": int(os.getenv("DATABASE_MAX_OVERFLOW", "20"))
        },
        "ai": {
            "provider": os.getenv("AI_PROVIDER", "dashscope"),
            "openai": {
                "api_key": os.getenv("OPENAI_API_KEY", ""),
                "base_url": os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
                "model": os.getenv("OPENAI_MODEL", "gpt-3.5-turbo")
            },
            "dashscope": {
                "api_key": os.getenv("DASHSCOPE_API_KEY", "sk-xxxx"),
                "model": os.getenv("DASHSCOPE_MODEL", "qwen-turbo")
            },
            "local": {
                "model_path": os.getenv("LOCAL_MODEL_PATH", "./data/models"),
                "model_name": os.getenv("LOCAL_MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2")
            },
            "vector_store": {
                "type": os.getenv("VECTOR_STORE_TYPE", "faiss"),
                "path": os.getenv("VECTOR_STORE_PATH", "./data/vector_store"),
                "embedding_model": os.getenv("VECTOR_STORE_EMBEDDING_MODEL", "sentence-transformers/all-MiniLM-L6-v2")
            },
            "parameters": {
                "max_tokens": int(os.getenv("AI_MAX_TOKENS", "2048")),
                "temperature": float(os.getenv("AI_TEMPERATURE", "0.7")),
                "top_p": float(os.getenv("AI_TOP_P", "0.9"))
            }
        },
        "rag": {
            "chunk_size": int(os.getenv("RAG_CHUNK_SIZE", "1000")),
            "chunk_overlap": int(os.getenv("RAG_CHUNK_OVERLAP", "200")),
            "similarity_threshold": float(os.getenv("RAG_SIMILARITY_THRESHOLD", "0.7")),
            "max_context_length": int(os.getenv("RAG_MAX_CONTEXT_LENGTH", "4000")),
            "enable_rerank": os.getenv("RAG_ENABLE_RERANK", "true").lower() == "true"
        },
        "git": {
            "token": os.getenv("GITLAB_TOKEN", ""),
            "temp_dir": os.getenv("GIT_TEMP_DIR", "/tmp/git_repos")
        },
        "logging": {
            "level": os.getenv("LOG_LEVEL", "INFO"),
            "file": os.getenv("LOG_FILE", "ai_service.log")
        }
    }

def set_config(config):
    """设置全局配置"""
    global _global_config
    _global_config = config

def get_config():
    """获取全局配置"""
    global _global_config
    if _global_config is None:
        raise RuntimeError("配置未初始化，请先调用 set_config()")
    return _global_config

def get_database_config():
    """获取数据库配置"""
    config = get_config()
    db_conf = config.database
    return {
        "host": db_conf.host,
        "port": db_conf.port,
        "username": db_conf.username,
        "password": db_conf.password,
        "database": db_conf.database,
        "charset": db_conf.charset,
        "pool_size": db_conf.pool_size,
        "max_overflow": db_conf.max_overflow
    }

def get_ai_config():
    """获取 AI 配置"""
    config = get_config()
    ai_conf = config.ai
    return {
        "provider": ai_conf.provider,
        "openai_api_key": ai_conf.openai.api_key,
        "openai_base_url": ai_conf.openai.base_url,
        "openai_model": ai_conf.openai.model,
        "dashscope_api_key": ai_conf.dashscope.api_key,
        "dashscope_model": ai_conf.dashscope.model,
        "local_model_path": ai_conf.local.model_path,
        "local_model_name": ai_conf.local.model_name,
        "vector_store_type": ai_conf.vector_store.type,
        "vector_store_path": ai_conf.vector_store.path,
        "embedding_model": ai_conf.vector_store.embedding_model,
        "vector_store_embedding_model": ai_conf.vector_store.embedding_model,
        "max_tokens": ai_conf.parameters.max_tokens,
        "temperature": ai_conf.parameters.temperature,
        "top_p": ai_conf.parameters.top_p
    }

def get_rag_config():
    """获取 RAG 配置"""
    config = get_config()
    rag_conf = config.rag
    return {
        "chunk_size": rag_conf.chunk_size,
        "chunk_overlap": rag_conf.chunk_overlap,
        "similarity_threshold": rag_conf.similarity_threshold,
        "max_context_length": rag_conf.max_context_length,
        "enable_rerank": rag_conf.enable_rerank
    }

def get_service_config():
    """获取服务配置"""
    config = get_config()
    service_conf = config.service
    return {
        "host": service_conf.host,
        "port": service_conf.port,
        "max_workers": service_conf.max_workers,
        "max_message_length": service_conf.max_message_length,
        "enable_reflection": service_conf.enable_reflection
    }
