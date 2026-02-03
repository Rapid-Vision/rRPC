# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Any, Awaitable, Dict, List, Optional, Protocol, Union
from .models import (
    EmptyModel,
    TextModel,
    FlagsModel,
    NestedModel,
    PayloadModel,
)


class RPCHandlers(Protocol):

    def test_empty(self) -> Union[EmptyModel, Awaitable[EmptyModel]]:
        ...

    def test_no_return(self) -> Union[None, Awaitable[None]]:
        ...

    def test_basic(self, text: TextModel, flag: bool, count: int, note: Optional[str] = None) -> Union[TextModel, Awaitable[TextModel]]:
        ...

    def test_list_map(self, texts: List[TextModel], flags: Dict[str, str]) -> Union[NestedModel, Awaitable[NestedModel]]:
        ...

    def test_optional(self, text: Optional[TextModel] = None, flag: Optional[bool] = None) -> Union[FlagsModel, Awaitable[FlagsModel]]:
        ...

    def test_validation_error(self, text: TextModel) -> Union[TextModel, Awaitable[TextModel]]:
        ...

    def test_unauthorized_error(self) -> Union[EmptyModel, Awaitable[EmptyModel]]:
        ...

    def test_forbidden_error(self) -> Union[EmptyModel, Awaitable[EmptyModel]]:
        ...

    def test_not_implemented_error(self) -> Union[EmptyModel, Awaitable[EmptyModel]]:
        ...

    def test_custom_error(self) -> Union[EmptyModel, Awaitable[EmptyModel]]:
        ...

    def test_map_return(self) -> Union[Dict[str, TextModel], Awaitable[Dict[str, TextModel]]]:
        ...

    def test_json(self, data: Any) -> Union[Any, Awaitable[Any]]:
        ...

    def test_raw(self, payload: Any) -> Union[Any, Awaitable[Any]]:
        ...

    def test_mixed_payload(self, payload: PayloadModel) -> Union[PayloadModel, Awaitable[PayloadModel]]:
        ...
