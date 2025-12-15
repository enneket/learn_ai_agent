"""配置管理"""
import os
from typing import Optional, Dict, Any
from pydantic import BaseModel

class Config(BaseModel):
    """HelloAgents配置类"""
    
    # LLM配置
    default_model: str = "gpt-3.5-turbo"
    default_provider: str = "openai"
    temperature: float = 0.7
    max_tokens: Optional[int] = None
    
    # 系统配置
    debug: bool = False
    log_level: str = "INFO"
    
    # 其他配置
    max_history_length: int = 100
    
    @classmethod
    def from_env(cls) -> "Config":
        """从环境变量创建配置"""
        return cls(
            debug=os.getenv("DEBUG", "false").lower() == "true",
            log_level=os.getenv("LOG_LEVEL", "INFO"),
            temperature=float(os.getenv("TEMPERATURE", "0.7")),
            max_tokens=int(os.getenv("MAX_TOKENS")) if os.getenv("MAX_TOKENS") else None,
        )
    
    def to_dict(self) -> Dict[str, Any]:
        """转换为字典格式"""
        return {
            "default_model": self.default_model,
            "default_provider": self.default_provider,
            "temperature": self.temperature,
            "max_tokens": self.max_tokens,
            "debug": self.debug,
            "log_level": self.log_level,
            "max_history_length": self.max_history_length
        }
