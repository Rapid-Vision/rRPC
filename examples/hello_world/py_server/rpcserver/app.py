# THIS CODE IS GENERATED

from __future__ import annotations

import inspect
from typing import Any

from fastapi import FastAPI
from fastapi.responses import JSONResponse, Response
from pydantic import BaseModel, ValidationError

from .errors import (
    ERROR_TYPE_CUSTOM,
    ERROR_TYPE_VALIDATION,
    RPCErrorException,
    error_payload,
    error_dict,
)
from .handlers import RPCHandlers
from .models import (
    GreetingMessageModel,
    HelloWorldParams,
)


def _normalize_prefix(prefix: str) -> str:
    prefix = prefix.strip()
    if prefix == "":
        return ""
    if not prefix.startswith("/"):
        prefix = "/" + prefix
    return prefix.rstrip("/")


def _encode_payload(value: Any) -> Any:
    if isinstance(value, BaseModel):
        try:
            return value.model_dump()
        except AttributeError:
            return value.dict()
    if isinstance(value, dict):
        return {k: _encode_payload(v) for k, v in value.items()}
    if isinstance(value, list):
        return [_encode_payload(item) for item in value]
    if isinstance(value, tuple):
        return tuple(_encode_payload(item) for item in value)
    return value


def create_app(handlers: RPCHandlers, prefix: str = "/rpc") -> FastAPI:
    app = FastAPI()
    prefix = _normalize_prefix(prefix)
    @app.post(f"{prefix}/hello_world")
    async def hello_world(params: HelloWorldParams):
        try:
            result = handlers.hello_world(name=params.name, surname=params.surname, )
            if inspect.isawaitable(result):
                result = await result
        except ValidationError as err:
            return JSONResponse(
                status_code=400,
                content=error_payload(ERROR_TYPE_VALIDATION, str(err)),
            )
        except RPCErrorException as err:
            return JSONResponse(
                status_code=err.status_code,
                content=error_dict(err.error),
            )
        except Exception as err:
            return JSONResponse(
                status_code=500,
                content=error_payload(ERROR_TYPE_CUSTOM, str(err)),
            )
        return JSONResponse(
            content={"greeting_message": _encode_payload(result)}
        )
    return app
