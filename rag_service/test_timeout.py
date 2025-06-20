#!/usr/bin/env python3
"""
超时功能测试脚本
"""

import time
import signal
import functools
import requests
import json

def timeout_handler(signum, frame):
    raise TimeoutError("操作超时")

def timeout(seconds=60):
    """
    超时装饰器，在指定秒数后中断函数执行
    """
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            # 设置信号处理器
            old_handler = signal.signal(signal.SIGALRM, timeout_handler)
            signal.alarm(seconds)
            
            try:
                result = func(*args, **kwargs)
                return result
            finally:
                # 恢复原来的信号处理器
                signal.alarm(0)
                signal.signal(signal.SIGALRM, old_handler)
        return wrapper
    return decorator

@timeout(5)  # 5秒超时测试
def long_running_function():
    """模拟长时间运行的函数"""
    print("开始长时间运行...")
    time.sleep(10)  # 睡眠10秒
    print("长时间运行完成")
    return "完成"

def test_timeout_decorator():
    """测试超时装饰器"""
    print("测试超时装饰器...")
    try:
        result = long_running_function()
        print(f"结果: {result}")
    except TimeoutError as e:
        print(f"超时异常: {e}")
    except Exception as e:
        print(f"其他异常: {e}")

def test_rag_service_timeout():
    """测试RAG服务超时"""
    print("测试RAG服务超时...")
    
    # 准备测试数据
    test_data = {
        "git_url": "https://gitlab.com/test/repo.git",
        "branch": "main",
        "diff_content": "diff --git a/test.py b/test.py\n@@ -1,1 +1,1 @@\n-print('old')\n+print('new')",
        "query": "测试查询",
        "gitlab_token": "invalid_token"  # 使用无效token来触发超时
    }
    
    try:
        # 设置5秒超时
        response = requests.post(
            "http://localhost:8000/api/code-analysis",
            json=test_data,
            timeout=5
        )
        print(f"响应状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
    except requests.exceptions.Timeout:
        print("HTTP请求超时")
    except requests.exceptions.ConnectionError:
        print("连接错误，RAG服务可能未启动")
    except Exception as e:
        print(f"其他异常: {e}")

def test_go_client_timeout():
    """测试Go客户端超时"""
    print("测试Go客户端超时...")
    
    # 这里可以添加Go客户端的超时测试
    # 由于这是Python脚本，我们只打印说明
    print("Go客户端超时测试需要在Go代码中进行")
    print("可以通过设置无效的RAG服务URL来测试超时")

if __name__ == "__main__":
    print("=== 超时功能测试 ===")
    
    # 测试Python装饰器超时
    test_timeout_decorator()
    
    print("\n" + "="*50 + "\n")
    
    # 测试RAG服务超时
    test_rag_service_timeout()
    
    print("\n" + "="*50 + "\n")
    
    # 测试Go客户端超时
    test_go_client_timeout()
    
    print("\n=== 测试完成 ===") 