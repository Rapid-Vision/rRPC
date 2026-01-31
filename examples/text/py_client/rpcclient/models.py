# THIS CODE IS GENERATED

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional


@dataclass
class TextModel:
    title: Optional[str]
    data: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "TextModel":
        return TextModel(
            title=None if data.get("title") is None else data.get("title"),
            data=data.get("data"),
        )


@dataclass
class SliceModel:
    begin: int
    end: int

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "SliceModel":
        return SliceModel(
            begin=data.get("begin"),
            end=data.get("end"),
        )


@dataclass
class StatsModel:
    ascii: bool
    word_count: Dict[str, int]
    total_words: int
    sentences: List[SliceModel]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "StatsModel":
        return StatsModel(
            ascii=data.get("ascii"),
            word_count={k: v for k, v in data.get("word_count").items()},
            total_words=data.get("total_words"),
            sentences=[SliceModel.from_dict(item) for item in data.get("sentences")],
        )
