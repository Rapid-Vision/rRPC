from __future__ import annotations

import json
from typing import Any, Dict, List, Optional

from fastapi import Request
from fastapi.responses import JSONResponse, Response

from rpcserver import (
    CustomRPCError,
    ForbiddenRPCError,
    InputRPCError,
    NotImplementedRPCError,
    RPCHandlers,
    UnauthorizedRPCError,
    ValidationRPCError,
    create_app,
)
from rpcserver.models import (
    EmptyModel,
    FlagsModel,
    NestedModel,
    PayloadModel,
    TextModel,
)

BEARER_TOKEN = "test_token"


class Service(RPCHandlers):
    def test_empty(self) -> EmptyModel:
        return EmptyModel()

    def test_no_return(self) -> None:
        return None

    def test_basic(
        self,
        text: TextModel,
        flag: bool,
        count: int,
        note: Optional[str] = None,
    ) -> TextModel:
        _ = flag
        _ = count
        title = text.title
        if title is None and note is not None:
            title = note
        return TextModel(title=title, body=text.body.strip())

    def test_list_map(self, texts: List[TextModel], flags: Dict[str, str]) -> NestedModel:
        _ = flags
        meta = dict(flags)
        out_flags = FlagsModel(
            enabled=True,
            retries=len(texts),
            labels=["ok"],
            meta=meta,
        )
        return NestedModel(
            text=texts[0],
            flags=out_flags,
            items=texts,
            lookup={"first": texts[0]},
        )

    def test_optional(
        self, text: Optional[TextModel] = None, flag: Optional[bool] = None
    ) -> FlagsModel:
        enabled = bool(flag)
        return FlagsModel(
            enabled=enabled,
            retries=0,
            labels=[],
            meta={},
        )

    def test_validation_error(self, text: TextModel) -> TextModel:
        if text.body.strip() == "":
            raise ValidationRPCError("body is required")
        return text

    def test_unauthorized_error(self) -> EmptyModel:
        raise UnauthorizedRPCError("missing token")

    def test_forbidden_error(self) -> EmptyModel:
        raise ForbiddenRPCError("not allowed")

    def test_not_implemented_error(self) -> EmptyModel:
        raise NotImplementedRPCError("not implemented")

    def test_custom_error(self) -> EmptyModel:
        raise CustomRPCError("custom failure")

    def test_map_return(self) -> Dict[str, TextModel]:
        return {"a": TextModel(title=None, body="mapped")}

    def test_json(self, data: Any) -> Any:
        return data

    def test_raw(self, payload: Any) -> Any:
        if not json.dumps(payload):
            raise ValidationRPCError("payload is not valid json")
        return payload

    def test_mixed_payload(self, payload: PayloadModel) -> PayloadModel:
        return payload


app = create_app(Service())


@app.middleware("http")
async def auth_middleware(request: Request, call_next):
    if request.headers.get("Authorization") != f"Bearer {BEARER_TOKEN}":
        return JSONResponse(
            status_code=401, content={"type": "auth", "message": "missing or invalid token"}
        )
    response = await call_next(request)
    if response.status_code == 204:
        return Response(status_code=204)
    return response


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("server:app", host="127.0.0.1", port=8080, reload=False)
