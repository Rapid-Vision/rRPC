# THIS CODE IS GENERATED

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional
from .models import (
    GreetingMessageModel,
)


class RPCHandlers(ABC):

    @abstractmethod
    def hello_world(self, name: str, surname: Optional[str] = None) -> GreetingMessageModel:
        raise NotImplementedError
