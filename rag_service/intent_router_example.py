#!/usr/bin/env python3
"""
意图识别和智能路由使用示例
"""
import asyncio
import logging
import sys
import os
from pathlib import Path

# 添加项目根目录到 Python 路径
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

# 设置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

async def main():
    """演示意图识别和智能路由功能"""
    try:
        # 设置环境变量
        os.environ.setdefault("DASHSCOPE_API_KEY", "sk-xxxx")
        os.environ.setdefault("DASHSCOPE_MODEL", "qwen-turbo")
        
        # 初始化配置
        from services.config_manager import set_config
        from dynaconf import Dynaconf
        settings = Dynaconf(
            settings_files=["config/settings.json"],
            environments=True,
            default_env="testing"
        )
        set_config(settings)
        
        # 获取智能体管理器
        from services.ai.agent_manager import get_agent_manager
        agent_manager = get_agent_manager()
        
        print("🤖 意图识别和智能路由演示")
        print("=" * 50)
        
        # 演示智能聊天
        print("\n📝 智能聊天演示:")
        print("-" * 30)
        
        messages = [
            "这段代码有SQL注入漏洞吗？",
            "如何设计一个高并发的微服务架构？",
            "这个函数的测试覆盖率够吗？",
            "请帮我优化这段代码的性能",
            "请全面分析这段代码，包括安全性、架构设计和测试策略",
            "设计一个安全的微服务架构，并考虑测试方案"
        ]
        
        for message in messages:
            print(f"\n用户: {message}")
            result = await agent_manager.smart_chat(message)
            
            if len(result['agent_names']) > 1:
                print(f"🤖 多专家协作 ({', '.join(result['agent_names'])}): {result['response'][:100]}...")
                print(f"   协作模式: {result['collaboration_mode']} | 意图: {result['intent']} | 置信度: {result['confidence']:.2f}")
            else:
                print(f"🤖 {result['agent_names'][0]}: {result['response'][:100]}...")
                print(f"   意图: {result['intent']} | 置信度: {result['confidence']:.2f}")
        
        # 演示智能代码审查
        print("\n\n🔍 智能代码审查演示:")
        print("-" * 30)
        
        code_examples = [
            {
                "code": "def login(username, password):\n    query = f\"SELECT * FROM users WHERE username='{username}' AND password='{password}'\"\n    return execute_query(query)",
                "description": "SQL注入漏洞代码"
            },
            {
                "code": "class UserService:\n    def __init__(self):\n        self.db = Database()\n    \n    def get_user(self, user_id):\n        return self.db.query(f\"SELECT * FROM users WHERE id={user_id}\")",
                "description": "架构设计问题代码"
            }
        ]
        
        for example in code_examples:
            print(f"\n代码示例: {example['description']}")
            print(f"代码:\n{example['code']}")
            
            result = await agent_manager.smart_code_review(
                code_changes=example['code'],
                context="这是一个用户认证相关的代码片段",
                query="请分析这段代码的问题"
            )
            
            print(f"🤖 {result['agent_name']}: {result['review'][:150]}...")
            print(f"   选择理由: {result['reasoning']}")
        
        print("\n\n✅ 演示完成！")
        print("\n💡 使用说明:")
        print("- 在 gRPC 请求中设置 agent_role='auto' 启用智能路由")
        print("- 或者直接调用 agent_manager.smart_chat() 和 smart_code_review() 方法")
        print("- 系统会自动分析用户意图并选择最合适的专业Agent或多个Agent协作")
        print("\n🔄 协作模式:")
        print("- single: 单个Agent处理")
        print("- sequential: 多个Agent按顺序处理，后一个基于前一个的结果")
        print("- parallel: 多个Agent并行处理，最后合并结果")
        print("\n🎯 协作场景:")
        print("- 代码审查 + 安全检查：code_reviewer + security_expert")
        print("- 架构设计 + 测试策略：architect + qa_engineer")
        print("- 全面代码分析：code_reviewer + security_expert + architect + qa_engineer")
        
    except Exception as e:
        logger.error(f"演示失败: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    asyncio.run(main())
