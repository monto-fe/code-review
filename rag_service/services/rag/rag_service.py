"""
RAG 增强检索服务
"""
import asyncio
import logging
import hashlib
from typing import List, Dict, Any, Optional, Tuple
from pathlib import Path
import os
import shutil
import re

from git import Repo
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.schema import Document
from urllib.parse import urlparse

from services.rag.vector_store import vector_store_manager
from services.config_manager import get_rag_config

logger = logging.getLogger(__name__)


class RAGService:
    """RAG 增强检索服务"""
    
    def __init__(self):
        rag_config = get_rag_config()
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=rag_config["chunk_size"],
            chunk_overlap=rag_config["chunk_overlap"]
        )
    
    def _get_repository_id(self, git_url: str, branch: str) -> str:
        """生成仓库唯一标识"""
        return hashlib.md5(f"{git_url}:{branch}".encode()).hexdigest()
    
    def _get_git_url_with_token(self, git_url: str, gitlab_token: str) -> str:
        """获取带认证的 Git URL"""
        if git_url.startswith('git@'):
            # SSH 格式转换为 HTTP 格式
            parts = git_url.split(':', 2)
            if len(parts) == 3:
                host, port, path = parts
                host = host.replace('git@', '')
                return f"http://oauth2:{gitlab_token}@{host}:{port}/{path}"
            else:
                host = git_url.split('@')[1].split(':')[0]
                path = git_url.split(':')[1]
                return f"http://oauth2:{gitlab_token}@{host}/{path}"
        else:
            # HTTP/HTTPS 格式
            parsed_url = urlparse(git_url)
            if parsed_url.scheme in ['http', 'https']:
                netloc = parsed_url.netloc
                if ':' in netloc:
                    host, port = netloc.split(':')
                    return f"http://oauth2:{gitlab_token}@{host}:{port}{parsed_url.path}"
                else:
                    return f"http://oauth2:{gitlab_token}@{netloc}{parsed_url.path}"
            else:
                raise ValueError(f"不支持的URL格式: {git_url}")
    
    def _extract_files_from_diff(self, diff_content: str) -> List[str]:
        """从 diff 内容中提取变更的文件路径"""
        files = []
        for line in diff_content.split('\n'):
            if line.startswith('diff --git'):
                parts = line.split()
                if len(parts) >= 3:
                    file_path = parts[2].replace('b/', '')
                    files.append(file_path)
            elif line.startswith('---') or line.startswith('+++'):
                if 'a/' in line:
                    file_path = line.split('a/')[-1].strip()
                    if file_path not in files:
                        files.append(file_path)
                elif 'b/' in line:
                    file_path = line.split('b/')[-1].strip()
                    if file_path not in files:
                        files.append(file_path)
        return list(set(files))
    
    def _extract_query_from_diff(self, diff_content: str) -> str:
        """从 diff 内容中提取关键信息构建查询"""
        query_parts = []
        
        # 提取文件类型
        file_extensions = set()
        file_pattern = r'diff --git a/.*?\.(\w+)'
        extensions = re.findall(file_pattern, diff_content)
        for ext in extensions:
            if ext.lower() in ['py', 'js', 'ts', 'java', 'go', 'cpp', 'c', 'php', 'rb', 'rs']:
                file_extensions.add(ext.lower())
        
        if file_extensions:
            query_parts.append(f"编程语言: {', '.join(file_extensions)}")
        
        # 提取函数名
        function_pattern = r'\+def\s+(\w+)'
        functions = re.findall(function_pattern, diff_content)
        if functions:
            query_parts.append(f"函数: {', '.join(functions)}")
        
        # 提取类名
        class_pattern = r'\+class\s+(\w+)'
        classes = re.findall(class_pattern, diff_content)
        if classes:
            query_parts.append(f"类: {', '.join(classes)}")
        
        # 提取语义关键词
        semantic_keywords = []
        keywords_map = {
            'API接口': ['api', 'http', 'request', 'response', 'endpoint'],
            '数据库操作': ['database', 'db', 'sql', 'query', 'select', 'insert', 'update', 'delete'],
            '认证授权': ['auth', 'login', 'password', 'token', 'jwt', 'oauth'],
            '错误处理': ['error', 'exception', 'try', 'catch', 'finally'],
            '配置管理': ['config', 'setting', 'env', 'environment']
        }
        
        for category, keywords in keywords_map.items():
            if any(keyword in diff_content.lower() for keyword in keywords):
                semantic_keywords.append(category)
        
        if semantic_keywords:
            query_parts.append(f"功能类型: {', '.join(semantic_keywords)}")
        
        return " ".join(query_parts)
    
    async def process_repository(
        self, 
        git_url: str, 
        branch: str, 
        gitlab_token: str,
        force_refresh: bool = False
    ) -> str:
        """处理代码仓库，建立向量索引"""
        repository_id = self._get_repository_id(git_url, branch)
        collection_name = f"repo_{repository_id}"
        
        # 检查是否已存在且不需要刷新
        if not force_refresh:
            doc_count = await vector_store_manager.get_document_count(collection_name)
            if doc_count > 0:
                logger.info(f"仓库 {git_url}:{branch} 已存在 {doc_count} 个文档，跳过处理")
                return repository_id
        
        # 克隆仓库
        temp_dir = Path("/tmp/git_repos") / repository_id
        temp_dir.mkdir(parents=True, exist_ok=True)
        
        try:
            # 清理旧目录
            if temp_dir.exists():
                shutil.rmtree(temp_dir)
            
            # 克隆仓库
            auth_git_url = self._get_git_url_with_token(git_url, gitlab_token)
            
            # 设置 Git 环境变量
            os.environ['GIT_SSL_NO_VERIFY'] = 'true'
            os.environ['GIT_HTTP_LOW_SPEED_TIME'] = '300'
            os.environ['GIT_HTTP_LOW_SPEED_LIMIT'] = '1000'
            
            # 克隆指定分支
            repo = Repo.clone_from(
                auth_git_url,
                temp_dir,
                branch=branch,
                depth=1,
                single_branch=True
            )
            
            # 读取代码文件
            documents = []
            for root, dirs, files in os.walk(temp_dir):
                # 跳过隐藏目录和常见的不需要索引的目录
                dirs[:] = [d for d in dirs if not d.startswith('.') and d not in ['node_modules', '__pycache__', 'venv', 'env']]
                
                for file in files:
                    if self._should_index_file(file):
                        file_path = Path(root) / file
                        try:
                            with open(file_path, 'r', encoding='utf-8') as f:
                                content = f.read()
                            
                            # 创建文档
                            relative_path = file_path.relative_to(temp_dir)
                            doc = Document(
                                page_content=content,
                                metadata={
                                    "file_path": str(relative_path),
                                    "file_name": file,
                                    "file_type": file_path.suffix,
                                    "repository_id": repository_id,
                                    "git_url": git_url,
                                    "branch": branch
                                }
                            )
                            documents.append(doc)
                        except Exception as e:
                            logger.warning(f"读取文件 {file_path} 失败: {e}")
                            continue
            
            # 分割文档
            split_documents = []
            for doc in documents:
                chunks = self.text_splitter.split_documents([doc])
                split_documents.extend(chunks)
            
            # 添加到向量存储
            if split_documents:
                await vector_store_manager.add_documents(collection_name, split_documents)
                logger.info(f"成功处理仓库 {git_url}:{branch}，添加了 {len(split_documents)} 个文档块")
            else:
                logger.warning(f"仓库 {git_url}:{branch} 没有找到可索引的文件")
            
            return repository_id
            
        finally:
            # 清理临时目录
            if temp_dir.exists():
                shutil.rmtree(temp_dir)
    
    def _should_index_file(self, filename: str) -> bool:
        """判断文件是否应该被索引"""
        # 文件扩展名白名单
        indexable_extensions = {
            '.py', '.js', '.ts', '.jsx', '.tsx', '.java', '.go', '.cpp', '.c', '.h',
            '.php', '.rb', '.rs', '.swift', '.kt', '.scala', '.clj', '.hs',
            '.md', '.txt', '.json', '.yaml', '.yml', '.xml', '.html', '.css',
            '.sql', '.sh', '.bat', '.ps1'
        }
        
        # 文件名黑名单
        blacklisted_files = {
            'package-lock.json', 'yarn.lock', 'pnpm-lock.yaml',
            'go.sum', 'composer.lock', 'Pipfile.lock'
        }
        
        file_path = Path(filename)
        
        # 检查扩展名
        if file_path.suffix.lower() not in indexable_extensions:
            return False
        
        # 检查文件名
        if filename in blacklisted_files:
            return False
        
        return True
    
    async def search_relevant_documents(
        self,
        repository_id: str,
        query: str,
        k: int = 4,
        similarity_threshold: float = None
    ) -> List[Document]:
        """搜索相关文档"""
        collection_name = f"repo_{repository_id}"
        rag_config = get_rag_config()
        threshold = similarity_threshold or rag_config["similarity_threshold"]
        
        # 执行相似度搜索
        documents = await vector_store_manager.similarity_search(
            collection_name=collection_name,
            query=query,
            k=k
        )
        
        # 应用相似度阈值过滤
        filtered_documents = []
        for doc in documents:
            # 这里简化处理，实际应用中需要计算相似度分数
            filtered_documents.append(doc)
        
        return filtered_documents
    
    async def get_code_context(
        self,
        git_url: str,
        branch: str,
        diff_content: str,
        query: str = "",
        gitlab_token: str = None
    ) -> Dict[str, Any]:
        """获取代码上下文信息"""
        # 处理仓库
        repository_id = await self.process_repository(git_url, branch, gitlab_token)
        
        # 构建搜索查询
        search_query = query
        if not search_query:
            search_query = self._extract_query_from_diff(diff_content)
        
        # 获取 RAG 配置
        rag_config = get_rag_config()
        
        # 搜索相关文档
        relevant_docs = await self.search_relevant_documents(
            repository_id=repository_id,
            query=search_query,
            k=rag_config["max_context_length"] // rag_config["chunk_size"]
        )
        
        # 提取变更文件
        changed_files = self._extract_files_from_diff(diff_content)
        
        # 构建上下文
        context_parts = []
        
        # 添加相关文档
        for doc in relevant_docs:
            context_parts.append(f"文件: {doc.metadata.get('file_path', 'unknown')}\n内容:\n{doc.page_content}")
        
        # 添加变更文件信息
        if changed_files:
            context_parts.append(f"变更文件: {', '.join(changed_files)}")
        
        context = "\n\n".join(context_parts)
        
        return {
            "repository_id": repository_id,
            "search_query": search_query,
            "context": context,
            "relevant_documents": len(relevant_docs),
            "changed_files": changed_files
        }
    
    async def delete_repository(self, repository_id: str) -> bool:
        """删除仓库的向量索引"""
        collection_name = f"repo_{repository_id}"
        try:
            # 这里简化处理，实际应用中需要实现删除逻辑
            logger.info(f"删除仓库索引: {repository_id}")
            return True
        except Exception as e:
            logger.error(f"删除仓库索引失败: {e}")
            return False


# 全局 RAG 服务实例
rag_service = RAGService()
