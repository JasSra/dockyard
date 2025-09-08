import { NextResponse } from 'next/server'
import { getProject } from '../../../../lib/store'
import { Agent } from '../../../../lib/agent'
import { CONFIG } from '../../../../lib/config'

export async function POST(_: Request, { params }: { params: { id: string }}) {
  const p = getProject(params.id)
  if (!p) return NextResponse.json({ error: 'not found' }, { status: 404 })
  // Allocate a port and write nginx
  const { hostPort } = await Agent.allocatePort(p.id)
  const host = `${p.id}.${CONFIG.APPS_ROOT}.${CONFIG.DOMAIN}`
  // bring up container using the standard tag built by deploy
  await Agent.composeUp(p.id, `dy-${p.id}:latest`, p.kv || {}, 3000, hostPort)
  await Agent.applyNginx(p.id, host, hostPort)
  return NextResponse.json({ ok: true, url: `http://${host}` })
}
