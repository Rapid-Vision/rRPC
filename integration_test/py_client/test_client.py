import unittest

from rpc_client import (
    RPCClient,
    EmptyModel,
    PayloadModel,
    TextModel,
    CustomRPCError,
    InputRPCError,
    ValidationRPCError,
    UnauthorizedRPCError,
    ForbiddenRPCError,
    NotImplementedRPCError,
)


class RPCClientTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.rpc = RPCClient("http://localhost:8080", headers={"Authorization": "Bearer test_token"})

    def test_empty(self) -> None:
        empty = self.rpc.test_empty()
        self.assertIsInstance(empty, EmptyModel)

    def test_basic(self) -> None:
        text = TextModel(title=None, body="  hello  ")
        basic = self.rpc.test_basic(text=text, flag=True, count=3, note="note")
        self.assertEqual(basic.body, "hello")
        self.assertEqual(basic.title, "note")

    def test_list_map(self) -> None:
        nested = self.rpc.test_list_map(
            texts=[TextModel(title="t1", body="b1"), TextModel(title="t2", body="b2")],
            flags={"mode": "fast"},
        )
        self.assertIsNotNone(nested.flags)
        self.assertEqual(nested.flags.retries, 2)
        self.assertEqual(nested.flags.meta.get("mode"), "fast")
        self.assertIsInstance(nested.lookup.get("first"), TextModel)

    def test_optional(self) -> None:
        optional = self.rpc.test_optional(text=None, flag=None)
        self.assertFalse(optional.enabled)

    def test_validation_error(self) -> None:
        with self.assertRaises(ValidationRPCError):
            self.rpc.test_validation_error(text=TextModel(title=None, body=""))

    def test_input_error(self) -> None:
        with self.assertRaises(InputRPCError):
            self.rpc.test_basic(text="bad", flag=True, count=1, note=None)

    def test_unauthorized_error(self) -> None:
        with self.assertRaises(UnauthorizedRPCError):
            self.rpc.test_unauthorized_error()

    def test_auth_middleware_missing_token(self) -> None:
        rpc = RPCClient("http://localhost:8080")
        with self.assertRaises(UnauthorizedRPCError):
            rpc.test_empty()

    def test_auth_middleware_invalid_token(self) -> None:
        rpc = RPCClient(
            "http://localhost:8080",
            headers={"Authorization": "Bearer bad_token"},
        )
        with self.assertRaises(UnauthorizedRPCError):
            rpc.test_empty()

    def test_forbidden_error(self) -> None:
        with self.assertRaises(ForbiddenRPCError):
            self.rpc.test_forbidden_error()

    def test_not_implemented_error(self) -> None:
        with self.assertRaises(NotImplementedRPCError):
            self.rpc.test_not_implemented_error()

    def test_custom_error(self) -> None:
        with self.assertRaises(CustomRPCError):
            self.rpc.test_custom_error()

    def test_map_return(self) -> None:
        mapped = self.rpc.test_map_return()
        self.assertIsInstance(mapped, dict)
        self.assertIsInstance(mapped.get("a"), TextModel)
        self.assertEqual(mapped["a"].body, "mapped")

    def test_json(self) -> None:
        payload = {"count": 2, "tags": ["a", "b"]}
        result = self.rpc.test_json(data=payload)
        self.assertIsInstance(result, dict)
        self.assertEqual(result.get("count"), 2)

    def test_raw(self) -> None:
        payload = {"ok": True}
        result = self.rpc.test_raw(payload=payload)
        self.assertIsInstance(result, dict)
        self.assertEqual(result.get("ok"), True)

    def test_mixed_payload(self) -> None:
        payload = {"data": {"value": "x"}, "raw_data": {"id": 1}}
        result = self.rpc.test_mixed_payload(payload=payload)
        self.assertIsInstance(result, PayloadModel)
        self.assertEqual(result.data.get("value"), "x")
        self.assertEqual(result.raw_data.get("id"), 1)


if __name__ == "__main__":
    unittest.main()
