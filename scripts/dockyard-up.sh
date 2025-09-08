#!/usr/bin/env bash
# Boot the Dockyard stack with simple flags.
# Writes .env overrides for docker compose and starts services.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INFRA_DIR="$ROOT_DIR/infra"
ENV_FILE="$INFRA_DIR/.env"

NGINX_PORT=${NGINX_PORT:-80}
CONTROLLER_PORT=${CONTROLLER_PORT:-3000}
AGENT_PORT=${AGENT_PORT:-8080}
DOMAIN=${DOMAIN:-localhost}
APPS_ROOT=${APPS_ROOT:-apps}
ADMIN_SECRET=${ADMIN_SECRET:-admin}
AGENT_SHARED_SECRET=${AGENT_SHARED_SECRET:-devsecret}

usage(){
  cat <<EOF
Usage: $0 [--nginx-port N] [--controller-port N] [--agent-port N] [--domain NAME] [--apps-root NAME] [--admin-secret S] [--agent-secret S]

Examples:
  $0 --nginx-port 8081 --controller-port 4000 --agent-port 9090 --domain localhost --apps-root apps
EOF
}

# Parse flags
while [[ $# -gt 0 ]]; do
  case "$1" in
    --nginx-port) NGINX_PORT="$2"; shift 2;;
    --controller-port) CONTROLLER_PORT="$2"; shift 2;;
    --agent-port) AGENT_PORT="$2"; shift 2;;
    --domain) DOMAIN="$2"; shift 2;;
    --apps-root) APPS_ROOT="$2"; shift 2;;
    --admin-secret) ADMIN_SECRET="$2"; shift 2;;
    --agent-secret) AGENT_SHARED_SECRET="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown option: $1" >&2; usage; exit 1;;
  esac
done

mkdir -p "$INFRA_DIR"
cat > "$ENV_FILE" <<ENV
NGINX_PORT=$NGINX_PORT
CONTROLLER_PORT=$CONTROLLER_PORT
AGENT_PORT=$AGENT_PORT
DOMAIN=$DOMAIN
APPS_ROOT=$APPS_ROOT
ADMIN_SECRET=$ADMIN_SECRET
AGENT_SHARED_SECRET=$AGENT_SHARED_SECRET
ENV

echo "[info] Wrote $ENV_FILE" >&2

cd "$INFRA_DIR"
exec docker compose up -d
