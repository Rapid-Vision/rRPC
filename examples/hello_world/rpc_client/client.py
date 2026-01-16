from __future__ import annotations

from dataclasses import dataclass
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
class GreetingMessageModel:
    message: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "GreetingMessageModel":
        return GreetingMessageModel(
            message=data.get("message"),
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

    def hello_world(self, name: str, surname: Optional[str] = None) -> GreetingMessageModel:
        payload = {
            "name": name,
            "surname": surname,
        }
        data = self._request("hello_world", payload)
        value = data.get("greeting_message") if isinstance(data, dict) else data
        return GreetingMessageModel.from_dict(value)
