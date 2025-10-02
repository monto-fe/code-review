#!/usr/bin/env python3
"""
生成 gRPC 协议文件
"""
import os
import subprocess
import sys
from pathlib import Path

def generate_proto_files():
    """生成 gRPC 协议文件"""
    # 获取项目根目录
    project_root = Path(__file__).parent.parent
    proto_dir = project_root / "proto"
    output_dir = project_root / "proto"
    
    # 确保输出目录存在
    output_dir.mkdir(exist_ok=True)
    
    # 生成 Python 文件
    proto_file = proto_dir / "ai_service.proto"
    
    if not proto_file.exists():
        print(f"错误: 协议文件不存在: {proto_file}")
        return False
    
    try:
        # 使用 grpc_tools.protoc 生成 Python 文件
        cmd = [
            sys.executable, "-m", "grpc_tools.protoc",
            f"--proto_path={proto_dir}",
            f"--python_out={output_dir}",
            f"--grpc_python_out={output_dir}",
            str(proto_file)
        ]
        
        print(f"执行命令: {' '.join(cmd)}")
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            print("gRPC 协议文件生成成功")
            return True
        else:
            print(f"生成失败: {result.stderr}")
            return False
            
    except Exception as e:
        print(f"生成协议文件时出错: {e}")
        return False

if __name__ == "__main__":
    success = generate_proto_files()
    sys.exit(0 if success else 1)
