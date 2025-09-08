Dockyard Controller

Dev quickstart
- Install deps: npm install
- Run in dev: npm run dev
- Typecheck: npm run typecheck

Environment
- ADMIN_SECRET (default: admin)
- AGENT_URL (default: http://localhost:8080)
- DOMAIN (default: {{ DOMAIN }})
- APPS_ROOT (default: {{ APPS_ROOT }})

Notes
- If you see errors like "Cannot find module 'next/server'" or Node types missing, run npm install to generate node_modules and types.# Controller

Minimal Next.js app with API stubs:
- POST /api/projects
- POST /api/projects/:id/deploy
- POST /api/projects/:id/kv
- GET  /api/projects/:id/kv
- GET  /api/projects/:id/status

Dev run:
- Use docker-compose in ../../infra
