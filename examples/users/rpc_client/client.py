from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Optional
import json
import urllib.error
import urllib.request


@dataclass
class UserModel:
    id: int
    username: str
    name: str
    surname: Optional[str]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "UserModel":
        return UserModel(
            id=data.get("id"),
            username=data.get("username"),
            name=data.get("name"),
            surname=None if data.get("surname") is None else data.get("surname"),
        )


@dataclass
class GroupModel:
    name: str
    users: List[UserModel]

    @staticmethod
    def from_dict(data: Dict[str, Any]) -> "GroupModel":
        return GroupModel(
            name=data.get("name"),
            users=[UserModel.from_dict(item) for item in data.get("users")],
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

    def get_user(self, user_id: int) -> UserModel:
        payload = {
            "user_id": user_id,
        }
        data = self._request("get_user", payload)
        value = data.get("user") if isinstance(data, dict) else data
        return UserModel.from_dict(value)

    def list_users(self) -> List[UserModel]:
        payload = None
        data = self._request("list_users", payload)
        value = data.get("result") if isinstance(data, dict) else data
        return [UserModel.from_dict(item) for item in value]

    def create_user(self, name: str, surname: Optional[str] = None) -> UserModel:
        payload = {
            "name": name,
            "surname": surname,
        }
        data = self._request("create_user", payload)
        value = data.get("user") if isinstance(data, dict) else data
        return UserModel.from_dict(value)

    def get_username_map(self) -> Dict[str, UserModel]:
        payload = None
        data = self._request("get_username_map", payload)
        value = data.get("result") if isinstance(data, dict) else data
        return {k: UserModel.from_dict(v) for k, v in value.items()}

    def find_group_by_name(self, name: str) -> GroupModel:
        payload = {
            "name": name,
        }
        data = self._request("find_group_by_name", payload)
        value = data.get("group") if isinstance(data, dict) else data
        return GroupModel.from_dict(value)
