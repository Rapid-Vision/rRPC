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
    EmptyModel,
    TextModel,
    FlagsModel,
    NestedModel,
    PayloadModel,
    TestBasicParams,
    TestListMapParams,
    TestOptionalParams,
    TestValidationErrorParams,
    TestJsonParams,
    TestRawParams,
    TestMixedPayloadParams,
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
    @app.post(f"{prefix}/test_empty")
    async def test_empty():
        try:
            result = handlers.test_empty()
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
            content={"empty": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_no_return")
    async def test_no_return():
        try:
            result = handlers.test_no_return()
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
        return Response(status_code=204)
    @app.post(f"{prefix}/test_basic")
    async def test_basic(params: TestBasicParams):
        try:
            result = handlers.test_basic(text=params.text, flag=params.flag, count=params.count, note=params.note, )
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
            content={"text": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_list_map")
    async def test_list_map(params: TestListMapParams):
        try:
            result = handlers.test_list_map(texts=params.texts, flags=params.flags, )
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
            content={"nested": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_optional")
    async def test_optional(params: TestOptionalParams):
        try:
            result = handlers.test_optional(text=params.text, flag=params.flag, )
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
            content={"flags": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_validation_error")
    async def test_validation_error(params: TestValidationErrorParams):
        try:
            result = handlers.test_validation_error(text=params.text, )
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
            content={"text": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_unauthorized_error")
    async def test_unauthorized_error():
        try:
            result = handlers.test_unauthorized_error()
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
            content={"empty": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_forbidden_error")
    async def test_forbidden_error():
        try:
            result = handlers.test_forbidden_error()
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
            content={"empty": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_not_implemented_error")
    async def test_not_implemented_error():
        try:
            result = handlers.test_not_implemented_error()
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
            content={"empty": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_custom_error")
    async def test_custom_error():
        try:
            result = handlers.test_custom_error()
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
            content={"empty": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_map_return")
    async def test_map_return():
        try:
            result = handlers.test_map_return()
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
            content={"result": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_json")
    async def test_json(params: TestJsonParams):
        try:
            result = handlers.test_json(data=params.data, )
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
            content={"json": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_raw")
    async def test_raw(params: TestRawParams):
        try:
            result = handlers.test_raw(payload=params.payload, )
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
            content={"raw": _encode_payload(result)}
        )
    @app.post(f"{prefix}/test_mixed_payload")
    async def test_mixed_payload(params: TestMixedPayloadParams):
        try:
            result = handlers.test_mixed_payload(payload=params.payload, )
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
            content={"payload": _encode_payload(result)}
        )
    return app
