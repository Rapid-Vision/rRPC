import unittest
from unittest import mock

from pydantic import ValidationError

from rpclient_pydantic import (
    RPCClient,
    FlagsModel,
)


class RPCClientPydanticTest(unittest.TestCase):
    def test_validates_inputs_before_request(self) -> None:
        rpc = RPCClient("http://localhost:8080")
        with mock.patch.object(
            RPCClient,
            "_request",
            side_effect=AssertionError("request should not be called"),
        ):
            with self.assertRaises(ValidationError):
                rpc.test_basic(
                    text={"title": "missing-body"},
                    flag=True,
                    count=1,
                )

    def test_accepts_optional_nullable_fields(self) -> None:
        rpc = RPCClient(
            "http://localhost:8080", headers={"Authorization": "Bearer test_token"}
        )
        result = rpc.test_optional(text=None, flag=None)
        self.assertIsInstance(result, FlagsModel)


if __name__ == "__main__":
    unittest.main()
