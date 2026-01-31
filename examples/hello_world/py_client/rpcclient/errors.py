# THIS CODE IS GENERATED

from dataclasses import dataclass
from typing import Literal

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
