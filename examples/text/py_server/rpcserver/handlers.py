# THIS CODE IS GENERATED

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional
from .models import (
    TextModel,
    SliceModel,
    StatsModel,
)


class RPCHandlers(ABC):

    @abstractmethod
    def submit_text(self, text: TextModel) -> int:
        raise NotImplementedError

    @abstractmethod
    def compute_stats(self, text_id: int) -> StatsModel:
        raise NotImplementedError
