from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional
import json
import urllib.error
import urllib.request


@dataclass
class GreetingMessageModel:
    message: str

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "GreetingMessageModel":
        return GreetingMessageModel(
            message=data.get("message"),
        )


class RPCClient:
    def __init__(self, base_url: str) -> None:
        self.base_url = base_url.rstrip("/")

    def _request(self, path: str, payload: Optional[Dict[str, Any]]) -> Any:
        url = f"{self.base_url}/{path}"
        data = None
        if payload is not None:
            data = json.dumps(payload).encode("utf-8")
        req = urllib.request.Request(url, data=data, method="POST", headers={"Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(req) as resp:
                body = resp.read()
        except urllib.error.HTTPError as err:
            detail = err.read().decode("utf-8")
            raise RuntimeError(f"rpc error: {detail}") from err
        if not body:
            return None
        return json.loads(body.decode("utf-8"))

    def hello_world(self, name: str, surname: Optional[str] = None) -> GreetingMessageModel:
        payload = {
            "name": name,
            "surname": surname,
        }
        data = self._request("hello_world", payload)
        value = data.get("greeting_message") if isinstance(data, dict) else data
        return GreetingMessageModel.from_dict(value)
