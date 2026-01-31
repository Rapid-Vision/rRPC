# THIS CODE IS GENERATED

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional


@dataclass
class EmptyModel:
    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "EmptyModel":
        _ = data
        return EmptyModel()


@dataclass
class TextModel:
    title: Optional[str]
    body: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "TextModel":
        return TextModel(
            title=None if data.get("title") is None else data.get("title"),
            body=data.get("body"),
        )


@dataclass
class FlagsModel:
    enabled: bool
    retries: int
    labels: List[str]
    meta: Dict[str, str]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "FlagsModel":
        return FlagsModel(
            enabled=data.get("enabled"),
            retries=data.get("retries"),
            labels=[item for item in data.get("labels")],
            meta={k: v for k, v in data.get("meta").items()},
        )


@dataclass
class NestedModel:
    text: TextModel
    flags: Optional[FlagsModel]
    items: List[TextModel]
    lookup: Dict[str, TextModel]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "NestedModel":
        return NestedModel(
            text=TextModel.from_dict(data.get("text")),
            flags=None if data.get("flags") is None else FlagsModel.from_dict(data.get("flags")),
            items=[TextModel.from_dict(item) for item in data.get("items")],
            lookup={k: TextModel.from_dict(v) for k, v in data.get("lookup").items()},
        )


@dataclass
class PayloadModel:
    data: Any
    raw_data: Any

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "PayloadModel":
        return PayloadModel(
            data=data.get("data"),
            raw_data=data.get("raw_data"),
        )
