# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Awaitable, Dict, List, Optional, Protocol, Union
from .models import (
    TextModel,
    SliceModel,
    StatsModel,
)


class RPCHandlers(Protocol):

    def submit_text(self, text: TextModel) -> Union[int, Awaitable[int]]:
        ...

    def compute_stats(self, text_id: int) -> Union[StatsModel, Awaitable[StatsModel]]:
        ...
