"""
AI 模型管理器 - 简化版本
"""
import asyncio
import logging
from typing import Dict, Any, Optional, List, Union
from abc import ABC, abstractmethod
import json
import os

# OpenAI
try:
    import openai
    from openai import AsyncOpenAI
    OPENAI_AVAILABLE = True
except ImportError:
    OPENAI_AVAILABLE = False

# 通义千问
try:
    import dashscope
    from dashscope import Generation
    DASHSCOPE_AVAILABLE = True
except ImportError:
    DASHSCOPE_AVAILABLE = False

# 本地模型
try:
    from transformers import AutoTokenizer, AutoModelForCausalLM, pipeline
    import torch
    TRANSFORMERS_AVAILABLE = True
except ImportError:
    TRANSFORMERS_AVAILABLE = False

logger = logging.getLogger(__name__)


class BaseAIModel(ABC):
    """AI 模型基类"""
    
    def __init__(self, model_name: str, **kwargs):
        self.model_name = model_name
        self.config = kwargs
    
    @abstractmethod
    async def generate(self, messages: List[Dict[str, str]], **kwargs) -> str:
        """生成文本"""
        pass
    
    @abstractmethod
    async def chat(self, message: str, system_prompt: str = "", **kwargs) -> str:
        """聊天对话"""
        pass


class OpenAIModel(BaseAIModel):
    """OpenAI 模型"""
    
    def __init__(self, model_name: str, api_key: str, base_url: str = None, **kwargs):
        super().__init__(model_name, **kwargs)
        if not OPENAI_AVAILABLE:
            raise ImportError("OpenAI 库未安装")
        
        self.client = AsyncOpenAI(
            api_key=api_key,
            base_url=base_url or os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1")
        )
    
    async def generate(self, messages: List[Dict[str, str]], **kwargs) -> str:
        """生成文本"""
        try:
            response = await self.client.chat.completions.create(
                model=self.model_name,
                messages=messages,
                max_tokens=kwargs.get('max_tokens', int(os.getenv("MAX_TOKENS", "2048"))),
                temperature=kwargs.get('temperature', float(os.getenv("TEMPERATURE", "0.7"))),
                top_p=kwargs.get('top_p', float(os.getenv("TOP_P", "0.9")))
            )
            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"OpenAI 生成失败: {e}")
            raise
    
    async def chat(self, message: str, system_prompt: str = "", **kwargs) -> str:
        """聊天对话"""
        messages = []
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        messages.append({"role": "user", "content": message})
        
        return await self.generate(messages, **kwargs)


class DashScopeModel(BaseAIModel):
    """通义千问模型"""
    
    def __init__(self, model_name: str, api_key: str, **kwargs):
        super().__init__(model_name, **kwargs)
        if not DASHSCOPE_AVAILABLE:
            raise ImportError("DashScope 库未安装")
        
        dashscope.api_key = api_key
    
    async def generate(self, messages: List[Dict[str, str]], **kwargs) -> str:
        """生成文本"""
        try:
            # 转换消息格式
            prompt = self._format_messages(messages)
            
            # 确保参数类型正确
            max_tokens = kwargs.get('max_tokens', 2048)
            if isinstance(max_tokens, str):
                max_tokens = int(max_tokens)
            
            temperature = kwargs.get('temperature', 0.7)
            if isinstance(temperature, str):
                temperature = float(temperature)
                
            top_p = kwargs.get('top_p', 0.9)
            if isinstance(top_p, str):
                top_p = float(top_p)
            
            response = Generation.call(
                model=self.model_name,
                prompt=prompt,
                max_tokens=max_tokens,
                temperature=temperature,
                top_p=top_p,
                result_format='message'
            )
            
            if response.status_code == 200:
                # DashScope API 响应格式：response.output.choices[0].message.content
                if hasattr(response.output, 'choices') and response.output.choices:
                    return response.output.choices[0].message.content
                elif hasattr(response.output, 'text') and response.output.text:
                    return response.output.text
                else:
                    raise Exception("DashScope API 响应格式异常：无法获取文本内容")
            else:
                raise Exception(f"DashScope API 错误: {response.message}")
        except Exception as e:
            logger.error(f"DashScope 生成失败: {e}")
            raise
    
    async def chat(self, message: str, system_prompt: str = "", **kwargs) -> str:
        """聊天对话"""
        messages = []
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        messages.append({"role": "user", "content": message})
        
        return await self.generate(messages, **kwargs)
    
    def _format_messages(self, messages: List[Dict[str, str]]) -> str:
        """格式化消息为通义千问格式"""
        formatted = []
        for msg in messages:
            role = msg["role"]
            content = msg["content"]
            if role == "system":
                formatted.append(f"系统: {content}")
            elif role == "user":
                formatted.append(f"用户: {content}")
            elif role == "assistant":
                formatted.append(f"助手: {content}")
        return "\n".join(formatted)


class LocalModel(BaseAIModel):
    """本地模型"""
    
    def __init__(self, model_name: str, model_path: str = None, **kwargs):
        super().__init__(model_name, **kwargs)
        if not TRANSFORMERS_AVAILABLE:
            raise ImportError("Transformers 库未安装")
        
        self.model_path = model_path or model_name
        self.tokenizer = None
        self.model = None
        self.pipeline = None
        self._load_model()
    
    def _load_model(self):
        """加载本地模型"""
        try:
            logger.info(f"加载本地模型: {self.model_path}")
            self.tokenizer = AutoTokenizer.from_pretrained(self.model_path)
            self.model = AutoModelForCausalLM.from_pretrained(
                self.model_path,
                torch_dtype=torch.float16 if torch.cuda.is_available() else torch.float32,
                device_map="auto" if torch.cuda.is_available() else "cpu"
            )
            
            self.pipeline = pipeline(
                "text-generation",
                model=self.model,
                tokenizer=self.tokenizer,
                max_length=int(os.getenv("MAX_TOKENS", "2048")),
                temperature=float(os.getenv("TEMPERATURE", "0.7")),
                do_sample=True
            )
            logger.info("本地模型加载完成")
        except Exception as e:
            logger.error(f"本地模型加载失败: {e}")
            raise
    
    async def generate(self, messages: List[Dict[str, str]], **kwargs) -> str:
        """生成文本"""
        try:
            # 转换消息格式
            prompt = self._format_messages(messages)
            
            # 在异步环境中运行同步的生成过程
            loop = asyncio.get_event_loop()
            result = await loop.run_in_executor(
                None,
                lambda: self.pipeline(
                    prompt,
                    max_length=kwargs.get('max_tokens', int(os.getenv("MAX_TOKENS", "2048"))),
                    temperature=kwargs.get('temperature', float(os.getenv("TEMPERATURE", "0.7"))),
                    top_p=kwargs.get('top_p', float(os.getenv("TOP_P", "0.9"))),
                    num_return_sequences=1
                )
            )
            
            return result[0]['generated_text'][len(prompt):].strip()
        except Exception as e:
            logger.error(f"本地模型生成失败: {e}")
            raise
    
    async def chat(self, message: str, system_prompt: str = "", **kwargs) -> str:
        """聊天对话"""
        messages = []
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        messages.append({"role": "user", "content": message})
        
        return await self.generate(messages, **kwargs)
    
    def _format_messages(self, messages: List[Dict[str, str]]) -> str:
        """格式化消息"""
        formatted = []
        for msg in messages:
            role = msg["role"]
            content = msg["content"]
            if role == "system":
                formatted.append(f"<|system|>\n{content}")
            elif role == "user":
                formatted.append(f"<|user|>\n{content}")
            elif role == "assistant":
                formatted.append(f"<|assistant|>\n{content}")
        formatted.append("<|assistant|>\n")
        return "\n".join(formatted)


class AIModelManager:
    """AI 模型管理器"""
    
    def __init__(self):
        self.models: Dict[str, BaseAIModel] = {}
        self._initialize_models()
    
    def _initialize_models(self):
        """初始化模型"""
        # OpenAI 模型
        openai_api_key = os.getenv("OPENAI_API_KEY")
        if openai_api_key and OPENAI_AVAILABLE:
            try:
                openai_model = OpenAIModel(
                    model_name=os.getenv("OPENAI_MODEL", "gpt-3.5-turbo"),
                    api_key=openai_api_key,
                    base_url=os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1")
                )
                self.models["openai"] = openai_model
                logger.info("OpenAI 模型初始化完成")
            except Exception as e:
                logger.warning(f"OpenAI 模型初始化失败: {e}")
        
        # 通义千问模型
        dashscope_api_key = os.getenv("DASHSCOPE_API_KEY")
        if dashscope_api_key and DASHSCOPE_AVAILABLE:
            try:
                dashscope_model = DashScopeModel(
                    model_name=os.getenv("DASHSCOPE_MODEL", "qwen-turbo"),
                    api_key=dashscope_api_key
                )
                self.models["dashscope"] = dashscope_model
                logger.info("DashScope 模型初始化完成")
            except Exception as e:
                logger.warning(f"DashScope 模型初始化失败: {e}")
        
        # 本地模型
        local_model_path = os.getenv("LOCAL_MODEL_PATH")
        if local_model_path and TRANSFORMERS_AVAILABLE:
            try:
                local_model = LocalModel(
                    model_name=os.getenv("LOCAL_MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2"),
                    model_path=local_model_path
                )
                self.models["local"] = local_model
                logger.info("本地模型初始化完成")
            except Exception as e:
                logger.warning(f"本地模型初始化失败: {e}")
    
    def get_model(self, provider: str = None) -> BaseAIModel:
        """获取模型实例"""
        # 如果指定了provider，直接返回
        if provider and provider in self.models:
            return self.models[provider]
        
        # 如果没有指定provider，使用配置中的默认provider
        if not provider:
            from services.config_manager import get_ai_config
            ai_config = get_ai_config()
            provider = ai_config.get("provider", "dashscope")
            if provider in self.models:
                return self.models[provider]
        
        # 按优先级返回可用模型
        for model_name in ["openai", "dashscope", "local"]:
            if model_name in self.models:
                return self.models[model_name]
        
        raise ValueError("没有可用的 AI 模型")
    
    async def generate(self, messages: List[Dict[str, str]], provider: str = None, **kwargs) -> str:
        """生成文本"""
        model = self.get_model(provider)
        return await model.generate(messages, **kwargs)
    
    async def chat(self, message: str, system_prompt: str = "", provider: str = None, **kwargs) -> str:
        """聊天对话"""
        model = self.get_model(provider)
        if model:
            return await model.chat(message, system_prompt, **kwargs)
        else:
            raise ValueError("没有可用的 AI 模型")
    
    def list_models(self) -> List[str]:
        """列出可用的模型"""
        return list(self.models.keys())


# 全局模型管理器实例（单例模式）
_model_manager_instance = None

def get_model_manager():
    """获取模型管理器单例"""
    global _model_manager_instance
    if _model_manager_instance is None:
        _model_manager_instance = AIModelManager()
    return _model_manager_instance

# 为了向后兼容，保留 model_manager 变量
model_manager = get_model_manager()
