# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Dict, List, Optional

from pydantic import BaseModel


class EmptyModel(BaseModel):
    pass


class TextModel(BaseModel):
    title: Optional[str] = None
    body: str


class FlagsModel(BaseModel):
    enabled: bool
    retries: int
    labels: List[str]
    meta: Dict[str, str]


class NestedModel(BaseModel):
    text: TextModel
    flags: Optional[FlagsModel] = None
    items: List[TextModel]
    lookup: Dict[str, TextModel]


class PayloadModel(BaseModel):
    data: Any
    raw_data: Any


class TestBasicParams(BaseModel):
    text: TextModel
    flag: bool
    count: int
    note: Optional[str] = None


class TestListMapParams(BaseModel):
    texts: List[TextModel]
    flags: Dict[str, str]


class TestOptionalParams(BaseModel):
    text: Optional[TextModel] = None
    flag: Optional[bool] = None


class TestValidationErrorParams(BaseModel):
    text: TextModel


class TestJsonParams(BaseModel):
    data: Any


class TestRawParams(BaseModel):
    payload: Any


class TestMixedPayloadParams(BaseModel):
    payload: PayloadModel
