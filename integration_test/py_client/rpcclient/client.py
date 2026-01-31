# THIS CODE IS GENERATED

from __future__ import annotations

from dataclasses import asdict, is_dataclass
from typing import Any, Dict, List, Optional
import json
import urllib.error
import urllib.request

from .errors import RPCError, RPCErrorException, _ERROR_EXCEPTIONS
from .models import (
    EmptyModel,
    TextModel,
    FlagsModel,
    NestedModel,
    PayloadModel,
)


class RPCClient:
    def __init__(
        self,
        base_url: str,
        prefix: str = "/rpc",
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None,
    ) -> None:
        self.base_url = self._normalize_base_url(base_url)
        self.prefix = self._normalize_prefix(prefix)
        self.headers = headers or {}
        self.timeout = timeout

    @staticmethod
    def _normalize_base_url(base_url: str) -> str:
        base_url = base_url.strip()
        if "://" not in base_url:
            base_url = "http://" + base_url
        return base_url.rstrip("/")

    @staticmethod
    def _normalize_prefix(prefix: str) -> str:
        prefix = prefix.strip()
        if prefix == "":
            return ""
        if not prefix.startswith("/"):
            prefix = "/" + prefix
        return prefix.rstrip("/")

    def _request(self, path: str, payload: Optional[Dict[str, Any]]) -> Any:
        if self.prefix:
            url = f"{self.base_url}{self.prefix}/{path}"
        else:
            url = f"{self.base_url}/{path}"
        data = None
        if payload is not None:
            data = json.dumps(self._encode_payload(payload)).encode("utf-8")
        headers = {**self.headers, "Content-Type": "application/json"}
        req = urllib.request.Request(url, data=data, method="POST", headers=headers)
        try:
            if self.timeout is None:
                with urllib.request.urlopen(req) as resp:
                    body = resp.read()
            else:
                with urllib.request.urlopen(req, timeout=self.timeout) as resp:
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
            raise RPCErrorException(
                RPCError(type="custom", message=f"rpc error: status {err.code}")
            ) from err
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

    def test_empty(self) -> EmptyModel:
        payload = None
        data = self._request("test_empty", payload)
        value = data.get("empty") if isinstance(data, dict) else data
        return EmptyModel.from_dict(value)

    def test_no_return(self) -> None:
        payload = None
        data = self._request("test_no_return", payload)
        return None

    def test_basic(self, text: TextModel, flag: bool, count: int, note: Optional[str] = None) -> TextModel:
        payload = {
            "text": text,
            "flag": flag,
            "count": count,
            "note": note,
        }
        data = self._request("test_basic", payload)
        value = data.get("text") if isinstance(data, dict) else data
        return TextModel.from_dict(value)

    def test_list_map(self, texts: List[TextModel], flags: Dict[str, str]) -> NestedModel:
        payload = {
            "texts": texts,
            "flags": flags,
        }
        data = self._request("test_list_map", payload)
        value = data.get("nested") if isinstance(data, dict) else data
        return NestedModel.from_dict(value)

    def test_optional(self, text: Optional[TextModel] = None, flag: Optional[bool] = None) -> FlagsModel:
        payload = {
            "text": text,
            "flag": flag,
        }
        data = self._request("test_optional", payload)
        value = data.get("flags") if isinstance(data, dict) else data
        return FlagsModel.from_dict(value)

    def test_validation_error(self, text: TextModel) -> TextModel:
        payload = {
            "text": text,
        }
        data = self._request("test_validation_error", payload)
        value = data.get("text") if isinstance(data, dict) else data
        return TextModel.from_dict(value)

    def test_unauthorized_error(self) -> EmptyModel:
        payload = None
        data = self._request("test_unauthorized_error", payload)
        value = data.get("empty") if isinstance(data, dict) else data
        return EmptyModel.from_dict(value)

    def test_forbidden_error(self) -> EmptyModel:
        payload = None
        data = self._request("test_forbidden_error", payload)
        value = data.get("empty") if isinstance(data, dict) else data
        return EmptyModel.from_dict(value)

    def test_not_implemented_error(self) -> EmptyModel:
        payload = None
        data = self._request("test_not_implemented_error", payload)
        value = data.get("empty") if isinstance(data, dict) else data
        return EmptyModel.from_dict(value)

    def test_custom_error(self) -> EmptyModel:
        payload = None
        data = self._request("test_custom_error", payload)
        value = data.get("empty") if isinstance(data, dict) else data
        return EmptyModel.from_dict(value)

    def test_map_return(self) -> Dict[str, TextModel]:
        payload = None
        data = self._request("test_map_return", payload)
        value = data.get("result") if isinstance(data, dict) else data
        return {k: TextModel.from_dict(v) for k, v in value.items()}

    def test_json(self, data: Any) -> Any:
        payload = {
            "data": data,
        }
        data = self._request("test_json", payload)
        value = data.get("json") if isinstance(data, dict) else data
        return value

    def test_raw(self, payload: Any) -> Any:
        payload = {
            "payload": payload,
        }
        data = self._request("test_raw", payload)
        value = data.get("raw") if isinstance(data, dict) else data
        return value

    def test_mixed_payload(self, payload: PayloadModel) -> PayloadModel:
        payload = {
            "payload": payload,
        }
        data = self._request("test_mixed_payload", payload)
        value = data.get("payload") if isinstance(data, dict) else data
        return PayloadModel.from_dict(value)
