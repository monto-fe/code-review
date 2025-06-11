from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict
from rag_engine import analyze_gitlab_code
import os
from dotenv import load_dotenv
from pathlib import Path

# 加载上级目录的.env文件
env_path = Path(__file__).parent / '.env'
load_dotenv(dotenv_path=env_path)

# 打印env_path
print("--------------------------------")
print(env_path)
print(os.getenv("OPENAI_API_KEY"))
print("--------------------------------")

app = FastAPI(title="Code Review RAG Service")

# 从环境变量获取配置
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
GITLAB_TOKEN = os.getenv("GITLAB_TOKEN")
RAG_SERVICE_HOST = os.getenv("RAG_SERVICE_HOST", "0.0.0.0")
RAG_SERVICE_PORT = int(os.getenv("RAG_SERVICE_PORT", "8000"))
VECTOR_STORE_TYPE = os.getenv("VECTOR_STORE_TYPE", "faiss")
VECTOR_STORE_PATH = os.getenv("VECTOR_STORE_PATH", "/app/data/vector_store")
GIT_TEMP_DIR = os.getenv("GIT_TEMP_DIR", "/tmp/git_repos")

class CodeReviewRequest(BaseModel):
    git_url: str
    branch: str
    diff_content: str
    query: Optional[str] = None
    gitlab_token: Optional[str] = None

class CodeAnalysisResponse(BaseModel):
    code_changes: str
    context: str
    changed_files: List[str]
    query: Optional[str]

@app.post("/api/code-analysis", response_model=CodeAnalysisResponse)
async def code_analysis(request: CodeReviewRequest):
    try:
        result = analyze_gitlab_code(
            git_url=request.git_url,
            branch=request.branch,
            diff_content=request.diff_content,
            query=request.query,
            gitlab_token=request.gitlab_token
        )
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
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