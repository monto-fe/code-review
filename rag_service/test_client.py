#!/usr/bin/env python3
"""
gRPC AI 服务完整测试客户端
"""
import asyncio
import grpc
import sys
import json
from pathlib import Path

# 添加 proto 目录到 Python 路径
proto_dir = Path(__file__).parent / "proto"
sys.path.insert(0, str(proto_dir))

try:
    from proto import ai_service_pb2
    from proto import ai_service_pb2_grpc
except ImportError as e:
    print(f"❌ 无法导入 gRPC 协议文件: {e}")
    print("请先运行: python3 scripts/generate_proto.py")
    sys.exit(1)


class AIServiceTester:
    """AI 服务测试器"""
    
    def __init__(self, host="localhost", port=50051):
        self.host = host
        self.port = port
        self.channel = None
        self.stub = None
    
    def connect(self):
        """连接到 gRPC 服务"""
        try:
            self.channel = grpc.insecure_channel(f"{self.host}:{self.port}")
            self.stub = ai_service_pb2_grpc.AIServiceStub(self.channel)
            print(f"✅ 已连接到 gRPC 服务: {self.host}:{self.port}")
            return True
        except Exception as e:
            print(f"❌ 连接失败: {e}")
            return False
    
    def disconnect(self):
        """断开连接"""
        if self.channel:
            self.channel.close()
            print("🔌 已断开连接")
    
    def test_health_check(self):
        """测试健康检查"""
        print("\n🏥 测试健康检查...")
        try:
            request = ai_service_pb2.HealthCheckRequest()
            response = self.stub.HealthCheck(request)
            print(f"✅ 健康检查成功: {response.status}")
            print(f"   状态码: {response.status}")
            return True
        except Exception as e:
            print(f"❌ 健康检查失败: {e}")
            return False
    
    def test_code_review(self):
        """测试代码审查"""
        print("\n🔍 测试代码审查...")
        try:
            request = ai_service_pb2.CodeReviewRequest(
                git_url="http://165.xxx.xxx.72:9980/usms/api/-/merge_requests/3",
                branch="main",
                diff_content="""
diff --git a/test.py b/test.py
index 1234567..abcdefg 100644
--- a/test.py
+++ b/test.py
@@ -1,3 +1,4 @@
 def hello():
     print("Hello World")
+    print("New line added")
""",
                agent_role="code_reviewer",
                gitlab_token="glpat-xxxxx"
            )
            response = self.stub.CodeReview(request)
            print(f"✅ 代码审查成功:")
            print(f"   智能体: {response.agent_role}")
            print(f"   建议: {response.suggestions}")
            print(f"   评分: {response.confidence_score}")
            return True
        except grpc.RpcError as e:
            print(f"❌ 代码审查失败: {e.details()}")
            return False
        except Exception as e:
            print(f"❌ 代码审查异常: {e}")
            return False
    
    def test_chat(self):
        """测试聊天功能"""
        print("\n💬 测试聊天功能...")
        try:
            request = ai_service_pb2.ChatRequest(
                message="你好，请介绍一下你的功能",
                agent_role="general_assistant",
                context="这是一个测试对话"
            )
            response = self.stub.Chat(request)
            print(f"✅ 聊天成功:")
            print(f"   智能体: {response.agent_role}")
            print(f"   回复: {response.response}")
            return True
        except grpc.RpcError as e:
            print(f"❌ 聊天失败: {e.details()}")
            if "没有可用的 AI 模型" in e.details():
                print("   💡 提示: 请在项目根目录的 .env 文件中配置 AI API 密钥")
            return False
        except Exception as e:
            print(f"❌ 聊天异常: {e}")
            return False
    
    def test_document_search(self):
        """测试文档搜索"""
        print("\n📚 测试文档搜索...")
        try:
            request = ai_service_pb2.DocumentSearchRequest(
                repository_id="test-repo-123",
                query="如何实现用户认证",
                limit=5
            )
            response = self.stub.DocumentSearch(request)
            print(f"✅ 文档搜索成功:")
            print(f"   找到 {len(response.documents)} 个相关文档")
            for i, doc in enumerate(response.documents[:3]):  # 只显示前3个
                print(f"   文档 {i+1}: {doc.file_path}")
                print(f"   内容: {doc.content[:100]}...")
            return True
        except grpc.RpcError as e:
            print(f"❌ 文档搜索失败: {e.details()}")
            return False
        except Exception as e:
            print(f"❌ 文档搜索异常: {e}")
            return False
    
    def test_available_agents(self):
        """测试获取可用智能体"""
        print("\n🤖 测试获取可用智能体...")
        try:
            request = ai_service_pb2.ChatRequest(
                message="请列出所有可用的智能体角色",
                agent_role="general_assistant"
            )
            response = self.stub.Chat(request)
            print(f"✅ 智能体信息获取成功:")
            print(f"   回复: {response.response}")
            return True
        except grpc.RpcError as e:
            print(f"❌ 获取智能体信息失败: {e.details()}")
            return False
        except Exception as e:
            print(f"❌ 获取智能体信息异常: {e}")
            return False
    
    def test_all_agents(self):
        """测试所有智能体"""
        print("\n🎭 测试所有智能体...")
        agents = [
            "code_reviewer",
            "architect", 
            "security_expert",
            "qa_engineer",
            "general_assistant"
        ]
        
        results = []
        for agent in agents:
            print(f"\n   测试智能体: {agent}")
            try:
                request = ai_service_pb2.ChatRequest(
                    message=f"你好，我是{agent}",
                    agent_role=agent
                )
                response = self.stub.Chat(request)
                print(f"   ✅ {agent} 响应成功")
                results.append((agent, True))
            except grpc.RpcError as e:
                print(f"   ❌ {agent} 响应失败: {e.details()}")
                results.append((agent, False))
            except Exception as e:
                print(f"   ❌ {agent} 异常: {e}")
                results.append((agent, False))
        
        return results
    
    def run_all_tests(self):
        """运行所有测试"""
        print("🚀 开始测试 gRPC AI 服务")
        print("=" * 60)
        
        # 连接到服务
        if not self.connect():
            print("❌ 无法连接到服务，请确保服务正在运行")
            return False
        
        try:
            # 运行测试
            tests = [
                ("健康检查", self.test_health_check),
                ("聊天功能", self.test_chat),
                ("代码审查", self.test_code_review),
                ("文档搜索", self.test_document_search),
                ("智能体信息", self.test_available_agents),
            ]
            
            results = []
            for test_name, test_func in tests:
                print(f"\n{'='*20} {test_name} {'='*20}")
                try:
                    success = test_func()
                    results.append((test_name, success))
                except Exception as e:
                    print(f"❌ {test_name} 测试异常: {e}")
                    results.append((test_name, False))
            
            # 测试所有智能体
            print(f"\n{'='*20} 所有智能体 {'='*20}")
            agent_results = self.test_all_agents()
            
            # 显示测试结果
            print("\n" + "="*60)
            print("📊 测试结果汇总:")
            print("="*60)
            
            passed = 0
            for test_name, success in results:
                status = "✅ 通过" if success else "❌ 失败"
                print(f"{test_name:15} : {status}")
                if success:
                    passed += 1
            
            print(f"\n智能体测试结果:")
            agent_passed = 0
            for agent, success in agent_results:
                status = "✅ 通过" if success else "❌ 失败"
                print(f"  {agent:15} : {status}")
                if success:
                    agent_passed += 1
            
            print(f"\n总计: {passed}/{len(results)} 个基础测试通过")
            print(f"智能体: {agent_passed}/{len(agent_results)} 个智能体可用")
            
            if passed == len(results) and agent_passed > 0:
                print("🎉 服务运行正常！")
            elif passed > 0:
                print("⚠️  服务部分功能正常，但可能需要配置 AI API 密钥")
            else:
                print("❌ 服务存在问题，请检查配置")
            
            return passed > 0
        
        finally:
            # 断开连接
            self.disconnect()


def print_config_help():
    """打印配置帮助"""
    print("""
🔧 配置 AI API 密钥

如果测试失败，请在项目根目录的 .env 文件中配置 AI API 密钥：

1. OpenAI API:
   "ai": {
     "openai": {
       "api_key": "your_openai_api_key_here"
     }
   }

2. 通义千问 API:
   "ai": {
     "dashscope": {
       "api_key": "your_dashscope_api_key_here"
     }
   }

配置完成后重启服务:
   python3 main.py
""")


def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(description="gRPC AI 服务测试客户端")
    parser.add_argument("--host", default="localhost", help="服务主机地址")
    parser.add_argument("-p", "--port", type=int, default=50051, help="服务端口")
    parser.add_argument("--config-help", action="store_true", help="显示配置帮助")
    
    args = parser.parse_args()
    
    if args.config_help:
        print_config_help()
        return
    
    # 创建测试器
    tester = AIServiceTester(host=args.host, port=args.port)
    
    # 运行测试
    try:
        success = tester.run_all_tests()
        if not success:
            print("\n" + "="*60)
            print_config_help()
    except KeyboardInterrupt:
        print("\n\n⏹️  测试被用户中断")
    except Exception as e:
        print(f"\n❌ 测试运行异常: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
