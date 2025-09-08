# Dockyard — Minimal PaaS (Controller + Agent + CLI)

This is an MVP mini-PaaS with:

- Controller (Next.js) providing UI + API and OpsBot
- Agent (Go) running on a Docker/Nginx host
- CLI (Go) for project deploys and operations
- Template registry for common stacks

Key features

- Automatic template detection (Next.js SSR/static, Static HTML, ASP.NET Core 8/9).
- Per-project KV (in-memory with JSON snapshots).
- Single Nginx fronts all web apps. Random GUID subdomains or custom domains.
- Port allocator with persistence, TTL cleanup, and event log.
- OpsBot (Ollama/OpenAI) manages deployments and recovery playbooks via tools.

Quick start

- See each package README under `controller/`, `agent/`, and `cli/`.
- For development, use the dev docker-compose under `infra/`.

Dev quickstart

- Requirements: Docker. Optional: Node.js for local Controller dev, Go for running tests locally.
- One-liner Ubuntu setup:
    - `scripts/install-prereqs-ubuntu.sh`
    - `scripts/dockyard-up.sh --nginx-port 8080 --controller-port 3000 --agent-port 8088 --domain localhost --apps-root apps`
    - `scripts/dockyard-status.sh` (to see services)
    - `scripts/dockyard-down.sh` (to stop)
    - `scripts/dockyard-clean.sh` (to clear confs; set `DO_PRUNE=1` to prune images)
- Alternatively: from `infra/`, run `docker compose up` (or use the devcontainer tasks).
- Services:
    - nginx on :${NGINX_PORT:-80}
    - controller on :${CONTROLLER_PORT:-3000}
    - agent on :${AGENT_PORT:-8080}
- Environment you can override in compose:
    - `AGENT_SHARED_SECRET=devsecret`, `ADMIN_SECRET=admin`, `DOMAIN=localhost`, `APPS_ROOT=apps`
- First run steps:
    - Open <http://localhost:3000>
    - Create a project, then start it via the Start API/UI. This will allocate a port, write an nginx vhost, and trigger nginx reload. Compose up is a placeholder for now.
    - Visit <http://{project}.{APPS_ROOT}.{DOMAIN}> (e.g., <http://myproj.apps.localhost>)

Placeholders

- DOMAIN, APPS_ROOT, CONTROLLER_URL, ADMIN_SECRET, AGENT_SHARED_SECRET, RANDOM_GUID,
  DEFAULT_CONTAINER_PORT(3000), PORT_RANGE_START(20000), PORT_RANGE_END(29999),
  DEFAULT_TTL_MINUTES(1440), AGENT_PUBLIC_IP, AGENT_SSH_USER.

Acceptance targets
- Deploy templates and reach app at http://{GUID}.{APPS_ROOT}.{DOMAIN}
- KV injected as env and visible at `/env` route.
- OpsBot can fix 502 by reallocating port and reloading nginx.

License
MIT // { SPECULATION }
# dockyard