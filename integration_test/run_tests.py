#!/usr/bin/env python3

import os
import signal
import socket
import subprocess
import sys
import time
from pathlib import Path


def run(cmd: list[str], cwd: Path) -> None:
    subprocess.run(cmd, cwd=cwd, check=True)


def wait_for_port(host: str, port: int, timeout: float) -> None:
    deadline = time.time() + timeout
    while time.time() < deadline:
        try:
            with socket.create_connection((host, port), timeout=0.2):
                return
        except OSError:
            time.sleep(0.1)
    raise RuntimeError(f"server did not start within {timeout:.1f}s")


def main() -> int:
    root = Path(__file__).resolve().parents[1]
    workdir = Path(__file__).resolve().parent
    rrpc = root / "rRPC"

    run(
        ["go", "build", "-o", "rRPC"],
        cwd=root,
    )

    run(
        [str(rrpc), "server", "-o", "./go_server", "-f", "test.rrpc"],
        cwd=workdir,
    )
    run(
        [str(rrpc), "client", "--lang", "go", "-o", "./go_client", "-f", "test.rrpc"],
        cwd=workdir,
    )
    run(
        [str(rrpc), "client", "-o", "./py_client", "-f", "test.rrpc"],
        cwd=workdir,
    )
    run(
        [str(rrpc), "openapi", "-o", ".", "-f", "test.rrpc"],
        cwd=workdir,
    )

    server = subprocess.Popen(
        ["go", "run", "."],
        cwd=workdir / "go_server",
        start_new_session=True,
    )
    try:
        wait_for_port("127.0.0.1", 8080, timeout=5.0)
        run(["go", "test", "."], cwd=workdir / "go_client")
        run(
            [sys.executable, "-m", "unittest", "test_client.py"],
            cwd=workdir / "py_client",
        )
    finally:
        try:
            os.killpg(server.pid, signal.SIGTERM)
        except ProcessLookupError:
            sys.exit(1)
        try:
            server.wait(timeout=3.0)
        except subprocess.TimeoutExpired:
            try:
                os.killpg(server.pid, signal.SIGKILL)
            except ProcessLookupError:
                sys.exit(1)
            server.wait(timeout=3.0)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
