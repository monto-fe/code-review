import os
from typing import List, Optional, Dict
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.vectorstores import FAISS
from langchain_community.embeddings import HuggingFaceEmbeddings
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
        return
    
    dependencies = [
        "faiss-cpu",
        "sentence-transformers"
    ]
    
    missing_deps = []
    
    for package_name in dependencies:
        try:
            if package_name == "sentence-transformers":
                import sentence_transformers
            elif package_name == "faiss-cpu":
                import faiss
        except ImportError:
            missing_deps.append(package_name)
    
    # 只在有缺失依赖时才安装
    if missing_deps:
        for package_name in missing_deps:
            try:
                # 使用 subprocess 安装依赖
                subprocess.check_call([
                    sys.executable, 
                    "-m", 
                    "pip", 
                    "install", 
                    package_name
                ])
                
                # 验证安装
                if package_name == "sentence-transformers":
                    import sentence_transformers
                elif package_name == "faiss-cpu":
                    import faiss
            except Exception as e:
                try:
                    # 尝试使用 pip 直接安装
                    subprocess.check_call([
                        "pip", 
                        "install", 
                        package_name
                    ])
                except Exception as e2:
                    raise Exception(f"无法安装 {package_name}，请手动运行: pip install {package_name}")

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
                # 从 "a/file.js b/file.js" 格式中提取文件路径
                file_path = parts[2].replace('b/', '')
                files.append(file_path)
        elif line.startswith('---') or line.startswith('+++'):
            # 从 "--- a/file.js" 或 "+++ b/file.js" 格式中提取文件路径
            if 'a/' in line:
                file_path = line.split('a/')[-1].strip()
                if file_path not in files:
                    files.append(file_path)
            elif 'b/' in line:
                file_path = line.split('b/')[-1].strip()
                if file_path not in files:
                    files.append(file_path)
    return list(set(files))  # 去重

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
            
            # 优化克隆：只克隆指定分支，深度为1
            print(f"开始克隆分支 {branch}，深度为1")
            repo = Repo.clone_from(
                auth_git_url, 
                temp_dir,
                branch=branch,  # 只克隆指定分支
                depth=1,  # 只克隆最新一次提交
                single_branch=True  # 只克隆一个分支
            )
            print(f"克隆完成，repo: {repo}")
            
            # 验证分支是否正确
            current_branch = repo.active_branch.name
            print(f"当前分支: {current_branch}")
            if current_branch != branch:
                print(f"警告：克隆的分支 {current_branch} 与请求的分支 {branch} 不匹配")
                
        except Exception as e:
            print(f"克隆仓库失败: {str(e)}")
            # 如果单分支克隆失败，尝试完整克隆
            print("尝试完整克隆...")
            repo = Repo.clone_from(auth_git_url, temp_dir)
            # 切换到指定分支
            repo.git.checkout(branch)
            print(f"完整克隆完成，repo: {repo}")
        
        # 从diff中提取变更的文件
        changed_files = extract_files_from_diff(diff_content)
        print(f"changed_files: {changed_files}")
        
        # 只读取变更的文件，而不是整个分支
        branch_files = {}
        if changed_files:
            for file_path in changed_files:
                try:
                    # 尝试多种可能的路径
                    possible_paths = [
                        os.path.join(temp_dir, file_path),
                        os.path.join(temp_dir, file_path.replace('a/', '')),
                        os.path.join(temp_dir, file_path.replace('b/', '')),
                        os.path.join(temp_dir, os.path.basename(file_path))
                    ]
                    
                    file_found = False
                    for full_path in possible_paths:
                        if os.path.exists(full_path):
                            with open(full_path, 'r', encoding='utf-8') as f:
                                branch_files[file_path] = f.read()
                            print(f"读取文件: {file_path} (路径: {full_path})")
                            file_found = True
                            break
                    
                    if not file_found:
                        print(f"文件不存在，尝试的路径: {possible_paths}")
                        
                except Exception as e:
                    print(f"读取文件 {file_path} 失败: {str(e)}")
                    continue
        
        # 如果没有找到变更文件，尝试读取一些关键文件作为上下文
        if not branch_files:
            print("未找到变更文件，尝试读取关键文件作为上下文")
            # 读取一些常见的配置文件
            context_files = [
                'README.md', 'requirements.txt', 'pyproject.toml', 
                'package.json', 'Dockerfile', 'docker-compose.yml',
                'app.py', 'main.py', 'index.js', 'index.ts'
            ]
            
            for file_name in context_files:
                try:
                    full_path = os.path.join(temp_dir, file_name)
                    if os.path.exists(full_path):
                        with open(full_path, 'r', encoding='utf-8') as f:
                            branch_files[file_name] = f.read()
                        print(f"读取上下文文件: {file_name}")
                except Exception as e:
                    print(f"读取上下文文件 {file_name} 失败: {str(e)}")
                    continue
        
        print(f"读取到 {len(branch_files)} 个文件")
        
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
        
        # 处理文件内容用于上下文
        all_texts = []
        for file_path, content in branch_files.items():
            all_texts.append(f"文件: {file_path}\n内容:\n{content}\n")
        
        if not all_texts:
            print("警告：没有找到任何文件内容，将使用diff内容作为上下文")
            all_texts.append(f"代码变更:\n{diff_content}\n")
        
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
                # 尝试不同的加载方式
                try:
                    vectorstore = FAISS.load_local(
                        vector_store_path,
                        embeddings
                    )
                except TypeError:
                    # 如果新版本不支持 allow_dangerous_deserialization 参数，尝试旧方式
                    vectorstore = FAISS.load_local(
                        vector_store_path,
                        embeddings,
                        allow_dangerous_deserialization=True
                    )
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