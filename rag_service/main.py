"""
AI 服务主入口
"""
import asyncio
import logging
import signal
import sys
import os
from pathlib import Path
from dynaconf import Dynaconf
from dotenv import load_dotenv

# 添加项目根目录到 Python 路径
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

# 加载项目根目录的 .env 文件
project_root_env = project_root.parent / ".env"
if project_root_env.exists():
    load_dotenv(project_root_env)
    print(f"已加载环境变量文件: {project_root_env}")
else:
    print(f"未找到环境变量文件: {project_root_env}")

# 从环境变量加载配置
from services.config_manager import load_env_config
env_config = load_env_config()

# 使用环境变量配置
print("使用环境变量配置")
settings = Dynaconf(
    **env_config,
    environments=False
)

# 配置日志
log_level = getattr(logging, settings.logging.level.upper(), logging.INFO)
log_file = getattr(settings.logging, 'file', 'ai_service.log')
logging.basicConfig(
    level=log_level,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler(log_file)
    ]
)

logger = logging.getLogger(__name__)

# 设置全局配置
from services.config_manager import set_config
set_config(settings)


class AIService:
    """AI 服务主类"""
    
    def __init__(self):
        self.settings = settings
        self.grpc_server = None
        self.running = False
    
    # 配置方法已移至 config_manager.py，直接使用全局方法
    
    async def initialize(self):
        """初始化服务"""
        try:
            logger.info(f"正在初始化 AI 服务 (环境: {self.settings.current_env})...")
            
            # 设置环境变量
            from services.config_manager import get_ai_config
            ai_config = get_ai_config()
            os.environ.setdefault("AI_PROVIDER", ai_config["provider"])
            os.environ.setdefault("OPENAI_API_KEY", ai_config["openai_api_key"])
            os.environ.setdefault("DASHSCOPE_API_KEY", ai_config["dashscope_api_key"])
            os.environ.setdefault("DASHSCOPE_MODEL", ai_config["dashscope_model"])
            os.environ.setdefault("VECTOR_STORE_TYPE", ai_config["vector_store_type"])
            os.environ.setdefault("VECTOR_STORE_PATH", ai_config["vector_store_path"])
            os.environ.setdefault("GITLAB_TOKEN", self.settings.git.token)
            
            # 初始化数据库
            logger.info("初始化数据库连接...")
            from services.database.connection import db_manager
            db_manager.init_sync_engine()
            db_manager.init_async_engine()
            await db_manager.create_tables()
            
            # 初始化其他组件
            logger.info("初始化 AI 模型管理器...")
            from services.ai.model_manager import model_manager
            available_models = model_manager.list_models()
            logger.info(f"可用 AI 模型: {available_models}")
            
            logger.info("初始化智能体管理器...")
            from services.ai.agent_manager import get_agent_manager
            agent_manager = get_agent_manager()
            agents = agent_manager.list_agents()
            logger.info(f"可用智能体: {[agent.name for agent in agents]}")
            
            logger.info("初始化 RAG 服务...")
            from services.rag.rag_service import rag_service
            logger.info("RAG 服务初始化完成")
            
            # 初始化 gRPC 服务器
            logger.info("初始化 gRPC 服务器...")
            from services.grpc.ai_service_impl import grpc_server
            self.grpc_server = grpc_server
            
            logger.info("AI 服务初始化完成")
            
        except Exception as e:
            logger.error(f"服务初始化失败: {e}")
            raise
    
    async def start(self):
        """启动服务"""
        try:
            await self.initialize()
            
            # 设置信号处理
            self._setup_signal_handlers()
            
            # 启动 gRPC 服务器
            from services.config_manager import get_service_config
            service_config = get_service_config()
            logger.info(f"启动 gRPC 服务器 {service_config['host']}:{service_config['port']}...")
            self.running = True
            
            # 启动 gRPC 服务器（不等待终止）
            await self.grpc_server.start_async()
            
            # 主循环，检查运行状态
            while self.running:
                try:
                    await asyncio.sleep(0.1)
                except KeyboardInterrupt:
                    logger.info("收到键盘中断信号")
                    self.running = False
                    break
            
        except Exception as e:
            logger.error(f"服务启动失败: {e}")
            await self.stop()
            raise
    
    async def stop(self):
        """停止服务"""
        if self.running:
            logger.info("正在停止 AI 服务...")
            self.running = False
            
            try:
                # 停止 gRPC 服务器（带超时）
                if self.grpc_server:
                    await asyncio.wait_for(
                        self.grpc_server.stop(), 
                        timeout=10.0
                    )
                
                # 关闭数据库连接（带超时）
                from services.database.connection import db_manager
                await asyncio.wait_for(
                    db_manager.close(), 
                    timeout=5.0
                )
                
                logger.info("AI 服务已停止")
                
            except asyncio.TimeoutError:
                logger.warning("服务停止超时，强制退出")
            except Exception as e:
                logger.error(f"停止服务时出错: {e}")
    
    def _setup_signal_handlers(self):
        """设置信号处理器"""
        def signal_handler(signum, frame):
            logger.info(f"收到信号 {signum}，正在关闭服务...")
            # 设置停止标志，让主循环处理
            self.running = False
        
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)


async def main():
    """主函数"""
    service = AIService()
    
    try:
        logger.info(f"🚀 启动 AI 服务 (环境: {service.settings.current_env})")
        await service.start()
    except KeyboardInterrupt:
        logger.info("收到键盘中断信号")
    except Exception as e:
        logger.error(f"服务运行异常: {e}")
        import traceback
        traceback.print_exc()
    finally:
        logger.info("开始清理资源...")
        await service.stop()
        logger.info("服务已完全停止")


if __name__ == "__main__":
    try:
        # 运行服务
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n收到中断信号，正在退出...")
    except Exception as e:
        print(f"程序异常退出: {e}")
        import traceback
        traceback.print_exc()
    finally:
        print("程序已退出")
