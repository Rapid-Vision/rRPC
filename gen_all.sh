#!/usr/bin/env bash

set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

run_make() {
  local dir="$1"
  echo "==> ${dir}"
  (cd "${dir}" && make)
}

run_make "${root_dir}/examples/hello_world"
run_make "${root_dir}/examples/text"
run_make "${root_dir}/integration_test"

echo "All make targets completed."
