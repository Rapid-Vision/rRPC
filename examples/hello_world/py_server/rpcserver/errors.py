# THIS CODE IS GENERATED

from __future__ import annotations

from typing import Dict

from pydantic import BaseModel

ERROR_TYPE_VALIDATION = "validation"
ERROR_TYPE_INPUT = "input"
ERROR_TYPE_UNAUTHORIZED = "unauthorized"
ERROR_TYPE_FORBIDDEN = "forbidden"
ERROR_TYPE_NOT_IMPLEMENTED = "not_implemented"
ERROR_TYPE_CUSTOM = "custom"

ERROR_STATUS: Dict[str, int] = {
    ERROR_TYPE_VALIDATION: 400,
    ERROR_TYPE_INPUT: 400,
    ERROR_TYPE_UNAUTHORIZED: 401,
    ERROR_TYPE_FORBIDDEN: 403,
    ERROR_TYPE_NOT_IMPLEMENTED: 501,
    ERROR_TYPE_CUSTOM: 500,
}


class RPCError(BaseModel):
    type: str
    message: str


class RPCErrorException(Exception):
    def __init__(self, error: RPCError, status_code: int) -> None:
        super().__init__(error.message)
        self.error = error
        self.status_code = status_code


class ValidationRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(RPCError(type=ERROR_TYPE_VALIDATION, message=message), 400)


class InputRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(RPCError(type=ERROR_TYPE_INPUT, message=message), 400)


class UnauthorizedRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(RPCError(type=ERROR_TYPE_UNAUTHORIZED, message=message), 401)


class ForbiddenRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(RPCError(type=ERROR_TYPE_FORBIDDEN, message=message), 403)


class NotImplementedRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(
            RPCError(type=ERROR_TYPE_NOT_IMPLEMENTED, message=message), 501
        )


class CustomRPCError(RPCErrorException):
    def __init__(self, message: str) -> None:
        super().__init__(RPCError(type=ERROR_TYPE_CUSTOM, message=message), 500)


def error_payload(error_type: str, message: str) -> Dict[str, str]:
    return {"type": error_type, "message": message}


def error_dict(error: RPCError) -> Dict[str, str]:
    try:
        return error.model_dump()
    except AttributeError:
        return error.dict()
