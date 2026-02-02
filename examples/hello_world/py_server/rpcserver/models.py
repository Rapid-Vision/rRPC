# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Dict, List, Optional

from pydantic import BaseModel


class GreetingMessageModel(BaseModel):
    message: str


class HelloWorldParams(BaseModel):
    name: str
    surname: Optional[str] = None
