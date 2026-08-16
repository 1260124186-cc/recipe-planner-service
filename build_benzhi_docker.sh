#!/usr/bin/env bash
set -euo pipefail

platform="${1:-linux/arm64}"
image="recipe-planner-service:benzhi"
container="recipe-planner-benzhi-check"
port="18080"

cleanup() {
  docker rm -f "${container}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

docker build --platform "${platform}" -f benzhi.Dockerfile -t "${image}" .
docker run --rm --platform "${platform}" --entrypoint go "${image}" build ./...
docker run -d --rm --platform "${platform}" --name "${container}" -p "${port}:8080" "${image}" >/dev/null

for _ in $(seq 1 20); do
  if curl --fail --silent --show-error "http://127.0.0.1:${port}/health"; then
    exit 0
  fi
  sleep 1
done

echo "health check did not become ready" >&2
exit 1
