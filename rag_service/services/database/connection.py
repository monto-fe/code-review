"""
数据库连接管理
"""
import asyncio
from typing import AsyncGenerator
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
from contextlib import asynccontextmanager
import logging

from services.database.models import Base
from services.config_manager import get_database_config

logger = logging.getLogger(__name__)


class DatabaseManager:
    """数据库管理器"""
    
    def __init__(self):
        self.engine = None
        self.async_engine = None
        self.session_factory = None
        self.async_session_factory = None
    
    def init_sync_engine(self):
        """初始化同步数据库引擎"""
        db_config = get_database_config()
        database_url = (
            f"mysql+pymysql://{db_config['username']}:{db_config['password']}"
            f"@{db_config['host']}:{db_config['port']}/{db_config['database']}"
            f"?charset={db_config['charset']}"
        )
        
        self.engine = create_engine(
            database_url,
            pool_size=db_config['pool_size'],
            max_overflow=db_config['max_overflow'],
            pool_pre_ping=True,
            echo=False
        )
        
        self.session_factory = sessionmaker(
            bind=self.engine,
            autocommit=False,
            autoflush=False
        )
        
        logger.info("同步数据库引擎初始化完成")
    
    def init_async_engine(self):
        """初始化异步数据库引擎"""
        db_config = get_database_config()
        async_database_url = (
            f"mysql+aiomysql://{db_config['username']}:{db_config['password']}"
            f"@{db_config['host']}:{db_config['port']}/{db_config['database']}"
            f"?charset={db_config['charset']}"
        )
        
        self.async_engine = create_async_engine(
            async_database_url,
            pool_size=db_config['pool_size'],
            max_overflow=db_config['max_overflow'],
            pool_pre_ping=True,
            echo=False
        )
        
        self.async_session_factory = async_sessionmaker(
            bind=self.async_engine,
            class_=AsyncSession,
            autocommit=False,
            autoflush=False
        )
        
        logger.info("异步数据库引擎初始化完成")
    
    async def create_tables(self):
        """创建数据库表"""
        if self.async_engine:
            async with self.async_engine.begin() as conn:
                await conn.run_sync(Base.metadata.create_all)
            logger.info("数据库表创建完成")
        else:
            Base.metadata.create_all(bind=self.engine)
            logger.info("数据库表创建完成")
    
    def get_sync_session(self):
        """获取同步数据库会话"""
        if not self.session_factory:
            self.init_sync_engine()
        return self.session_factory()
    
    def get_async_session(self):
        """获取异步数据库会话"""
        if not self.async_session_factory:
            self.init_async_engine()
        return self.async_session_factory()
    
    @asynccontextmanager
    async def get_async_session_context(self) -> AsyncGenerator[AsyncSession, None]:
        """异步数据库会话上下文管理器"""
        if not self.async_session_factory:
            self.init_async_engine()
        
        async with self.async_session_factory() as session:
            try:
                yield session
                await session.commit()
            except Exception:
                await session.rollback()
                raise
            finally:
                await session.close()
    
    async def close(self):
        """关闭数据库连接"""
        if self.async_engine:
            await self.async_engine.dispose()
        if self.engine:
            self.engine.dispose()
        logger.info("数据库连接已关闭")


# 全局数据库管理器实例
db_manager = DatabaseManager()


def get_sync_db():
    """获取同步数据库会话（依赖注入）"""
    db = db_manager.get_sync_session()
    try:
        yield db
    finally:
        db.close()


async def get_async_db():
    """获取异步数据库会话（依赖注入）"""
    async with db_manager.get_async_session_context() as session:
        yield session
