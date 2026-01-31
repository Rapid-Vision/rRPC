# THIS CODE IS GENERATED

from __future__ import annotations


from dataclasses import dataclass

from typing import Any, Dict, List, Optional
@dataclass
class GreetingMessageModel:
    message: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "GreetingMessageModel":
        return GreetingMessageModel(
            message=data.get("message"),
        )

