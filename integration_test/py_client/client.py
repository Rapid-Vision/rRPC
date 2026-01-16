#!/usr/bin/env python3

import sys

from rpc_client import (
    RPCClient,
    EmptyModel,
    TextModel,
    CustomRPCError,
    InputRPCError,
    ValidationRPCError,
    UnauthorizedRPCError,
    ForbiddenRPCError,
    NotImplementedRPCError,
)


def run_test(name: str, fn) -> bool:
    try:
        fn()
    except AssertionError as exc:
        print(f"FAIL {name}: {exc}")
        return False
    except Exception as exc:
        print(f"FAIL {name}: unexpected error: {exc}")
        return False
    print(f"OK {name}")
    return True


def main() -> None:
    rpc = RPCClient("http://localhost:8080", headers={"Authorization": "Bearer test_token"})
    passed = 0
    total = 0

    def test_empty() -> None:
        empty = rpc.test_empty()
        assert isinstance(empty, EmptyModel)

    def test_basic() -> None:
        text = TextModel(title=None, body="  hello  ")
        basic = rpc.test_basic(text=text, flag=True, count=3, note="note")
        assert basic.body == "hello"
        assert basic.title == "note"

    def test_list_map() -> None:
        nested = rpc.test_list_map(
            texts=[TextModel(title="t1", body="b1"), TextModel(title="t2", body="b2")],
            flags={"mode": "fast"},
        )
        assert nested.flags is not None
        assert nested.flags.retries == 2
        assert nested.flags.meta.get("mode") == "fast"
        assert isinstance(nested.lookup.get("first"), TextModel)

    def test_optional() -> None:
        optional = rpc.test_optional(text=None, flag=None)
        assert optional.enabled is False

    def test_validation_error() -> None:
        try:
            rpc.test_validation_error(text=TextModel(title=None, body=""))
        except ValidationRPCError:
            return
        raise AssertionError("expected ValidationRPCError")

    def test_input_error() -> None:
        try:
            rpc.test_basic(text="bad", flag=True, count=1, note=None)
        except InputRPCError:
            return
        raise AssertionError("expected InputRPCError")

    def test_unauthorized_error() -> None:
        try:
            rpc.test_unauthorized_error()
        except UnauthorizedRPCError:
            return
        raise AssertionError("expected UnauthorizedRPCError")

    def test_forbidden_error() -> None:
        try:
            rpc.test_forbidden_error()
        except ForbiddenRPCError:
            return
        raise AssertionError("expected ForbiddenRPCError")

    def test_not_implemented_error() -> None:
        try:
            rpc.test_not_implemented_error()
        except NotImplementedRPCError:
            return
        raise AssertionError("expected NotImplementedRPCError")

    def test_custom_error() -> None:
        try:
            rpc.test_custom_error()
        except CustomRPCError:
            return
        raise AssertionError("expected CustomRPCError")

    def test_map_return() -> None:
        mapped = rpc.test_map_return()
        assert isinstance(mapped, dict)
        assert isinstance(mapped.get("a"), TextModel)
        assert mapped["a"].body == "mapped"

    tests = [
        ("empty", test_empty),
        ("basic", test_basic),
        ("list_map", test_list_map),
        ("optional", test_optional),
        ("validation_error", test_validation_error),
        ("input_error", test_input_error),
        ("unauthorized_error", test_unauthorized_error),
        ("forbidden_error", test_forbidden_error),
        ("not_implemented_error", test_not_implemented_error),
        ("custom_error", test_custom_error),
        ("map_return", test_map_return),
    ]

    for name, fn in tests:
        total += 1
        if run_test(name, fn):
            passed += 1

    print(f"passed {passed}/{total}")
    if passed != total:
        sys.exit(1)


if __name__ == "__main__":
    main()
