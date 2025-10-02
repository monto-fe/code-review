"""
AI 智能体管理器
"""
import asyncio
import logging
from typing import Dict, Any, Optional, List
from dataclasses import dataclass
import json

# 模型管理器将在运行时导入

logger = logging.getLogger(__name__)


@dataclass
class AgentRole:
    """智能体角色定义"""
    role_key: str
    name: str
    description: str
    system_prompt: str
    model_provider: str = "openai"
    temperature: float = 0.7
    max_tokens: int = 2048
    config: Dict[str, Any] = None
    
    def __post_init__(self):
        if self.config is None:
            self.config = {}


class AgentManager:
    """AI 智能体管理器"""
    
    def __init__(self):
        self.agents: Dict[str, AgentRole] = {}
        self.model_manager = None
        self._initialize_default_agents()
        # 初始化模型管理器
        self._initialize_model_manager()
    
    def _initialize_default_agents(self):
        """初始化默认智能体"""
        default_agents = {
            "code_reviewer": AgentRole(
                role_key="code_reviewer",
                name="代码审查专家",
                description="专业的代码审查专家，擅长代码质量分析、安全漏洞检测和最佳实践建议",
                system_prompt=(
                    "你是一位资深的代码审查专家，具有丰富的编程经验和代码质量评估能力。\n"
                    "请从以下几个方面提供专业的审查建议：\n"
                    "1. 代码质量：可读性、可维护性、性能优化\n"
                    "2. 安全性：潜在的安全漏洞、输入验证、权限控制\n"
                    "3. 最佳实践：设计模式、代码规范、架构合理性\n"
                    "4. 潜在问题：bug风险、资源泄漏、性能瓶颈\n\n"
                    "请以结构化的方式输出审查结果，包括：\n"
                    "- 总体评价\n"
                    "- 具体问题点\n"
                    "- 改进建议\n"
                    "- 代码示例（如适用）"
                ),
                model_provider=None,  # 使用配置中的默认provider
                temperature=0.3
            ),
            "architect": AgentRole(
                role_key="architect",
                name="系统架构师",
                description="系统架构设计专家，专注于系统设计、技术选型和架构优化",
                system_prompt=(
                    "你是一位经验丰富的系统架构师，擅长系统设计、技术选型、性能优化和架构演进。\n\n"
                    "请从架构角度分析代码变更：\n"
                    "1. 系统设计：模块划分、接口设计、数据流\n"
                    "2. 技术选型：技术栈合理性、性能考虑、可扩展性\n"
                    "3. 架构模式：设计模式应用、架构风格、解耦设计\n"
                    "4. 性能优化：算法复杂度、资源使用、并发处理\n"
                    "5. 可维护性：代码组织、文档完整性、测试覆盖\n\n"
                    "请提供架构层面的专业建议。"
                ),
                model_provider=None,  # 使用配置中的默认provider
                temperature=0.5
            ),
            "security_expert": AgentRole(
                role_key="security_expert",
                name="安全专家",
                description="网络安全专家，专注于代码安全漏洞检测和防护建议",
                system_prompt=(
                    "你是一位网络安全专家，擅长识别代码中的安全漏洞、潜在风险和攻击面。\n\n"
                    "请重点关注以下安全方面：\n"
                    "1. 输入验证：SQL注入、XSS、CSRF防护\n"
                    "2. 身份认证：密码安全、会话管理、多因子认证\n"
                    "3. 授权控制：权限验证、访问控制、最小权限原则\n"
                    "4. 数据保护：敏感数据加密、传输安全、存储安全\n"
                    "5. 安全配置：安全头设置、HTTPS配置、错误处理\n"
                    "6. 依赖安全：第三方库漏洞、版本管理、供应链安全\n\n"
                    "请提供详细的安全分析和防护建议。"
                ),
                model_provider=None,  # 使用配置中的默认provider
                temperature=0.2
            ),
            "qa_engineer": AgentRole(
                role_key="qa_engineer",
                name="测试工程师",
                description="质量保证专家，专注于测试用例设计和质量评估",
                system_prompt=(
                    "你是一位资深的测试工程师，擅长测试用例设计、质量评估和缺陷分析。\n\n"
                    "请从测试角度分析代码变更：\n"
                    "1. 功能测试：需求覆盖、边界条件、异常处理\n"
                    "2. 单元测试：测试覆盖率、测试用例设计、Mock使用\n"
                    "3. 集成测试：接口测试、数据流测试、系统集成\n"
                    "4. 性能测试：响应时间、并发处理、资源使用\n"
                    "5. 安全测试：漏洞扫描、渗透测试、安全配置\n"
                    "6. 用户体验：界面友好性、错误提示、操作流程\n\n"
                    "请提供测试建议和测试用例设计。"
                ),
                model_provider=None,  # 使用配置中的默认provider
                temperature=0.4
            ),
            "intent_router": AgentRole(
                role_key="intent_router",
                name="意图识别路由器",
                description="智能意图识别和Agent调度器，分析用户意图并路由到最合适的专业Agent或多个Agent协作",
                system_prompt=(
                    "你是一个智能意图识别路由器，负责分析用户输入并决定调用哪个或哪些专业Agent。\n\n"
                    "可用的专业Agent：\n"
                    "1. code_reviewer - 代码审查专家：代码质量、安全漏洞、最佳实践\n"
                    "2. architect - 系统架构师：系统设计、技术选型、架构优化\n"
                    "3. security_expert - 安全专家：安全漏洞检测、防护建议\n"
                    "4. qa_engineer - 测试工程师：测试用例设计、质量评估\n"
                    "5. general_assistant - 通用助手：一般编程问题、技术支持\n\n"
                    "分析规则：\n"
                    "- 包含'安全'、'漏洞'、'攻击'、'防护'、'加密'、'SQL注入'、'XSS'、'CSRF'等安全相关关键词 → security_expert\n"
                    "- 包含'架构'、'设计'、'系统'、'模块'、'接口'、'微服务'等关键词 → architect\n"
                    "- 包含'测试'、'用例'、'覆盖'、'验证'、'质量保证'等关键词 → qa_engineer\n"
                    "- 包含'审查'、'review'、'代码质量'、'bug'、'问题'、'优化'、'性能'等关键词 → code_reviewer\n"
                    "- 包含'全面'、'综合'、'多角度'、'协作'、'团队'等关键词 → 多个Agent协作\n"
                    "- 其他一般编程问题 → general_assistant\n\n"
                    "协作场景：\n"
                    "- 代码审查 + 安全检查：code_reviewer + security_expert\n"
                    "- 架构设计 + 测试策略：architect + qa_engineer\n"
                    "- 全面代码分析：code_reviewer + security_expert + architect + qa_engineer\n"
                    "- 安全架构设计：security_expert + architect\n"
                    "- 性能优化 + 测试：code_reviewer + qa_engineer\n\n"
                    "请分析用户输入，返回JSON格式：\n"
                    "{\n"
                    "  \"intent\": \"识别的意图类型\",\n"
                    "  \"target_agents\": [\"推荐的Agent角色键列表\"],\n"
                    "  \"collaboration_mode\": \"single\" | \"sequential\" | \"parallel\",\n"
                    "  \"confidence\": 0.95,\n"
                    "  \"reasoning\": \"选择理由\"\n"
                    "}\n\n"
                    "协作模式说明：\n"
                    "- single: 单个Agent处理\n"
                    "- sequential: 多个Agent按顺序处理，后一个Agent基于前一个的结果\n"
                    "- parallel: 多个Agent并行处理，最后合并结果"
                ),
                model_provider=None,
                temperature=0.1  # 低温度确保稳定的意图识别
            ),
            "general_assistant": AgentRole(
                role_key="general_assistant",
                name="通用助手",
                description="通用AI助手，提供各种编程和技术支持",
                system_prompt=(
                    "你是一位专业的编程助手，能够提供各种编程和技术支持。\n\n"
                    "你可以帮助用户：\n"
                    "1. 代码解释和优化\n"
                    "2. 技术问题解答\n"
                    "3. 最佳实践建议\n"
                    "4. 工具使用指导\n"
                    "5. 学习资源推荐\n\n"
                    "请以专业、友好的方式提供帮助。"
                ),
                model_provider=None,  # 使用配置中的默认provider
                temperature=0.7
            )
        }
        
        for role_key, agent in default_agents.items():
            self.register_agent(agent)
    
    def _initialize_model_manager(self):
        """初始化模型管理器"""
        try:
            from services.ai.model_manager import get_model_manager
            self.model_manager = get_model_manager()
            logger.info("模型管理器初始化完成")
        except Exception as e:
            logger.warning(f"模型管理器初始化失败: {e}")
            self.model_manager = None
    
    def register_agent(self, agent: AgentRole):
        """注册智能体"""
        self.agents[agent.role_key] = agent
        logger.info(f"智能体注册成功: {agent.name} ({agent.role_key})")
    
    def get_agent(self, role_key: str) -> Optional[AgentRole]:
        """获取智能体"""
        return self.agents.get(role_key)
    
    def list_agents(self) -> List[AgentRole]:
        """列出所有智能体"""
        return list(self.agents.values())
    
    async def chat_with_agent(
        self, 
        role_key: str, 
        message: str, 
        context: Dict[str, Any] = None,
        **kwargs
    ) -> str:
        """与指定智能体对话"""
        agent = self.get_agent(role_key)
        if not agent:
            raise ValueError(f"未找到智能体: {role_key}")
        
        # 构建完整的系统提示
        system_prompt = agent.system_prompt
        if context:
            context_str = self._format_context(context)
            system_prompt += f"\n\n当前上下文:\n{context_str}"
        
        try:
            # 确保模型管理器已初始化
            if self.model_manager is None:
                self._initialize_model_manager()
            
            if self.model_manager is None:
                raise ValueError("模型管理器初始化失败")
            
            # 如果智能体没有指定provider，使用None让模型管理器使用默认配置
            provider = agent.model_provider if agent.model_provider else None
            response = await self.model_manager.chat(
                message=message,
                system_prompt=system_prompt,
                provider=provider,
                temperature=kwargs.get('temperature', agent.temperature),
                max_tokens=kwargs.get('max_tokens', agent.max_tokens)
            )
            return response
        except Exception as e:
            logger.error(f"智能体 {role_key} 对话失败: {e}")
            raise
    
    async def code_review_with_agent(
        self,
        role_key: str,
        code_changes: str,
        context: str = "",
        query: str = "",
        **kwargs
    ) -> str:
        """使用指定智能体进行代码审查"""
        agent = self.get_agent(role_key)
        if not agent:
            raise ValueError(f"未找到智能体: {role_key}")
        
        # 构建代码审查提示
        review_prompt = (
            "请对以下代码变更进行审查：\n\n"
            f"代码变更：\n{code_changes}\n\n"
            f"相关上下文：\n{context}\n\n"
            f"查询要求：\n{query if query else '无特定要求'}\n\n"
            "请提供详细的审查建议。"
        )
        
        try:
            response = await self.chat_with_agent(
                role_key=role_key,
                message=review_prompt,
                context={"code_changes": code_changes, "context": context, "query": query},
                **kwargs
            )
            return response
        except Exception as e:
            logger.error(f"智能体 {role_key} 代码审查失败: {e}")
            raise
    
    def _format_context(self, context: Dict[str, Any]) -> str:
        """格式化上下文信息"""
        formatted = []
        for key, value in context.items():
            if isinstance(value, (dict, list)):
                formatted.append(f"{key}: {json.dumps(value, ensure_ascii=False, indent=2)}")
            else:
                formatted.append(f"{key}: {value}")
        return "\n".join(formatted)
    
    async def get_agent_capabilities(self, role_key: str) -> Dict[str, Any]:
        """获取智能体能力描述"""
        agent = self.get_agent(role_key)
        if not agent:
            return {}
        
        return {
            "role_key": agent.role_key,
            "name": agent.name,
            "description": agent.description,
            "model_provider": agent.model_provider,
            "temperature": agent.temperature,
            "max_tokens": agent.max_tokens,
            "capabilities": self._extract_capabilities(agent.system_prompt)
        }
    
    async def route_intent(self, message: str, context: Dict[str, Any] = None) -> Dict[str, Any]:
        """意图识别和Agent路由"""
        try:
            # 使用意图识别路由器分析用户输入
            intent_response = await self.chat_with_agent(
                role_key="intent_router",
                message=message,
                context=context or {},
                temperature=0.1
            )
            
            # 解析意图识别结果
            try:
                import json
                intent_data = json.loads(intent_response)
                
                # 验证返回的数据结构
                required_keys = ["intent", "target_agents", "collaboration_mode", "confidence", "reasoning"]
                if not all(key in intent_data for key in required_keys):
                    raise ValueError("意图识别返回数据格式不完整")
                
                # 确保目标Agents存在
                target_agents = intent_data["target_agents"]
                if not isinstance(target_agents, list) or len(target_agents) == 0:
                    raise ValueError("target_agents 必须是非空列表")
                
                # 验证每个Agent是否存在
                valid_agents = []
                for agent_key in target_agents:
                    if self.get_agent(agent_key):
                        valid_agents.append(agent_key)
                    else:
                        logger.warning(f"意图识别推荐了不存在的Agent: {agent_key}")
                
                # 如果没有有效Agent，回退到通用助手
                if not valid_agents:
                    valid_agents = ["general_assistant"]
                    intent_data["target_agents"] = valid_agents
                    intent_data["collaboration_mode"] = "single"
                    intent_data["reasoning"] = "没有有效的Agent，回退到通用助手"
                
                # 更新target_agents为有效列表
                intent_data["target_agents"] = valid_agents
                
                return intent_data
                
            except (json.JSONDecodeError, ValueError) as e:
                logger.warning(f"意图识别结果解析失败: {e}，使用默认路由")
                return {
                    "intent": "general_query",
                    "target_agents": ["general_assistant"],
                    "collaboration_mode": "single",
                    "confidence": 0.5,
                    "reasoning": "意图识别解析失败，使用默认路由"
                }
                
        except Exception as e:
            logger.error(f"意图识别失败: {e}，使用默认路由")
            return {
                "intent": "general_query",
                "target_agents": ["general_assistant"],
                "collaboration_mode": "single",
                "confidence": 0.3,
                "reasoning": f"意图识别失败: {str(e)}，使用默认路由"
            }
    
    async def smart_chat(self, message: str, context: Dict[str, Any] = None, **kwargs) -> Dict[str, Any]:
        """智能聊天 - 自动识别意图并路由到合适的Agent或多个Agent协作"""
        try:
            # 第一步：意图识别
            intent_result = await self.route_intent(message, context)
            target_agents = intent_result["target_agents"]
            collaboration_mode = intent_result["collaboration_mode"]
            
            logger.info(f"意图识别结果: {intent_result['intent']} -> {target_agents} (模式: {collaboration_mode}, 置信度: {intent_result['confidence']})")
            
            # 第二步：根据协作模式调用Agent(s)
            if collaboration_mode == "single" or len(target_agents) == 1:
                # 单个Agent处理
                target_agent = target_agents[0]
                response_text = await self.chat_with_agent(
                    role_key=target_agent,
                    message=message,
                    context=context,
                    **kwargs
                )
                
                return {
                    "response": response_text,
                    "intent": intent_result["intent"],
                    "target_agents": target_agents,
                    "collaboration_mode": collaboration_mode,
                    "confidence": intent_result["confidence"],
                    "reasoning": intent_result["reasoning"],
                    "agent_names": [self.get_agent(agent).name for agent in target_agents if self.get_agent(agent)],
                    "responses": {agent: response_text for agent in target_agents}
                }
                
            elif collaboration_mode == "sequential":
                # 顺序协作：多个Agent按顺序处理
                responses = {}
                current_message = message
                current_context = context or {}
                
                for i, agent_key in enumerate(target_agents):
                    logger.info(f"顺序协作 - Agent {i+1}/{len(target_agents)}: {agent_key}")
                    
                    # 为后续Agent添加上一个Agent的响应作为上下文
                    if i > 0:
                        current_context = current_context.copy()
                        current_context[f"previous_agent_{i-1}_response"] = responses[target_agents[i-1]]
                        current_message = f"基于前一个专家的分析，请继续提供专业建议：\n{current_message}"
                    
                    response = await self.chat_with_agent(
                        role_key=agent_key,
                        message=current_message,
                        context=current_context,
                        **kwargs
                    )
                    responses[agent_key] = response
                
                # 合并所有响应
                combined_response = self._combine_responses(responses, "sequential")
                
                return {
                    "response": combined_response,
                    "intent": intent_result["intent"],
                    "target_agents": target_agents,
                    "collaboration_mode": collaboration_mode,
                    "confidence": intent_result["confidence"],
                    "reasoning": intent_result["reasoning"],
                    "agent_names": [self.get_agent(agent).name for agent in target_agents if self.get_agent(agent)],
                    "responses": responses
                }
                
            elif collaboration_mode == "parallel":
                # 并行协作：多个Agent同时处理
                logger.info(f"并行协作 - 启动 {len(target_agents)} 个Agent")
                
                # 创建并行任务
                tasks = []
                for agent_key in target_agents:
                    task = self.chat_with_agent(
                        role_key=agent_key,
                        message=message,
                        context=context,
                        **kwargs
                    )
                    tasks.append((agent_key, task))
                
                # 等待所有任务完成
                responses = {}
                for agent_key, task in tasks:
                    try:
                        response = await task
                        responses[agent_key] = response
                    except Exception as e:
                        logger.error(f"Agent {agent_key} 并行处理失败: {e}")
                        responses[agent_key] = f"Agent {agent_key} 处理失败: {str(e)}"
                
                # 合并所有响应
                combined_response = self._combine_responses(responses, "parallel")
                
                return {
                    "response": combined_response,
                    "intent": intent_result["intent"],
                    "target_agents": target_agents,
                    "collaboration_mode": collaboration_mode,
                    "confidence": intent_result["confidence"],
                    "reasoning": intent_result["reasoning"],
                    "agent_names": [self.get_agent(agent).name for agent in target_agents if self.get_agent(agent)],
                    "responses": responses
                }
            
            else:
                raise ValueError(f"不支持的协作模式: {collaboration_mode}")
            
        except Exception as e:
            logger.error(f"智能聊天失败: {e}")
            # 回退到通用助手
            try:
                response_text = await self.chat_with_agent(
                    role_key="general_assistant",
                    message=message,
                    context=context,
                    **kwargs
                )
                return {
                    "response": response_text,
                    "intent": "fallback",
                    "target_agents": ["general_assistant"],
                    "collaboration_mode": "single",
                    "confidence": 0.1,
                    "reasoning": f"智能聊天失败，使用通用助手: {str(e)}",
                    "agent_names": ["通用助手"],
                    "responses": {"general_assistant": response_text}
                }
            except Exception as fallback_error:
                logger.error(f"通用助手也失败: {fallback_error}")
                return {
                    "response": f"抱歉，服务暂时不可用。错误信息: {str(e)}",
                    "intent": "error",
                    "target_agents": [],
                    "collaboration_mode": "none",
                    "confidence": 0.0,
                    "reasoning": f"所有Agent都失败: {str(e)}",
                    "agent_names": [],
                    "responses": {}
                }
    
    async def smart_code_review(self, code_changes: str, context: str = "", query: str = "", **kwargs) -> Dict[str, Any]:
        """智能代码审查 - 根据代码内容自动选择最合适的审查Agent"""
        try:
            # 分析代码内容，决定使用哪个审查Agent
            analysis_prompt = f"""
            请分析以下代码变更，决定使用哪个专业Agent进行审查：
            
            代码变更：
            {code_changes}
            
            上下文：
            {context}
            
            查询要求：
            {query if query else '无特定要求'}
            
            可用的审查Agent：
            1. code_reviewer - 通用代码审查
            2. security_expert - 安全相关代码
            3. architect - 架构相关代码
            4. qa_engineer - 测试相关代码
            
            请返回JSON格式：
            {{
                "target_agent": "推荐的Agent",
                "reasoning": "选择理由"
            }}
            """
            
            # 使用意图路由器分析
            intent_result = await self.route_intent(analysis_prompt)
            target_agent = intent_result["target_agent"]
            
            # 如果推荐的不是审查相关Agent，默认使用code_reviewer
            if target_agent not in ["code_reviewer", "security_expert", "architect", "qa_engineer"]:
                target_agent = "code_reviewer"
                intent_result["reasoning"] = "非审查Agent，使用通用代码审查专家"
            
            logger.info(f"智能代码审查选择: {target_agent} - {intent_result['reasoning']}")
            
            # 执行代码审查
            review_text = await self.code_review_with_agent(
                role_key=target_agent,
                code_changes=code_changes,
                context=context,
                query=query,
                **kwargs
            )
            
            return {
                "review": review_text,
                "target_agent": target_agent,
                "agent_name": self.get_agent(target_agent).name if self.get_agent(target_agent) else "未知",
                "reasoning": intent_result["reasoning"],
                "confidence": intent_result["confidence"]
            }
            
        except Exception as e:
            logger.error(f"智能代码审查失败: {e}")
            # 回退到通用代码审查
            try:
                review_text = await self.code_review_with_agent(
                    role_key="code_reviewer",
                    code_changes=code_changes,
                    context=context,
                    query=query,
                    **kwargs
                )
                return {
                    "review": review_text,
                    "target_agent": "code_reviewer",
                    "agent_name": "代码审查专家",
                    "reasoning": f"智能选择失败，使用通用代码审查: {str(e)}",
                    "confidence": 0.3
                }
            except Exception as fallback_error:
                logger.error(f"通用代码审查也失败: {fallback_error}")
                return {
                    "review": f"抱歉，代码审查服务暂时不可用。错误信息: {str(e)}",
                    "target_agent": "none",
                    "agent_name": "无",
                    "reasoning": f"所有审查Agent都失败: {str(e)}",
                    "confidence": 0.0
                }

    def _combine_responses(self, responses: Dict[str, str], mode: str) -> str:
        """合并多个Agent的响应"""
        if not responses:
            return "没有可用的响应"
        
        if len(responses) == 1:
            return list(responses.values())[0]
        
        combined = []
        
        if mode == "sequential":
            # 顺序协作：按顺序展示每个专家的分析
            combined.append("## 专家协作分析报告\n")
            for i, (agent_key, response) in enumerate(responses.items(), 1):
                agent_name = self.get_agent(agent_key).name if self.get_agent(agent_key) else agent_key
                combined.append(f"### {i}. {agent_name} 的分析")
                combined.append(response)
                combined.append("")  # 空行分隔
                
        elif mode == "parallel":
            # 并行协作：同时展示所有专家的观点
            combined.append("## 多专家协作分析报告\n")
            combined.append("以下是各位专家的专业分析：\n")
            
            for agent_key, response in responses.items():
                agent_name = self.get_agent(agent_key).name if self.get_agent(agent_key) else agent_key
                combined.append(f"### 🔍 {agent_name}")
                combined.append(response)
                combined.append("")  # 空行分隔
            
            # 添加综合总结
            combined.append("### 📋 综合总结")
            combined.append("基于以上各位专家的分析，建议综合考虑各个方面的建议，制定最适合的解决方案。")
            
        else:
            # 默认模式：简单合并
            for agent_key, response in responses.items():
                agent_name = self.get_agent(agent_key).name if self.get_agent(agent_key) else agent_key
                combined.append(f"**{agent_name}**: {response}")
                combined.append("")
        
        return "\n".join(combined)

    def _extract_capabilities(self, system_prompt: str) -> List[str]:
        """从系统提示中提取能力描述"""
        capabilities = []
        
        # 简单的关键词匹配来提取能力
        capability_keywords = {
            "代码质量": ["代码质量", "可读性", "可维护性", "性能"],
            "安全性": ["安全", "漏洞", "攻击", "防护"],
            "架构设计": ["架构", "设计", "模式", "解耦"],
            "测试": ["测试", "用例", "质量", "覆盖"],
            "性能": ["性能", "优化", "并发", "资源"],
            "最佳实践": ["最佳实践", "规范", "标准", "建议"]
        }
        
        for capability, keywords in capability_keywords.items():
            if any(keyword in system_prompt for keyword in keywords):
                capabilities.append(capability)
        
        return capabilities


# 全局智能体管理器实例（延迟初始化）
_agent_manager_instance = None

def get_agent_manager():
    """获取智能体管理器单例"""
    global _agent_manager_instance
    if _agent_manager_instance is None:
        _agent_manager_instance = AgentManager()
    return _agent_manager_instance

# 为了向后兼容，保留 agent_manager 变量
agent_manager = get_agent_manager()
