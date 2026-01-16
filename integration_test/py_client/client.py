#!/usr/bin/env python3

from rpc_client import (
    RPCClient,
    TextModel,
    ValidationRPCError,
    UnauthorizedRPCError,
    ForbiddenRPCError,
    NotImplementedRPCError,
)


def main() -> None:
    rpc = RPCClient("http://localhost:8080")

    empty = rpc.test_empty()
    print("empty:", empty)

    text = TextModel(title=None, body="  hello  ")
    basic = rpc.test_basic(text=text, flag=True, count=3, note="note")
    print("basic:", basic)

    nested = rpc.test_list_map(
        texts=[TextModel(title="t1", body="b1"), TextModel(title="t2", body="b2")],
        flags={"mode": "fast"},
    )
    print("nested:", nested)

    optional = rpc.test_optional(text=None, flag=None)
    print("optional:", optional)

    try:
        rpc.test_validation_error(text=TextModel(title=None, body=""))
    except ValidationRPCError as exc:
        print("validation:", exc)

    try:
        rpc.test_unauthorized_error()
    except UnauthorizedRPCError as exc:
        print("unauthorized:", exc)

    try:
        rpc.test_forbidden_error()
    except ForbiddenRPCError as exc:
        print("forbidden:", exc)

    try:
        rpc.test_not_implemented_error()
    except NotImplementedRPCError as exc:
        print("not_implemented:", exc)


if __name__ == "__main__":
    main()
