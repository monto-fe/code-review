#!/usr/bin/env python3
"""
gRPC AI 服务基础测试客户端
"""
import grpc
import sys
import os
from pathlib import Path

# 设置环境变量（模拟服务启动时的环境）
os.environ.setdefault("ENV", "testing")
os.environ.setdefault("AI_PROVIDER", "dashscope")
os.environ.setdefault("DASHSCOPE_API_KEY", "")
os.environ.setdefault("DASHSCOPE_MODEL", "qwen-turbo")
os.environ.setdefault("VECTOR_STORE_TYPE", "faiss")
os.environ.setdefault("VECTOR_STORE_PATH", "/app/data/vector_store")

# 添加 proto 目录到 Python 路径
proto_dir = Path(__file__).parent / "proto"
sys.path.insert(0, str(proto_dir))

def test_basic_connection():
    """测试基本连接"""
    print("🔍 测试 gRPC 服务基本连接...")
    
    try:
        # 导入协议文件
        from proto import ai_service_pb2
        from proto import ai_service_pb2_grpc
        print("✅ 协议文件导入成功")
    except ImportError as e:
        print(f"❌ 协议文件导入失败: {e}")
        return False
    
    try:
        # 创建连接
        channel = grpc.insecure_channel("localhost:50051")
        
        # 测试连接
        grpc.channel_ready_future(channel).result(timeout=5)
        print("✅ gRPC 服务连接成功")
        
        # 创建存根
        stub = ai_service_pb2_grpc.AIServiceStub(channel)
        print("✅ gRPC 存根创建成功")
        
        # 测试健康检查
        print("\n🏥 测试健康检查...")
        try:
            request = ai_service_pb2.HealthCheckRequest()
            response = stub.HealthCheck(request)
            print(f"✅ 健康检查成功: 状态码 {response.status}")
        except Exception as e:
            print(f"❌ 健康检查失败: {e}")
        
        # 测试聊天功能（预期会失败，因为没有配置 AI 模型）
        print("\n💬 测试聊天功能...")
        try:
            request = ai_service_pb2.ChatRequest(
                message="你好",
                agent_role="general_assistant"
            )
            response = stub.Chat(request)
            print(f"✅ 聊天成功: {response.response}")
        except grpc.RpcError as e:
            if "没有可用的 AI 模型" in e.details():
                print("⚠️  聊天失败: 没有配置 AI 模型（这是预期的）")
                print("   💡 要启用 AI 功能，请在 config/settings.json 中配置 API 密钥")
            else:
                print(f"❌ 聊天失败: {e.details()}")
        except Exception as e:
            print(f"❌ 聊天异常: {e}")
        
        # 测试代码审查功能
        print("\n🔍 测试代码审查...")
        try:
            request = ai_service_pb2.CodeReviewRequest(
                git_url="http://16.xxx.xxx.72:9980/usms/api",
                branch="main",
                diff_content="test diff",
                agent_role="code_reviewer",
                gitlab_token="glpat-snxxxx"
            )
            response = stub.CodeReview(request)
            print(f"✅ 代码审查成功: {response.suggestions}")
        except grpc.RpcError as e:
            if "没有可用的 AI 模型" in e.details():
                print("⚠️  代码审查失败: 没有配置 AI 模型（这是预期的）")
            else:
                print(f"❌ 代码审查失败: {e.details()}")
        except Exception as e:
            print(f"❌ 代码审查异常: {e}")
        
        # 测试文档搜索功能
        print("\n📚 测试文档搜索...")
        try:
            request = ai_service_pb2.DocumentSearchRequest(
                repository_id="test-repo-123",
                query="test query",
                limit=5
            )
            response = stub.DocumentSearch(request)
            print(f"✅ 文档搜索成功: 找到 {len(response.documents)} 个文档")
        except grpc.RpcError as e:
            if "没有可用的 AI 模型" in e.details():
                print("⚠️  文档搜索失败: 没有配置 AI 模型（这是预期的）")
            else:
                print(f"❌ 文档搜索失败: {e.details()}")
        except Exception as e:
            print(f"❌ 文档搜索异常: {e}")
        
        channel.close()
        print("\n🎉 基础连接测试完成！")
        return True
        
    except grpc.RpcError as e:
        print(f"❌ gRPC 调用失败: {e}")
        return False
    except Exception as e:
        print(f"❌ 连接失败: {e}")
        return False

def check_service_status():
    """检查服务状态"""
    print("🔍 检查服务状态...")
    
    import subprocess
    try:
        result = subprocess.run(
            ["lsof", "-i", ":50051"], 
            capture_output=True, 
            text=True
        )
        
        if result.returncode == 0:
            print("✅ 端口 50051 正在被使用")
            lines = result.stdout.strip().split('\n')
            if len(lines) > 1:
                print(f"   进程信息: {lines[1]}")
            return True
        else:
            print("❌ 端口 50051 未被使用，服务可能未启动")
            return False
    except FileNotFoundError:
        print("⚠️  无法使用 lsof 命令检查端口")
        return False

def print_config_instructions():
    """打印配置说明"""
    print("""
🔧 配置 AI API 密钥以启用完整功能

当前服务运行正常，但 AI 功能需要配置 API 密钥。

1. 编辑配置文件:
   vim config/settings.json

2. 在 "ai" 部分添加 API 密钥:
   {
     "ai": {
       "openai": {
         "api_key": "your_openai_api_key_here"
       },
       "dashscope": {
         "api_key": "your_dashscope_api_key_here"
       }
     }
   }

3. 重启服务:
   python3 main.py

支持的 AI 服务:
- OpenAI (GPT-3.5, GPT-4)
- 通义千问 (Qwen)
- 本地模型 (Hugging Face)
""")

def main():
    """主函数"""
    print("🚀 gRPC AI 服务基础测试")
    print("=" * 50)
    
    # 检查服务状态
    if not check_service_status():
        print("\n请先启动服务:")
        print("python3 main.py")
        return
    
    print("\n" + "=" * 50)
    
    # 测试连接
    if test_basic_connection():
        print("\n" + "=" * 50)
        print("📊 测试结果:")
        print("✅ gRPC 服务连接正常")
        print("✅ 协议文件加载正常")
        print("✅ 基础功能可用")
        print("⚠️  AI 功能需要配置 API 密钥")
        print("\n" + "=" * 50)
        print_config_instructions()
    else:
        print("\n❌ 测试失败，请检查服务配置")

if __name__ == "__main__":
    main()
