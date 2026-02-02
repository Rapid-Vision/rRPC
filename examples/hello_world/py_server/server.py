from rpcserver import create_app
from rpcserver import RPCHandlers
from rpcserver import GreetingMessageModel


class Service(RPCHandlers):
    def hello_world(
        self, name: str, surname: str | None = None
    ) -> GreetingMessageModel:
        if surname:
            full_name = f"{name} {surname}"
        else:
            full_name = name
        return GreetingMessageModel(message=f"Hello, {full_name}!")


app = create_app(Service())

if __name__ == "__main__":
    import uvicorn

    uvicorn.run("server:app", host="127.0.0.1", port=8080, reload=False)
