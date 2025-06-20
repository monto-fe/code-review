import os
from typing import List, Optional, Dict
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.vectorstores import FAISS
from langchain.embeddings import HuggingFaceEmbeddings
# from langchain.chains import LLMChain
# from langchain.prompts import PromptTemplate
# from langchain.llms import OpenAI
from git import Repo
# import tempfile
import shutil
import hashlib
# from pathlib import Path
from urllib.parse import urlparse
import sys
import subprocess
import re

# 代码审查模板
CODE_REVIEW_TEMPLATE = """
请对以下代码变更进行审查，并提供详细的建议：

代码变更：
{code_changes}

相关上下文：
{context}

请从以下几个方面提供建议：
1. 代码质量：代码的可读性、可维护性、性能等方面
2. 潜在问题：可能存在的bug、安全隐患、资源泄漏等
3. 改进建议：具体的优化方案和最佳实践
4. 示例代码：如果可能，提供改进后的代码示例

请以结构化的方式输出审查结果。
"""

def check_dependencies():
    """
    检查必要的依赖是否已安装
    """
    # 检查是否跳过依赖检查
    skip_deps_check = os.getenv("SKIP_DEPS_CHECK", "false").lower() == "true"
    if skip_deps_check:
        print("跳过依赖检查（SKIP_DEPS_CHECK=true）")
        return
    
    dependencies = [
        ("numpy", "numpy<2.0.0"),  # 指定安装 NumPy 1.x 版本
        ("torch", "torch==2.1.2"),  # 指定 torch 版本
        ("torchvision", "torchvision==0.16.2"),  # 指定 torchvision 版本
        "faiss-cpu",
        "sentence-transformers"
    ]
    
    missing_deps = []
    
    for dep in dependencies:
        try:
            if isinstance(dep, tuple):
                package_name, install_name = dep
            else:
                package_name = install_name = dep
                
            if package_name == "sentence-transformers":
                import sentence_transformers
                print("sentence-transformers 已安装")
            elif package_name == "numpy":
                import numpy
                print("numpy 已安装")
            elif package_name == "faiss-cpu":
                import faiss
                print("faiss-cpu 已安装")
            elif package_name == "torch":
                import torch
                print("torch 已安装")
            elif package_name == "torchvision":
                import torchvision
                print("torchvision 已安装")
        except ImportError:
            print(f"检测到缺失依赖: {package_name}")
            missing_deps.append((package_name, install_name))
    
    # 只在有缺失依赖时才安装
    if missing_deps:
        print(f"正在安装缺失的依赖: {[dep[0] for dep in missing_deps]}")
        for package_name, install_name in missing_deps:
            try:
                # 使用 subprocess 安装依赖
                subprocess.check_call([
                    sys.executable, 
                    "-m", 
                    "pip", 
                    "install", 
                    install_name
                ])
                print(f"{package_name} 安装完成")
                
                # 验证安装
                if package_name == "sentence-transformers":
                    import sentence_transformers
                elif package_name == "numpy":
                    import numpy
                elif package_name == "faiss-cpu":
                    import faiss
                elif package_name == "torch":
                    import torch
                elif package_name == "torchvision":
                    import torchvision
                print(f"{package_name} 验证成功")
            except Exception as e:
                print(f"安装 {package_name} 失败: {str(e)}")
                print(f"尝试使用备用方法安装 {package_name}...")
                try:
                    # 尝试使用 pip 直接安装
                    subprocess.check_call([
                        "pip", 
                        "install", 
                        install_name
                    ])
                    print(f"使用备用方法安装 {package_name} 成功")
                except Exception as e2:
                    print(f"备用安装方法也失败: {str(e2)}")
                    raise Exception(f"无法安装 {package_name}，请手动运行: pip install {install_name}")
    else:
        print("所有依赖已安装，跳过依赖检查")

def get_vector_store_path(git_url: str, branch: str) -> str:
    """
    根据git_url和branch生成向量存储的路径
    
    Args:
        git_url: Git仓库URL
        branch: 分支名
    
    Returns:
        str: 向量存储的路径
    """
    # 使用git_url和branch生成唯一的hash
    repo_hash = hashlib.md5(f"{git_url}:{branch}".encode()).hexdigest()
    # 在用户目录下创建.rag_cache目录
    cache_dir = os.path.expanduser("~/.rag_cache")
    os.makedirs(cache_dir, exist_ok=True)
    return os.path.join(cache_dir, f"vector_store_{repo_hash}")

def get_git_url_with_token(git_url: str, gitlab_token: str) -> str:
    """
    将GitLab token添加到git URL中
    
    Args:
        git_url: 原始git URL
        gitlab_token: GitLab访问令牌
    
    Returns:
        str: 带有token的git URL
    """
    # 统一转换为HTTP格式并添加token
    if git_url.startswith('git@'):
        # 从SSH格式转换为HTTP格式
        parts = git_url.split(':', 2)  # 最多分割2次
        if len(parts) == 3:
            host, port, path = parts
            # 移除git@前缀
            host = host.replace('git@', '')
            return f"http://oauth2:{gitlab_token}@{host}:{port}/{path}"
        else:
            # 标准SSH格式 (git@host:path)
            host = git_url.split('@')[1].split(':')[0]
            path = git_url.split(':')[1]
            return f"http://oauth2:{gitlab_token}@{host}/{path}"
    else:
        # 处理HTTP/HTTPS格式的URL
        parsed_url = urlparse(git_url)
        if parsed_url.scheme in ['http', 'https']:
            # 处理带端口号的情况
            netloc = parsed_url.netloc
            if ':' in netloc:
                host, port = netloc.split(':')
                return f"http://oauth2:{gitlab_token}@{host}:{port}{parsed_url.path}"
            else:
                return f"http://oauth2:{gitlab_token}@{netloc}{parsed_url.path}"
        else:
            raise ValueError(f"不支持的URL格式: {git_url}")

def read_branch_files(repo_path: str, branch: str) -> Dict[str, str]:
    """
    读取指定分支的所有文件内容
    
    Args:
        repo_path: 仓库本地路径
        branch: 分支名
    
    Returns:
        Dict[str, str]: 文件路径到文件内容的映射
    """
    repo = Repo(repo_path)
    repo.git.checkout(branch)
    
    files_content = {}
    for root, _, files in os.walk(repo_path):
        for file in files:
            # 跳过.git目录和临时文件
            if '.git' in root or file.startswith('.'):
                continue
                
            file_path = os.path.join(root, file)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    # 使用相对路径作为键
                    rel_path = os.path.relpath(file_path, repo_path)
                    files_content[rel_path] = f.read()
            except Exception as e:
                print(f"读取文件 {file_path} 失败: {str(e)}")
                continue
    
    return files_content

def extract_files_from_diff(diff_content: str) -> List[str]:
    """
    从diff内容中提取变更的文件路径
    
    Args:
        diff_content: diff字符串
    
    Returns:
        List[str]: 变更的文件路径列表
    """
    files = []
    for line in diff_content.split('\n'):
        if line.startswith('diff --git'):
            # 提取文件路径
            parts = line.split()
            if len(parts) >= 3:
                file_path = parts[2].replace('b/', '')
                files.append(file_path)
    return files

def extract_query_from_diff(diff_content: str) -> str:
    """
    从diff内容中提取关键信息构建查询
    
    Args:
        diff_content: diff字符串
    
    Returns:
        str: 构建的查询字符串
    """
    query_parts = []
    
    # 1. 提取文件类型和扩展名
    file_extensions = set()
    file_pattern = r'diff --git a/.*?\.(\w+)'
    extensions = re.findall(file_pattern, diff_content)
    for ext in extensions:
        if ext.lower() in ['py', 'js', 'ts', 'java', 'go', 'cpp', 'c', 'php', 'rb', 'rs']:
            file_extensions.add(ext.lower())
    
    if file_extensions:
        query_parts.append(f"编程语言: {', '.join(file_extensions)}")
    
    # 2. 提取函数名
    function_pattern = r'\+def\s+(\w+)'
    functions = re.findall(function_pattern, diff_content)
    if functions:
        query_parts.append(f"函数: {', '.join(functions)}")
    
    # 3. 提取类名
    class_pattern = r'\+class\s+(\w+)'
    classes = re.findall(class_pattern, diff_content)
    if classes:
        query_parts.append(f"类: {', '.join(classes)}")
    
    # 4. 提取方法名
    method_pattern = r'\+(\s+def\s+\w+)'
    methods = re.findall(method_pattern, diff_content)
    if methods:
        method_names = [m.strip().split()[1] for m in methods if len(m.strip().split()) > 1]
        if method_names:
            query_parts.append(f"方法: {', '.join(method_names)}")
    
    # 5. 提取语义关键词
    semantic_keywords = []
    
    # API相关
    if any(keyword in diff_content.lower() for keyword in ['api', 'http', 'request', 'response', 'endpoint']):
        semantic_keywords.append('API接口')
    
    # 数据库相关
    if any(keyword in diff_content.lower() for keyword in ['database', 'db', 'sql', 'query', 'select', 'insert', 'update', 'delete']):
        semantic_keywords.append('数据库操作')
    
    # 认证相关
    if any(keyword in diff_content.lower() for keyword in ['auth', 'login', 'password', 'token', 'jwt', 'oauth']):
        semantic_keywords.append('认证授权')
    
    # 错误处理
    if any(keyword in diff_content.lower() for keyword in ['error', 'exception', 'try', 'catch', 'finally']):
        semantic_keywords.append('错误处理')
    
    # 配置相关
    if any(keyword in diff_content.lower() for keyword in ['config', 'setting', 'env', 'environment']):
        semantic_keywords.append('配置管理')
    
    if semantic_keywords:
        query_parts.append(f"功能类型: {', '.join(semantic_keywords)}")
    
    return " ".join(query_parts)

def build_optimal_query(diff_content: str, original_query: str = None, merge_request_info: dict = None) -> str:
    """
    构建最优的查询字符串
    
    Args:
        diff_content: diff内容
        original_query: 原始查询（来自Go后端）
        merge_request_info: 合并请求信息
    
    Returns:
        str: 最优查询字符串
    """
    query_parts = []
    
    # 1. 从diff中提取关键信息
    diff_query = extract_query_from_diff(diff_content)
    if diff_query:
        query_parts.append(diff_query)
    
    # 2. 添加原始查询（如果提供）
    if original_query:
        query_parts.append(original_query)
    
    # 3. 添加MR信息（如果提供）
    if merge_request_info:
        title = merge_request_info.get('title', '')
        description = merge_request_info.get('description', '')
        
        if title:
            query_parts.append(f"功能: {title}")
        if description:
            # 提取描述中的关键词
            desc_keywords = extract_keywords_from_description(description)
            if desc_keywords:
                query_parts.append(f"描述关键词: {', '.join(desc_keywords)}")
    
    return " ".join(query_parts)

def extract_keywords_from_description(description: str) -> List[str]:
    """
    从MR描述中提取关键词
    
    Args:
        description: MR描述
    
    Returns:
        List[str]: 关键词列表
    """
    keywords = []
    
    # 常见的关键词模式
    keyword_patterns = [
        r'修复\s*(\w+)',  # 修复xxx
        r'添加\s*(\w+)',  # 添加xxx
        r'优化\s*(\w+)',  # 优化xxx
        r'重构\s*(\w+)',  # 重构xxx
        r'更新\s*(\w+)',  # 更新xxx
    ]
    
    for pattern in keyword_patterns:
        matches = re.findall(pattern, description)
        keywords.extend(matches)
    
    return list(set(keywords))  # 去重

def analyze_gitlab_code(git_url: str, branch: str, diff_content: str, query: Optional[str] = None, gitlab_token: Optional[str] = None) -> Dict[str, str]:
    """
    分析GitLab代码变更
    
    Args:
        git_url: Git仓库URL
        branch: 分支名
        diff_content: 格式化的diff字符串
        query: 可选的查询参数
        gitlab_token: GitLab访问令牌
    
    Returns:
        Dict[str, str]: 代码变更和上下文作为字典返回
    """
    # 检查依赖
    print("检查依赖")
    check_dependencies()
    print("依赖检查完成")
    
    # 检查GitLab token
    if not gitlab_token:
        gitlab_token = os.getenv("GITLAB_TOKEN")
        if not gitlab_token:
            raise ValueError("未提供GitLab token，请通过参数或环境变量GITLAB_TOKEN提供")
    
    # 使用配置的临时目录，确保使用绝对路径
    git_temp_dir = os.path.abspath(os.getenv("GIT_TEMP_DIR", "git_dir"))
    os.makedirs(git_temp_dir, exist_ok=True)
    
    # 创建唯一的子目录
    repo_name = hashlib.md5(f"{git_url}:{branch}".encode()).hexdigest()
    temp_dir = os.path.join(git_temp_dir, repo_name)
    
    print(f"使用临时目录: {temp_dir}")
    
    try:
        # 如果目录已存在，先删除
        if os.path.exists(temp_dir):
            shutil.rmtree(temp_dir)
        
        # 确保目录存在
        os.makedirs(temp_dir, exist_ok=True)
        
        # 获取带有token的git URL
        auth_git_url = get_git_url_with_token(git_url, gitlab_token)
        print(f"auth_git_url: {auth_git_url}")
        
        try:
            # 配置Git
            os.environ['GIT_SSL_NO_VERIFY'] = 'true'
            os.environ['GIT_HTTP_LOW_SPEED_TIME'] = '300'
            os.environ['GIT_HTTP_LOW_SPEED_LIMIT'] = '1000'
            
            # 克隆仓库
            repo = Repo.clone_from(auth_git_url, temp_dir)
            print(f"repo: {repo}")
        except Exception as e:
            print(f"克隆仓库失败: {str(e)}")
            raise Exception(f"无法克隆仓库: {str(e)}")
        
        # 读取整个分支的文件
        branch_files = read_branch_files(temp_dir, branch)
        print(f"读取到 {len(branch_files)} 个文件")
        
        # 从diff中提取变更的文件
        changed_files = extract_files_from_diff(diff_content)
        print(f"changed_files: {changed_files}")
        
        # 获取向量存储路径
        print(f"git_url: {git_url}")
        print(f"branch: {branch}")
        vector_store_path = get_vector_store_path(git_url, branch)
        print(f"vector_store_path: {vector_store_path}")
        
        # 将代码分割成块
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200
        )
        
        # 处理所有分支文件用于上下文
        all_texts = []
        for file_path, content in branch_files.items():
            all_texts.append(f"文件: {file_path}\n内容:\n{content}\n")
        
        if not all_texts:
            raise ValueError("没有找到任何文件内容")
        
        print(f"准备处理 {len(all_texts)} 个文本块")
        
        # 分割所有文本
        texts = text_splitter.split_text("\n".join(all_texts))
        print(f"分割后得到 {len(texts)} 个文本块")
        
        try:
            # 创建向量存储，使用开源的嵌入模型
            model_name = "sentence-transformers/all-MiniLM-L6-v2"  # 使用轻量级的开源模型
            print(f"使用模型: {model_name}")
            embeddings = HuggingFaceEmbeddings(model_name=model_name)
            print("模型加载成功")
            
            # 尝试加载现有的向量存储，如果不存在则创建新的
            try:
                vectorstore = FAISS.load_local(vector_store_path, embeddings)
                print(f"从 {vector_store_path} 加载向量存储")
            except Exception as e:
                print(f"加载向量存储失败: {str(e)}")
                print(f"创建新的向量存储到 {vector_store_path}")
                vectorstore = FAISS.from_texts(texts, embeddings)
                # 保存向量存储
                vectorstore.save_local(vector_store_path)
                print("向量存储创建并保存成功")
        except Exception as e:
            print(f"向量存储操作失败: {str(e)}")
            raise Exception(f"向量存储操作失败: {str(e)}")
        
        # 使用LLM生成代码审查建议
        try:
            # llm = OpenAI(temperature=0)
            #prompt = PromptTemplate(
            #    input_variables=["code_changes", "context"],
            #    template=CODE_REVIEW_TEMPLATE
            #)
            # chain = LLMChain(llm=llm, prompt=prompt)
            
            # 构建最优查询
            optimal_query = build_optimal_query(diff_content, query)
            print(f"构建的最优查询: {optimal_query}")
            
            # 获取相关上下文
            context_parts = []
            if optimal_query:
                # 使用最优查询获取相关上下文
                docs = vectorstore.similarity_search(optimal_query)
                context_parts.append("\n".join([doc.page_content for doc in docs]))
            
            # 添加变更文件的完整内容作为上下文
            for file_path in changed_files:
                if file_path in branch_files:
                    context_parts.append(f"文件 {file_path} 的完整内容:\n{branch_files[file_path]}")
            
            context = "\n\n".join(context_parts)
            
            # 生成审查结果
            # result = chain.run(
            #     code_changes=diff_content,
            #     context=context
            # )
            
            # 将代码变更和上下文作为字典返回
            result_dict = {
                "code_changes": diff_content,
                "context": context,
                "optimal_query": optimal_query  # 添加最优查询到返回结果中
            }
            
            return result_dict
        except Exception as e:
            print(f"生成审查结果失败: {str(e)}")
            raise Exception(f"生成审查结果失败: {str(e)}")
        
    finally:
        # 清理临时目录
        try:
            shutil.rmtree(temp_dir)
            print(f"清理临时目录: {temp_dir}")
        except Exception as e:
            print(f"清理临时目录失败: {str(e)}")