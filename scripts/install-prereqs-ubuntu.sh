#!/usr/bin/env bash
# Installs Docker Engine (community), docker compose plugin, and helpers on Ubuntu.
# Safe to run multiple times.
set -euo pipefail

if [[ $(id -u) -ne 0 ]]; then
  echo "[info] Not running as root; will use sudo where needed." >&2
  SUDO=sudo
else
  SUDO=""
fi

# Update apt and install dependencies
$SUDO apt-get update -y
$SUDO apt-get install -y ca-certificates curl gnupg lsb-release jq

# Install Docker Engine (Ubuntu official docker.io is fine for dev)
if ! command -v docker >/dev/null 2>&1; then
  echo "[info] Installing docker.io (Ubuntu)" >&2
  $SUDO apt-get install -y docker.io
fi

# Install docker compose plugin (package: docker-compose-plugin)
if ! docker compose version >/dev/null 2>&1; then
  echo "[info] Installing docker-compose-plugin" >&2
  $SUDO apt-get install -y docker-compose-plugin
fi

# Enable and start docker
$SUDO systemctl enable --now docker || true

# Add current user to docker group so compose can run without sudo
if getent group docker >/dev/null 2>&1; then
  if id -nG "$USER" | grep -qw docker; then
    echo "[info] User $USER already in docker group." >&2
  else
    echo "[info] Adding $USER to docker group (you must log out/in to take effect)" >&2
    $SUDO usermod -aG docker "$USER" || true
  fi
fi

echo "[ok] Prerequisites installed. You may need to re-login for docker group membership to apply." >&2
