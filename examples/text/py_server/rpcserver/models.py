# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Dict, List, Optional

from pydantic import BaseModel


class TextModel(BaseModel):
    title: Optional[str] = None
    data: str


class SliceModel(BaseModel):
    begin: int
    end: int


class StatsModel(BaseModel):
    ascii: bool
    word_count: Dict[str, int]
    total_words: int
    sentences: List[SliceModel]


class SubmitTextParams(BaseModel):
    text: TextModel


class ComputeStatsParams(BaseModel):
    text_id: int
