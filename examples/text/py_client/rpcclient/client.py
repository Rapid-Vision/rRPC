# THIS CODE IS GENERATED

from __future__ import annotations

from dataclasses import asdict, is_dataclass
from typing import Any, Dict, List, Optional
import json
import urllib.error
import urllib.request

from .errors import RPCError, RPCErrorException, _ERROR_EXCEPTIONS
from .models import (
    TextModel,
    SliceModel,
    StatsModel,
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
