from __future__ import annotations

from dataclasses import dataclass, asdict, is_dataclass
from typing import Any, Dict, List, Optional, Literal
import json
import urllib.error
import urllib.request


RPCErrorType = Literal[
    "custom",
    "validation",
    "input",
    "unauthorized",
    "forbidden",
    "not_implemented",
]


@dataclass
class RPCError:
    type: RPCErrorType
    message: str


class RPCErrorException(Exception):
    def __init__(self, error: RPCError) -> None:
        super().__init__(error.message)
        self.error = error


class CustomRPCError(RPCErrorException):
    pass


class ValidationRPCError(RPCErrorException):
    pass


class InputRPCError(RPCErrorException):
    pass


class UnauthorizedRPCError(RPCErrorException):
    pass


class ForbiddenRPCError(RPCErrorException):
    pass


class NotImplementedRPCError(RPCErrorException):
    pass


_ERROR_EXCEPTIONS = {
    "custom": CustomRPCError,
    "validation": ValidationRPCError,
    "input": InputRPCError,
    "unauthorized": UnauthorizedRPCError,
    "forbidden": ForbiddenRPCError,
    "not_implemented": NotImplementedRPCError,
}


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
            data = json.dumps(self._encode_payload(payload)).encode("utf-8")
        headers = {**self.headers, "Content-Type": "application/json"}
        req = urllib.request.Request(url, data=data, method="POST", headers=headers)
        try:
            with urllib.request.urlopen(req) as resp:
                body = resp.read()
        except urllib.error.HTTPError as err:
            with err as resp:
                detail = resp.read()
            try:
                parsed = json.loads(detail.decode("utf-8")) if detail else None
            except json.JSONDecodeError:
                parsed = None
            if parsed is not None:
                self._raise_if_error(parsed)
            text = detail.decode("utf-8", errors="replace")
            raise RuntimeError(f"rpc error: {text}") from err
        if not body:
            return None
        return json.loads(body.decode("utf-8"))

    def _raise_if_error(self, payload: Any) -> None:
        if not isinstance(payload, dict):
            return
        err_type = payload.get("type")
        message = payload.get("message")
        if not isinstance(err_type, str) or not isinstance(message, str):
            return
        exc_type = _ERROR_EXCEPTIONS.get(err_type)
        if exc_type is None:
            return
        raise exc_type(RPCError(type=err_type, message=message))

    def _encode_payload(self, value: Any) -> Any:
        if is_dataclass(value):
            return asdict(value)
        if isinstance(value, dict):
            return {k: self._encode_payload(v) for k, v in value.items()}
        if isinstance(value, list):
            return [self._encode_payload(item) for item in value]
        if isinstance(value, tuple):
            return tuple(self._encode_payload(item) for item in value)
        return value

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
