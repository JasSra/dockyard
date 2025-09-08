# Agent

Go HTTP server exposing Agent API. Uses templates under /templates and stores state in /data.

Endpoints (HMAC with AGENT_SHARED_SECRET):
- POST /v1/deploy
- POST /v1/compose/up
- POST /v1/compose/down
- POST /v1/nginx/apply
- POST /v1/ports/allocate
- POST /v1/ports/free
- GET  /v1/status
- GET  /v1/logs
- POST /v1/ops/run

Dev
- Run unit tests locally: `go test ./...`
- If your environment blocks local runs, you can validate in a container build:
	- docker build . (go will download modules during build)
