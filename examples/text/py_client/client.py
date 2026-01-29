#!/usr/bin/env python3

from rpcclient import RPCClient, TextModel


def main() -> None:
    rpc = RPCClient("http://localhost:8080")

    text = TextModel(
        title="Sample",
        data="Hello world. This is a short example.",
    )
    text_id = rpc.submit_text(text=text)
    stats = rpc.compute_stats(text_id=text_id)
    print("stats:", stats)


if __name__ == "__main__":
    main()
