from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional
import json
import urllib.error
import urllib.request


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


class RPCClient:
    def __init__(self, base_url: str, headers: Optional[Dict[str, str]] = None) -> None:
        self.base_url = base_url.rstrip("/")
        self.headers = headers or {}

    def _request(self, path: str, payload: Optional[Dict[str, Any]]) -> Any:
        url = f"{self.base_url}/rpc/{path}"
        data = None
        if payload is not None:
            data = json.dumps(payload).encode("utf-8")
        headers = {**self.headers, "Content-Type": "application/json"}
        req = urllib.request.Request(url, data=data, method="POST", headers=headers)
        try:
            with urllib.request.urlopen(req) as resp:
                body = resp.read()
        except urllib.error.HTTPError as err:
            detail = err.read().decode("utf-8")
            raise RuntimeError(f"rpc error: {detail}") from err
        if not body:
            return None
        return json.loads(body.decode("utf-8"))

    def submit_text(self, text: TextModel) -> int:
        payload = {
            "text": text,
        }
        data = self._request("submit_text", payload)
        value = data.get("int") if isinstance(data, dict) else data
        return value

    def compute_stats(self, text_id: int) -> StatsModel:
        payload = {
            "text_id": text_id,
        }
        data = self._request("compute_stats", payload)
        value = data.get("stats") if isinstance(data, dict) else data
        return StatsModel.from_dict(value)
