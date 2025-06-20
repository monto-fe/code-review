from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel, Field
from typing import Optional
from rag_engine import analyze_gitlab_code
import os
from dotenv import load_dotenv
from pathlib import Path
import asyncio
import time

# 加载上级目录的.env文件
env_path = Path(__file__).parent / '.env'
load_dotenv(dotenv_path=env_path)

# 打印env_path
print("--------------------------------")
print(env_path)
# print(os.getenv("OPENAI_API_KEY"))
print("--------------------------------")

app = FastAPI(
    title="Code Review RAG Service",
    description="基于 RAG 的代码审查服务，支持 GitLab 代码分析和审查建议生成",
    version="1.0.0"
)

# 从环境变量获取配置
# OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
GITLAB_TOKEN = os.getenv("GITLAB_TOKEN")
RAG_SERVICE_HOST = os.getenv("RAG_SERVICE_HOST", "0.0.0.0")
RAG_SERVICE_PORT = int(os.getenv("RAG_SERVICE_PORT", "8000"))
VECTOR_STORE_TYPE = os.getenv("VECTOR_STORE_TYPE", "faiss")
VECTOR_STORE_PATH = os.getenv("VECTOR_STORE_PATH", "/app/data/vector_store")
GIT_TEMP_DIR = os.getenv("GIT_TEMP_DIR", "/tmp/git_repos")

print(f"RAG_SERVICE_HOST: {RAG_SERVICE_HOST}")
print(f"RAG_SERVICE_PORT: {RAG_SERVICE_PORT}")
print(f"VECTOR_STORE_TYPE: {VECTOR_STORE_TYPE}")
print(f"VECTOR_STORE_PATH: {VECTOR_STORE_PATH}")
print(f"GIT_TEMP_DIR: {GIT_TEMP_DIR}")

class CodeReviewRequest(BaseModel):
    git_url: str = Field(..., description="Git仓库URL", example="https://gitlab.com/user/repo.git")
    branch: str = Field(..., description="分支名", example="main")
    diff_content: str = Field(..., description="格式化的diff字符串", example="diff --git a/file.py b/file.py\n@@ -1,3 +1,3 @@\n-def old_function():\n+def new_function():\n     return True")
    query: Optional[str] = Field(None, description="可选的查询参数，用于相似度搜索", example="查找数据库相关代码")
    gitlab_token: Optional[str] = Field(None, description="GitLab访问令牌", example="glpat-xxxxxxxx")

class CodeAnalysisResponse(BaseModel):
    review: str = Field(..., description="代码审查结果", example="代码审查建议：\n1. 代码质量：函数命名清晰，逻辑简单\n2. 潜在问题：无\n3. 改进建议：可以考虑添加类型注解")

# 超时控制函数
async def run_with_timeout(coro, timeout_seconds=60):
    """
    在指定时间内运行协程，超时则抛出异常
    """
    try:
        return await asyncio.wait_for(coro, timeout=timeout_seconds)
    except asyncio.TimeoutError:
        raise HTTPException(status_code=408, detail=f"请求超时({timeout_seconds}秒)")

# 同步函数包装器
def run_sync_with_timeout(func, *args, **kwargs):
    """
    在指定时间内运行同步函数，超时则抛出异常
    """
    import concurrent.futures
    import threading
    
    def target():
        try:
            result_container['result'] = func(*args, **kwargs)
        except Exception as e:
            result_container['error'] = e
    
    result_container = {}
    thread = threading.Thread(target=target)
    thread.daemon = True
    thread.start()
    thread.join(timeout=60)  # 60秒超时
    
    if thread.is_alive():
        raise HTTPException(status_code=408, detail="RAG分析超时(60秒)")
    
    if 'error' in result_container:
        raise result_container['error']
    
    return result_container['result']

@app.post("/api/code-analysis", response_model=CodeAnalysisResponse, tags=["代码审查"])
async def code_analysis(request: CodeReviewRequest):
    """
    代码审查接口
    
    基于 RAG (Retrieval-Augmented Generation) 技术对 GitLab 代码变更进行智能审查。
    
    - **git_url**: Git 仓库的 URL 地址
    - **branch**: 要分析的分支名
    - **diff_content**: 格式化的 diff 字符串，包含代码变更信息
    - **query**: 可选的查询参数，用于在代码库中搜索相关内容
    - **gitlab_token**: GitLab 访问令牌，用于认证访问仓库
    
    返回详细的代码审查建议，包括：
    - 代码质量评估
    - 潜在问题识别
    - 改进建议
    - 示例代码（如适用）
    """
    print(f"Received request: {request}")
    print(f"git_url: {request.git_url}")
    print(f"branch: {request.branch}")
    print(f"diff_content length: {len(request.diff_content)}")
    print(f"query: {request.query}")
    print(f"gitlab_token: {request.gitlab_token[:8] if request.gitlab_token else 'None'}...")
    
    try:
        # 检查OpenAI API密钥
        # if not OPENAI_API_KEY:
        #    raise HTTPException(status_code=500, detail="OpenAI API key not configured")
        
        # 获取GitLab token（优先使用请求中的token，其次使用环境变量中的token）
        gitlab_token = request.gitlab_token or GITLAB_TOKEN
        if not gitlab_token:
            raise HTTPException(status_code=500, detail="GitLab token not configured")
        
        # 使用超时控制运行RAG分析
        start_time = time.time()
        result = run_sync_with_timeout(
            analyze_gitlab_code,
            git_url=request.git_url,
            branch=request.branch,
            diff_content=request.diff_content,
            query=request.query,
            gitlab_token=gitlab_token
        )
        end_time = time.time()
        
        print(f"RAG分析完成，耗时: {end_time - start_time:.2f}秒")
        
        # 将字典结果转换为字符串
        if isinstance(result, dict):
            review_text = f"代码变更:\n{result.get('code_changes', '')}\n\n相关上下文:\n{result.get('context', '')}\n\n最优查询:\n{result.get('optimal_query', '')}"
        else:
            review_text = str(result)
        
        return CodeAnalysisResponse(review=review_text)
    except HTTPException:
        raise
    except Exception as e:
        print(f"Unexpected error: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health", tags=["健康检查"])
async def health_check():
    """
    健康检查接口
    
    返回服务状态和配置信息，用于监控服务是否正常运行。
    """
    return {
        "status": "healthy",
        "config": {
            "vector_store_type": VECTOR_STORE_TYPE,
            "vector_store_path": VECTOR_STORE_PATH,
            "git_temp_dir": GIT_TEMP_DIR
        }
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)