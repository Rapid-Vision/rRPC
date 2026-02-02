# THIS CODE IS GENERATED

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional
from .models import (
    EmptyModel,
    TextModel,
    FlagsModel,
    NestedModel,
    PayloadModel,
)


class RPCHandlers(ABC):

    @abstractmethod
    def test_empty(self) -> EmptyModel:
        raise NotImplementedError

    @abstractmethod
    def test_no_return(self) -> None:
        raise NotImplementedError

    @abstractmethod
    def test_basic(self, text: TextModel, flag: bool, count: int, note: Optional[str] = None) -> TextModel:
        raise NotImplementedError

    @abstractmethod
    def test_list_map(self, texts: List[TextModel], flags: Dict[str, str]) -> NestedModel:
        raise NotImplementedError

    @abstractmethod
    def test_optional(self, text: Optional[TextModel] = None, flag: Optional[bool] = None) -> FlagsModel:
        raise NotImplementedError

    @abstractmethod
    def test_validation_error(self, text: TextModel) -> TextModel:
        raise NotImplementedError

    @abstractmethod
    def test_unauthorized_error(self) -> EmptyModel:
        raise NotImplementedError

    @abstractmethod
    def test_forbidden_error(self) -> EmptyModel:
        raise NotImplementedError

    @abstractmethod
    def test_not_implemented_error(self) -> EmptyModel:
        raise NotImplementedError

    @abstractmethod
    def test_custom_error(self) -> EmptyModel:
        raise NotImplementedError

    @abstractmethod
    def test_map_return(self) -> Dict[str, TextModel]:
        raise NotImplementedError

    @abstractmethod
    def test_json(self, data: Any) -> Any:
        raise NotImplementedError

    @abstractmethod
    def test_raw(self, payload: Any) -> Any:
        raise NotImplementedError

    @abstractmethod
    def test_mixed_payload(self, payload: PayloadModel) -> PayloadModel:
        raise NotImplementedError
