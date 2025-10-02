"""
配置管理器 - 统一管理配置传递
"""
import os
from typing import Dict, Any, Optional

# 全局配置存储
_global_config = None

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
