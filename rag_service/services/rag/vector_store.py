"""
RAG 向量存储服务
"""
import os
import hashlib
import logging
from typing import List, Dict, Any, Optional, Tuple
import asyncio
from pathlib import Path
from abc import ABC, abstractmethod

# FAISS
try:
    import faiss
    import numpy as np
    from langchain_community.vectorstores import FAISS
    from langchain_huggingface import HuggingFaceEmbeddings
    FAISS_AVAILABLE = True
except ImportError:
    FAISS_AVAILABLE = False

# Chroma
try:
    import chromadb
    from chromadb.config import Settings
    CHROMA_AVAILABLE = True
except ImportError:
    CHROMA_AVAILABLE = False

from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.schema import Document
from services.config_manager import get_ai_config

logger = logging.getLogger(__name__)


class BaseVectorStore(ABC):
    """向量存储基类"""
    
    def __init__(self, collection_name: str, **kwargs):
        self.collection_name = collection_name
        self.config = kwargs
    
    @abstractmethod
    async def add_documents(self, documents: List[Document]) -> List[str]:
        """添加文档"""
        pass
    
    @abstractmethod
    async def similarity_search(
        self, 
        query: str, 
        k: int = 4, 
        filter: Dict[str, Any] = None
    ) -> List[Document]:
        """相似度搜索"""
        pass
    
    @abstractmethod
    async def delete_documents(self, ids: List[str]) -> bool:
        """删除文档"""
        pass
    
    @abstractmethod
    async def get_document_count(self) -> int:
        """获取文档数量"""
        pass


class FAISSVectorStore(BaseVectorStore):
    """FAISS 向量存储"""
    
    def __init__(self, collection_name: str, embedding_model: str = None, **kwargs):
        super().__init__(collection_name, **kwargs)
        if not FAISS_AVAILABLE:
            raise ImportError("FAISS 库未安装")
        
        ai_config = get_ai_config()
        self.embedding_model = embedding_model or ai_config.get("vector_store_embedding_model", "sentence-transformers/all-MiniLM-L6-v2")
        self.embeddings = None  # 延迟初始化
        self.vectorstore = None
        self._load_or_create_store()
    
    def _get_embeddings(self):
        """获取嵌入模型（延迟初始化）"""
        if self.embeddings is None:
            self.embeddings = HuggingFaceEmbeddings(model_name=self.embedding_model)
        return self.embeddings
    
    def _get_store_path(self) -> str:
        """获取存储路径"""
        ai_config = get_ai_config()
        vector_store_path = ai_config["vector_store_path"]
        
        # 如果是相对路径，转换为绝对路径
        if not os.path.isabs(vector_store_path):
            # 获取项目根目录 - 从当前文件位置计算
            # services/rag/vector_store.py -> 项目根目录
            project_root = Path(__file__).parent.parent.parent
            vector_store_path = project_root / vector_store_path
        
        store_dir = Path(vector_store_path) / self.collection_name
        store_dir.mkdir(parents=True, exist_ok=True)
        return str(store_dir)
    
    def _load_or_create_store(self):
        """加载或创建向量存储"""
        store_path = self._get_store_path()
        
        # 检查存储文件是否存在
        faiss_file = Path(store_path) / "index.faiss"
        pkl_file = Path(store_path) / "index.pkl"
        
        if faiss_file.exists() and pkl_file.exists():
            try:
                # 尝试加载现有存储
                try:
                    # 新版本 FAISS 支持 allow_dangerous_deserialization 参数
                    self.vectorstore = FAISS.load_local(
                        store_path,
                        self._get_embeddings(),
                        allow_dangerous_deserialization=True
                    )
                except TypeError:
                    # 旧版本 FAISS 不支持该参数
                    self.vectorstore = FAISS.load_local(
                        store_path,
                        self._get_embeddings()
                    )
                logger.info(f"从 {store_path} 加载 FAISS 向量存储")
            except Exception as e:
                logger.warning(f"加载现有存储失败: {e}，将在添加文档时重新创建")
                self.vectorstore = None
        else:
            logger.info(f"存储文件不存在，将在添加第一个文档时创建")
            self.vectorstore = None
    
    async def add_documents(self, documents: List[Document]) -> List[str]:
        """添加文档"""
        if not documents:
            return []
        
        # 如果向量存储为空，创建新的
        if self.vectorstore is None:
            loop = asyncio.get_event_loop()
            self.vectorstore = await loop.run_in_executor(
                None,
                lambda: FAISS.from_texts(
                    texts=[doc.page_content for doc in documents],
                    embedding=self._get_embeddings(),
                    metadatas=[doc.metadata for doc in documents]
                )
            )
            logger.info(f"创建新的 FAISS 向量存储，包含 {len(documents)} 个文档")
        else:
            # 在异步环境中运行同步操作
            loop = asyncio.get_event_loop()
            ids = await loop.run_in_executor(
                None,
                lambda: self.vectorstore.add_documents(documents)
            )
        
        # 保存存储
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(
            None,
            lambda: self.vectorstore.save_local(self._get_store_path())
        )
        
        logger.info(f"添加了 {len(documents)} 个文档到 FAISS 存储")
        return list(range(len(documents)))  # 返回简单的ID列表
    
    async def similarity_search(
        self, 
        query: str, 
        k: int = 4, 
        filter: Dict[str, Any] = None
    ) -> List[Document]:
        """相似度搜索"""
        try:
            # 检查是否有数据
            if self.vectorstore is None or self.vectorstore.index.ntotal == 0:
                logger.info("向量存储为空，返回空结果")
                return []
            
            loop = asyncio.get_event_loop()
            docs = await loop.run_in_executor(
                None,
                lambda: self.vectorstore.similarity_search(query, k=k)
            )
            
            # 应用过滤器
            if filter:
                docs = self._apply_filter(docs, filter)
            
            return docs
        except Exception as e:
            logger.warning(f"相似度搜索失败: {e}")
            return []
    
    async def delete_documents(self, ids: List[str]) -> bool:
        """删除文档"""
        try:
            # FAISS 不支持按ID删除，需要重建索引
            # 这里简化处理，实际应用中需要维护ID映射
            logger.warning("FAISS 不支持按ID删除文档，需要重建索引")
            return True
        except Exception as e:
            logger.error(f"删除文档失败: {e}")
            return False
    
    async def get_document_count(self) -> int:
        """获取文档数量"""
        if self.vectorstore is None:
            return 0
        return self.vectorstore.index.ntotal
    
    def _apply_filter(self, docs: List[Document], filter: Dict[str, Any]) -> List[Document]:
        """应用过滤器"""
        filtered_docs = []
        for doc in docs:
            if self._doc_matches_filter(doc, filter):
                filtered_docs.append(doc)
        return filtered_docs
    
    def _doc_matches_filter(self, doc: Document, filter: Dict[str, Any]) -> bool:
        """检查文档是否匹配过滤器"""
        metadata = doc.meta_data
        for key, value in filter.items():
            if key not in metadata or metadata[key] != value:
                return False
        return True


class ChromaVectorStore(BaseVectorStore):
    """Chroma 向量存储"""
    
    def __init__(self, collection_name: str, embedding_model: str = None, **kwargs):
        super().__init__(collection_name, **kwargs)
        if not CHROMA_AVAILABLE:
            raise ImportError("Chroma 库未安装")
        
        ai_config = get_ai_config()
        self.embedding_model = embedding_model or ai_config.get("vector_store_embedding_model", "sentence-transformers/all-MiniLM-L6-v2")
        self.embeddings = HuggingFaceEmbeddings(model_name=self.embedding_model)
        
        # 初始化 Chroma 客户端
        ai_config = get_ai_config()
        self.client = chromadb.Client(Settings(
            persist_directory=ai_config["vector_store_path"]
        ))
        
        # 获取或创建集合
        try:
            self.collection = self.client.get_collection(collection_name)
        except ValueError:
            self.collection = self.client.create_collection(collection_name)
    
    async def add_documents(self, documents: List[Document]) -> List[str]:
        """添加文档"""
        if not documents:
            return []
        
        # 准备数据
        texts = [doc.page_content for doc in documents]
        metadatas = [doc.meta_data for doc in documents]
        ids = [str(hashlib.md5(doc.page_content.encode()).hexdigest()) for doc in documents]
        
        # 生成嵌入
        loop = asyncio.get_event_loop()
        embeddings = await loop.run_in_executor(
            None,
            lambda: self.embeddings.embed_documents(texts)
        )
        
        # 添加到集合
        self.collection.add(
            embeddings=embeddings,
            documents=texts,
            metadatas=metadatas,
            ids=ids
        )
        
        logger.info(f"添加了 {len(documents)} 个文档到 Chroma 存储")
        return ids
    
    async def similarity_search(
        self, 
        query: str, 
        k: int = 4, 
        filter: Dict[str, Any] = None
    ) -> List[Document]:
        """相似度搜索"""
        # 生成查询嵌入
        loop = asyncio.get_event_loop()
        query_embedding = await loop.run_in_executor(
            None,
            lambda: self.embeddings.embed_query(query)
        )
        
        # 搜索
        results = self.collection.query(
            query_embeddings=[query_embedding],
            n_results=k,
            where=filter
        )
        
        # 转换为 Document 对象
        documents = []
        if results['documents'] and results['documents'][0]:
            for i, doc_text in enumerate(results['documents'][0]):
                metadata = results['metadatas'][0][i] if results['metadatas'] else {}
                documents.append(Document(
                    page_content=doc_text,
                    metadata=metadata
                ))
        
        return documents
    
    async def delete_documents(self, ids: List[str]) -> bool:
        """删除文档"""
        try:
            self.collection.delete(ids=ids)
            logger.info(f"删除了 {len(ids)} 个文档")
            return True
        except Exception as e:
            logger.error(f"删除文档失败: {e}")
            return False
    
    async def get_document_count(self) -> int:
        """获取文档数量"""
        return self.collection.count()


class VectorStoreManager:
    """向量存储管理器"""
    
    def __init__(self):
        self.stores: Dict[str, BaseVectorStore] = {}
        self._initialize_stores()
    
    def _initialize_stores(self):
        """初始化向量存储"""
        ai_config = get_ai_config()
        store_type = ai_config["vector_store_type"].lower()
        
        if store_type == "faiss" and FAISS_AVAILABLE:
            self.default_store_class = FAISSVectorStore
        elif store_type == "chroma" and CHROMA_AVAILABLE:
            self.default_store_class = ChromaVectorStore
        else:
            logger.warning(f"向量存储类型 {store_type} 不可用，使用 FAISS")
            self.default_store_class = FAISSVectorStore
    
    def get_store(self, collection_name: str) -> BaseVectorStore:
        """获取向量存储实例"""
        if collection_name not in self.stores:
            self.stores[collection_name] = self.default_store_class(collection_name)
        return self.stores[collection_name]
    
    async def add_documents(
        self, 
        collection_name: str, 
        documents: List[Document]
    ) -> List[str]:
        """添加文档"""
        store = self.get_store(collection_name)
        return await store.add_documents(documents)
    
    async def similarity_search(
        self,
        collection_name: str,
        query: str,
        k: int = 4,
        filter: Dict[str, Any] = None
    ) -> List[Document]:
        """相似度搜索"""
        store = self.get_store(collection_name)
        return await store.similarity_search(query, k, filter)
    
    async def delete_documents(self, collection_name: str, ids: List[str]) -> bool:
        """删除文档"""
        store = self.get_store(collection_name)
        return await store.delete_documents(ids)
    
    async def get_document_count(self, collection_name: str) -> int:
        """获取文档数量"""
        store = self.get_store(collection_name)
        return await store.get_document_count()


# 全局向量存储管理器实例
vector_store_manager = VectorStoreManager()
