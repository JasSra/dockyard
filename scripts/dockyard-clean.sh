#!/usr/bin/env bash
# Stops stack, removes generated nginx confs, and prunes containers/images (optional flags).
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INFRA_DIR="$ROOT_DIR/infra"
NGINX_DIR="$INFRA_DIR/nginx/conf.d"
DO_PRUNE=${DO_PRUNE:-0}

cd "$INFRA_DIR"
docker compose down || true

if [[ -d "$NGINX_DIR" ]]; then
  echo "[info] Cleaning conf.d/*.conf" >&2
  find "$NGINX_DIR" -maxdepth 1 -type f -name '*.conf' -print -delete || true
fi

if [[ "$DO_PRUNE" == "1" ]]; then
  echo "[warn] Pruning dangling containers/images/volumes" >&2
  docker system prune -f || true
fi

echo "[ok] Clean complete" >&2
