import crypto from 'crypto'
import { CONFIG } from './config'

function hmac(body: string) {
  const key = process.env.AGENT_SHARED_SECRET || 'devsecret'
  return crypto.createHmac('sha256', key).update(body).digest('hex')
}

async function post(path: string, payload: any) {
  const body = JSON.stringify(payload)
  const res = await fetch(`${CONFIG.AGENT_URL}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Signature': hmac(body),
      'X-Idempotency-Key': crypto.randomUUID(),
    },
    body,
    cache: 'no-store',
  })
  if (!res.ok) throw new Error(`agent ${path} ${res.status}`)
  return res.json()
}

export const Agent = {
  allocatePort: (projectId: string) => post('/v1/ports/allocate', { projectId }),
  freePort: (projectId: string) => post('/v1/ports/free', { projectId }),
  applyNginx: (projectId: string, host: string, upstreamPort: number) => post('/v1/nginx/apply', { projectId, host, upstreamPort }),
  deploy: (projectId: string, archiveUrl: string) => post('/v1/deploy', { projectId, archiveUrl }),
  async composeUp(projectId: string, imageTag: string, env: Record<string,string> = {}, containerPort = 3000, hostPort?: number) {
    const body: any = { projectId, imageTag, env, containerPort }
    if (hostPort) body.hostPort = hostPort
    return post('/v1/compose/up', body)
  },
  async composeDown(projectId: string) {
    const body = { projectId }
    return post('/v1/compose/down', body)
  },
}
