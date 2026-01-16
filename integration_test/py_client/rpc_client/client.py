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
class EmptyModelModel:
    pass


@dataclass
class TextModelModel:
    title: Optional[str]
    body: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "TextModelModel":
        return TextModelModel(
            title=None if data.get("title") is None else data.get("title"),
            body=data.get("body"),
        )


@dataclass
class FlagsModelModel:
    enabled: bool
    retries: int
    labels: List[str]
    meta: Dict[str, str]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "FlagsModelModel":
        return FlagsModelModel(
            enabled=data.get("enabled"),
            retries=data.get("retries"),
            labels=[item for item in data.get("labels")],
            meta={k: v for k, v in data.get("meta").items()},
        )


@dataclass
class NestedModelModel:
    text: TextModelModel
    flags: Optional[FlagsModelModel]
    items: List[TextModelModel]
    lookup: Dict[str, TextModelModel]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "NestedModelModel":
        return NestedModelModel(
            text=TextModelModel.from_dict(data.get("text")),
            flags=None if data.get("flags") is None else FlagsModelModel.from_dict(data.get("flags")),
            items=[TextModelModel.from_dict(item) for item in data.get("items")],
            lookup={k: TextModelModel.from_dict(v) for k, v in data.get("lookup").items()},
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
            detail = err.read()
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

    def test_empty(self) -> EmptyModelModel:
        payload = None
        data = self._request("test_empty", payload)
        value = data.get("empty_model") if isinstance(data, dict) else data
        return EmptyModelModel.from_dict(value)

    def test_basic(self, text: TextModelModel, flag: bool, count: int, note: Optional[str] = None) -> TextModelModel:
        payload = {
            "text": text,
            "flag": flag,
            "count": count,
            "note": note,
        }
        data = self._request("test_basic", payload)
        value = data.get("text_model") if isinstance(data, dict) else data
        return TextModelModel.from_dict(value)

    def test_list_map(self, texts: List[TextModelModel], flags: Dict[str, str]) -> NestedModelModel:
        payload = {
            "texts": texts,
            "flags": flags,
        }
        data = self._request("test_list_map", payload)
        value = data.get("nested_model") if isinstance(data, dict) else data
        return NestedModelModel.from_dict(value)

    def test_optional(self, text: Optional[TextModelModel] = None, flag: Optional[bool] = None) -> FlagsModelModel:
        payload = {
            "text": text,
            "flag": flag,
        }
        data = self._request("test_optional", payload)
        value = data.get("flags_model") if isinstance(data, dict) else data
        return FlagsModelModel.from_dict(value)

    def test_validation_error(self, text: TextModelModel) -> TextModelModel:
        payload = {
            "text": text,
        }
        data = self._request("test_validation_error", payload)
        value = data.get("text_model") if isinstance(data, dict) else data
        return TextModelModel.from_dict(value)

    def test_unauthorized_error(self) -> EmptyModelModel:
        payload = None
        data = self._request("test_unauthorized_error", payload)
        value = data.get("empty_model") if isinstance(data, dict) else data
        return EmptyModelModel.from_dict(value)

    def test_forbidden_error(self) -> EmptyModelModel:
        payload = None
        data = self._request("test_forbidden_error", payload)
        value = data.get("empty_model") if isinstance(data, dict) else data
        return EmptyModelModel.from_dict(value)

    def test_not_implemented_error(self) -> EmptyModelModel:
        payload = None
        data = self._request("test_not_implemented_error", payload)
        value = data.get("empty_model") if isinstance(data, dict) else data
        return EmptyModelModel.from_dict(value)
