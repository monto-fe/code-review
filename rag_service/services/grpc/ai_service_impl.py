"""
gRPC AI 服务实现
"""
import asyncio
import logging
import time
from typing import Dict, Any, List
import grpc
from concurrent import futures

# 导入生成的 gRPC 代码
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'proto'))

try:
    from proto import ai_service_pb2
    from proto import ai_service_pb2_grpc
except ImportError:
    # 如果 proto 文件还没有生成，创建占位符
    class ai_service_pb2:
        class CodeReviewRequest:
            pass
        class CodeReviewResponse:
            pass
        class ChatRequest:
            pass
        class ChatResponse:
            pass
        class DocumentSearchRequest:
            pass
        class DocumentSearchResponse:
            pass
        class Document:
            pass
        class HealthCheckRequest:
            pass
        class HealthCheckResponse:
            pass
        class ServingStatus:
            UNKNOWN = 0
            SERVING = 1
            NOT_SERVING = 2

    class ai_service_pb2_grpc:
        class AIServiceServicer:
            pass

# 这些将在运行时导入
from services.config_manager import get_service_config

logger = logging.getLogger(__name__)


class AIServiceServicer(ai_service_pb2_grpc.AIServiceServicer):
    """AI 服务 gRPC 实现"""
    
    def __init__(self):
        # 在运行时导入，确保环境变量已设置
        from services.ai.agent_manager import get_agent_manager
        from services.rag.rag_service import rag_service
        self.agent_manager = get_agent_manager()
        self.rag_service = rag_service
    
    async def CodeReview(self, request, context):
        """代码审查服务"""
        try:
            logger.info(f"收到代码审查请求: {request.git_url}:{request.branch}")
            
            # 获取智能体
            agent_role = request.agent_role or "code_reviewer"
            agent = self.agent_manager.get_agent(agent_role)
            if not agent:
                if context:
                    context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                    context.set_details(f"未找到智能体角色: {agent_role}")
                return ai_service_pb2.CodeReviewResponse()
            
            # 获取 RAG 上下文（如果可用）
            try:
                rag_context = await self.rag_service.get_code_context(
                    git_url=request.git_url,
                    branch=request.branch,
                    diff_content=request.diff_content,
                    query=request.query,
                    gitlab_token=request.gitlab_token
                )
            except Exception as e:
                logger.warning(f"RAG 上下文获取失败，使用基础模式: {e}")
                rag_context = {
                    "context": "",
                    "changed_files": [],
                    "relevant_documents": 0
                }
            
            # 构建上下文信息
            context_info = {
                "git_url": request.git_url,
                "branch": request.branch,
                "changed_files": rag_context.get("changed_files", []),
                "relevant_documents": rag_context.get("relevant_documents", 0)
            }
            
            # 执行代码审查
            start_time = time.time()
            
            # 如果指定了特定的agent_role，直接使用
            if request.agent_role and request.agent_role != "auto":
                review_result = await self.agent_manager.code_review_with_agent(
                    role_key=agent_role,
                    code_changes=request.diff_content,
                    context=rag_context.get("context", ""),
                    query=request.query,
                    temperature=agent.temperature,
                    max_tokens=agent.max_tokens
                )
                
                # 构建响应
                response = ai_service_pb2.CodeReviewResponse()
                response.review = review_result
                response.agent_role = agent_role
                response.confidence_score = 85  # 简化的置信度评分
                response.metadata["processing_time"] = str(time.time() - start_time)
                response.metadata["repository_id"] = rag_context.get("repository_id", "")
                response.metadata["search_query"] = rag_context.get("search_query", "")
                response.metadata["routing_mode"] = "direct"
                
            else:
                # 使用智能代码审查
                smart_result = await self.agent_manager.smart_code_review(
                    code_changes=request.diff_content,
                    context=rag_context.get("context", ""),
                    query=request.query
                )
                
                # 构建响应
                response = ai_service_pb2.CodeReviewResponse()
                response.review = smart_result["review"]
                response.agent_role = smart_result["target_agent"]
                response.confidence_score = int(smart_result["confidence"] * 100)
                response.metadata["processing_time"] = str(time.time() - start_time)
                response.metadata["repository_id"] = rag_context.get("repository_id", "")
                response.metadata["search_query"] = rag_context.get("search_query", "")
                response.metadata["agent_name"] = smart_result["agent_name"]
                response.metadata["reasoning"] = smart_result["reasoning"]
                response.metadata["routing_mode"] = "smart"
            
            processing_time = time.time() - start_time
            
            # 生成智能建议
            review_text = response.review
            suggestions = self._extract_suggestions_from_review(review_text)
            response.suggestions.extend(suggestions)
            
            logger.info(f"代码审查完成，耗时: {processing_time:.2f}秒")
            return response
            
        except Exception as e:
            logger.error(f"代码审查失败: {e}")
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"代码审查失败: {str(e)}")
            return ai_service_pb2.CodeReviewResponse()
    
    async def Chat(self, request, context):
        """聊天服务"""
        try:
            logger.info(f"收到聊天请求: {request.agent_role}")
            
            # 构建上下文
            context_info = {}
            if request.context:
                context_info = dict(request.context)
            
            start_time = time.time()
            
            # 如果指定了特定的agent_role，直接使用
            if request.agent_role and request.agent_role != "auto":
                agent_role = request.agent_role
                agent = self.agent_manager.get_agent(agent_role)
                if not agent:
                    if context:
                        context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                        context.set_details(f"未找到智能体角色: {agent_role}")
                    return ai_service_pb2.ChatResponse()
                
                # 使用指定的Agent
                response_text = await self.agent_manager.chat_with_agent(
                    role_key=agent_role,
                    message=request.message,
                    context=context_info,
                    temperature=agent.temperature,
                    max_tokens=agent.max_tokens
                )
                
                # 构建响应
                response = ai_service_pb2.ChatResponse()
                response.response = response_text
                response.agent_role = agent_role
                response.session_id = request.session_id
                response.metadata["processing_time"] = str(time.time() - start_time)
                response.metadata["model_provider"] = agent.model_provider or "unknown"
                response.metadata["routing_mode"] = "direct"
                
            else:
                # 使用智能路由
                smart_result = await self.agent_manager.smart_chat(
                    message=request.message,
                    context=context_info
                )
                
                # 构建响应
                response = ai_service_pb2.ChatResponse()
                response.response = smart_result["response"]
                response.agent_role = smart_result["target_agents"][0] if smart_result["target_agents"] else "unknown"
                response.session_id = request.session_id
                response.metadata["processing_time"] = str(time.time() - start_time)
                response.metadata["intent"] = smart_result["intent"]
                response.metadata["confidence"] = str(smart_result["confidence"])
                response.metadata["reasoning"] = smart_result["reasoning"]
                response.metadata["target_agents"] = ",".join(smart_result["target_agents"])
                response.metadata["agent_names"] = ",".join(smart_result["agent_names"])
                response.metadata["collaboration_mode"] = smart_result["collaboration_mode"]
                response.metadata["routing_mode"] = "smart"
            
            processing_time = time.time() - start_time
            logger.info(f"聊天完成，耗时: {processing_time:.2f}秒")
            return response
            
        except Exception as e:
            logger.error(f"聊天失败: {e}")
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"聊天失败: {str(e)}")
            return ai_service_pb2.ChatResponse()
    
    async def DocumentSearch(self, request, context):
        """文档搜索服务"""
        try:
            logger.info(f"收到文档搜索请求: {request.query}")
            
            # 搜索相关文档（如果可用）
            try:
                documents = await self.rag_service.search_relevant_documents(
                    repository_id=request.repository_id,
                    query=request.query,
                    k=request.limit or 4,
                    similarity_threshold=request.similarity_threshold or 0.7
                )
            except Exception as e:
                logger.warning(f"文档搜索失败，返回空结果: {e}")
                documents = []
            
            # 构建响应
            response = ai_service_pb2.DocumentSearchResponse()
            response.total_count = len(documents)
            
            if documents:
                response.max_similarity = 0.95  # 简化的相似度分数
                
                for doc in documents:
                    doc_proto = ai_service_pb2.Document()
                    doc_proto.id = str(hash(doc.page_content))
                    doc_proto.content = doc.page_content
                    doc_proto.file_path = doc.meta_data.get("file_path", "")
                    doc_proto.similarity_score = 0.9  # 简化的相似度分数
                    
                    # 添加元数据
                    for key, value in doc.meta_data.items():
                        doc_proto.metadata[key] = str(value)
                    
                    response.documents.append(doc_proto)
            
            logger.info(f"文档搜索完成，找到 {len(documents)} 个文档")
            return response
            
        except Exception as e:
            logger.error(f"文档搜索失败: {e}")
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(f"文档搜索失败: {str(e)}")
            return ai_service_pb2.DocumentSearchResponse()
    
    async def HealthCheck(self, request, context):
        """健康检查服务"""
        try:
            # 检查服务状态
            status = ai_service_pb2.HealthCheckResponse.SERVING
            
            # 检查数据库连接
            try:
                # 这里可以添加数据库健康检查
                pass
            except Exception:
                status = ai_service_pb2.HealthCheckResponse.NOT_SERVING
            
            # 检查 AI 模型
            available_models = []
            try:
                available_models = self.agent_manager.model_manager.list_models()
                if not available_models:
                    status = ai_service_pb2.HealthCheckResponse.NOT_SERVING
            except Exception as e:
                logger.warning(f"AI 模型检查失败: {e}")
                status = ai_service_pb2.HealthCheckResponse.NOT_SERVING
            
            # 构建响应
            response = ai_service_pb2.HealthCheckResponse()
            response.status = status
            response.message = "服务正常" if status == ai_service_pb2.HealthCheckResponse.SERVING else "服务异常"
            response.details["available_models"] = ",".join(available_models)
            response.details["service_version"] = "1.0.0"
            
            return response
            
        except Exception as e:
            logger.error(f"健康检查失败: {e}")
            response = ai_service_pb2.HealthCheckResponse()
            response.status = ai_service_pb2.HealthCheckResponse.NOT_SERVING
            response.message = f"健康检查失败: {str(e)}"
            return response

    def _extract_suggestions_from_review(self, review_text: str) -> List[str]:
        """从审查报告中提取具体建议"""
        suggestions = []
        
        # 简单的关键词匹配提取建议
        suggestion_patterns = {
            "测试相关": ["测试", "单元测试", "集成测试", "测试用例"],
            "安全相关": ["安全", "漏洞", "注入", "权限", "验证"],
            "性能相关": ["性能", "优化", "效率", "内存", "CPU"],
            "代码质量": ["重构", "可读性", "维护性", "命名", "注释"],
            "架构相关": ["架构", "设计模式", "解耦", "模块化", "接口"],
            "最佳实践": ["最佳实践", "规范", "标准", "约定", "风格"]
        }
        
        # 检查每个类别
        for category, keywords in suggestion_patterns.items():
            if any(keyword in review_text for keyword in keywords):
                suggestions.append(f"关注{category}方面的改进")
        
        # 如果没有找到具体建议，提供通用建议
        if not suggestions:
            suggestions = [
                "建议添加单元测试",
                "考虑代码重构以提高可读性",
                "检查潜在的安全漏洞"
            ]
        
        # 限制建议数量，避免过多
        return suggestions[:5]


class GRPCServer:
    """gRPC 服务器"""
    
    def __init__(self):
        self.server = None
        self.ai_service = AIServiceServicer()
    
    async def start(self):
        """启动 gRPC 服务器（阻塞模式）"""
        try:
            await self._setup_server()
            await self.server.start()
            logger.info(f"gRPC 服务器启动成功，监听地址: {self.listen_addr}")
            
            # 等待服务器关闭
            await self.server.wait_for_termination()
            
        except Exception as e:
            logger.error(f"gRPC 服务器启动失败: {e}")
            raise
    
    async def start_async(self):
        """启动 gRPC 服务器（非阻塞模式）"""
        try:
            await self._setup_server()
            await self.server.start()
            logger.info(f"gRPC 服务器启动成功，监听地址: {self.listen_addr}")
            
        except Exception as e:
            logger.error(f"gRPC 服务器启动失败: {e}")
            raise
    
    async def _setup_server(self):
        """设置 gRPC 服务器"""
        # 获取服务配置
        service_config = get_service_config()
        
        # 创建 gRPC 服务器
        self.server = grpc.aio.server(
            futures.ThreadPoolExecutor(max_workers=service_config["max_workers"]),
            options=[
                ('grpc.max_message_length', service_config["max_message_length"]),
                ('grpc.max_receive_message_length', service_config["max_message_length"]),
            ]
        )
        
        # 添加服务
        ai_service_pb2_grpc.add_AIServiceServicer_to_server(
            self.ai_service, self.server
        )
        
        # 添加反射服务（用于调试）
        if service_config["enable_reflection"]:
            from grpc_reflection.v1alpha import reflection
            SERVICE_NAMES = (
                ai_service_pb2.DESCRIPTOR.services_by_name['AIService'].full_name,
                reflection.SERVICE_NAME,
            )
            reflection.enable_server_reflection(SERVICE_NAMES, self.server)
        
        # 设置监听地址
        self.listen_addr = f"{service_config['host']}:{service_config['port']}"
        self.server.add_insecure_port(self.listen_addr)
    
    async def stop(self):
        """停止 gRPC 服务器"""
        if self.server:
            await self.server.stop(grace=5.0)
            logger.info("gRPC 服务器已停止")



# 全局 gRPC 服务器实例
grpc_server = GRPCServer()
