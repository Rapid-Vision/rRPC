from generated import RPCClient


def main() -> None:
    rpc = RPCClient("http://localhost:8080")
    greeting = rpc.hello_world(name="Ada", surname="Lovelace")
    print("greeting:", greeting)


if __name__ == "__main__":
    main()
