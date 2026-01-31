#!/usr/bin/env python3

import argparse
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


def parse_tests() -> set[str]:
    parser = argparse.ArgumentParser(description="Run rRPC integration tests")
    parser.add_argument(
        "--test",
        action="append",
        default=[],
        help="Comma-separated list of test suites: go, py, ts-all, ts-bare, ts-zod",
    )
    args = parser.parse_args()
    all_tests = {"go", "py", "ts-all", "ts-bare", "ts-zod"}
    if not args.test:
        return all_tests

    selected: set[str] = set()
    for entry in args.test:
        selected.update({item.strip() for item in entry.split(",") if item.strip()})

    invalid = selected.difference(all_tests)
    if invalid:
        raise SystemExit(f"unknown test suites: {', '.join(invalid)}")
    return selected


def codegen(
    rrpc: str,
    workdir: os.PathLike,
    root: os.PathLike,
):
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
        [
            str(rrpc),
            "client",
            "--lang",
            "ts",
            "-o",
            "./ts_client",
            "-f",
            "test.rrpc",
        ],
        cwd=workdir,
    )
    run(
        [
            str(rrpc),
            "client",
            "--lang",
            "ts",
            "--ts-zod",
            "--pkg",
            "rpcclientzod",
            "-o",
            "./ts_client",
            "-f",
            "test.rrpc",
        ],
        cwd=workdir,
    )

    run(
        [str(rrpc), "openapi", "-o", ".", "-f", "test.rrpc"],
        cwd=workdir,
    )


def main() -> int:
    selected = parse_tests()
    run_go = "go" in selected
    run_py = "py" in selected
    run_ts_all = "ts-all" in selected
    run_ts_bare = run_ts_all or "ts-bare" in selected
    run_ts_zod = run_ts_all or "ts-zod" in selected

    root = Path(__file__).resolve().parents[1]
    workdir = Path(__file__).resolve().parent
    rrpc = root / "rRPC"

    codegen(rrpc=rrpc, workdir=workdir, root=root)

    server = subprocess.Popen(
        ["go", "run", "."],
        cwd=workdir / "go_server",
        start_new_session=True,
    )
    try:
        wait_for_port("127.0.0.1", 8080, timeout=5.0)
        if run_go:
            print("Running go tests:")
            run(["go", "test", "."], cwd=workdir / "go_client")
            print("\n")

        if run_py:
            print("Running python tests:")
            run(
                [sys.executable, "-m", "unittest", "test_client.py"],
                cwd=workdir / "py_client",
            )
            print("\n")

        if run_ts_bare or run_ts_zod:
            print("Running typescript tests:")
            run(["bun", "install"], cwd=workdir / "ts_client")
            if run_ts_all:
                run(["bun", "test"], cwd=workdir / "ts_client")
            else:
                if run_ts_bare:
                    run(["bun", "test", "client.test.ts"], cwd=workdir / "ts_client")
                if run_ts_zod:
                    run(
                        ["bun", "test", "client_zod.test.ts"],
                        cwd=workdir / "ts_client",
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
