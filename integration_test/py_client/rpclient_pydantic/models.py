# THIS CODE IS GENERATED

from __future__ import annotations


from pydantic import BaseModel

from typing import Any, Dict, List, Optional


class EmptyModel(BaseModel):
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "EmptyModel":
        _ = data
        return cls()


class TextModel(BaseModel):
    title: Optional[str]
    body: str

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "TextModel":
        try:
            return cls.model_validate(data)
        except AttributeError:
            return cls.parse_obj(data)


class FlagsModel(BaseModel):
    enabled: bool
    retries: int
    labels: List[str]
    meta: Dict[str, str]

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "FlagsModel":
        try:
            return cls.model_validate(data)
        except AttributeError:
            return cls.parse_obj(data)


class NestedModel(BaseModel):
    text: TextModel
    flags: Optional[FlagsModel]
    items: List[TextModel]
    lookup: Dict[str, TextModel]

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "NestedModel":
        try:
            return cls.model_validate(data)
        except AttributeError:
            return cls.parse_obj(data)


class PayloadModel(BaseModel):
    data: Any
    raw_data: Any

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "PayloadModel":
        try:
            return cls.model_validate(data)
        except AttributeError:
            return cls.parse_obj(data)
