# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Awaitable, Dict, List, Optional, Protocol, Union
from .models import (
    GreetingMessageModel,
)


class RPCHandlers(Protocol):

    def hello_world(self, name: str, surname: Optional[str] = None) -> Union[GreetingMessageModel, Awaitable[GreetingMessageModel]]:
        ...
