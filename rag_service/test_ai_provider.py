#!/usr/bin/env python3
"""
测试AI模型选择功能
"""
import sys
import asyncio
from pathlib import Path

# 添加项目根目录到 Python 路径
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

from main import settings
from services.config_manager import set_config, get_ai_config

# 设置配置
set_config(settings)

async def test_ai_provider():
    """测试AI模型选择功能"""
    print("🔍 测试AI模型选择功能...")
    
    # 设置环境变量
    import os
    ai_config = get_ai_config()
    os.environ["AI_PROVIDER"] = ai_config["provider"]
    os.environ["DASHSCOPE_API_KEY"] = ai_config["dashscope_api_key"]
    os.environ["DASHSCOPE_MODEL"] = ai_config["dashscope_model"]
    
    print(f"当前配置的AI提供者: {ai_config['provider']}")
    
    # 测试模型管理器
    from services.ai.model_manager import model_manager
    print(f"可用模型: {model_manager.list_models()}")
    
    # 测试获取默认模型
    try:
        default_model = model_manager.get_model()
        print(f"✅ 默认模型: {type(default_model).__name__}")
    except Exception as e:
        print(f"❌ 获取默认模型失败: {e}")
    
    # 测试智能体管理器
    from services.ai.agent_manager import agent_manager
    print(f"智能体数量: {len(agent_manager.list_agents())}")
    
    # 测试智能体聊天
    try:
        response = await agent_manager.chat_with_agent(
            role_key="general_assistant",
            message="你好，请介绍一下你自己"
        )
        print(f"✅ 智能体聊天成功: {response[:100]}...")
    except Exception as e:
        print(f"❌ 智能体聊天失败: {e}")
        import traceback
        traceback.print_exc()
    
    # 测试不同提供者的切换
    print("\n🔄 测试提供者切换...")
    
    # 测试指定DashScope
    try:
        dashscope_model = model_manager.get_model("dashscope")
        print(f"✅ DashScope模型: {type(dashscope_model).__name__}")
    except Exception as e:
        print(f"❌ 获取DashScope模型失败: {e}")
    
    # 测试指定OpenAI（如果配置了）
    try:
        openai_model = model_manager.get_model("openai")
        print(f"✅ OpenAI模型: {type(openai_model).__name__}")
    except Exception as e:
        print(f"⚠️  获取OpenAI模型失败: {e} (这是预期的，因为没有配置API密钥)")

if __name__ == "__main__":
    asyncio.run(test_ai_provider())
