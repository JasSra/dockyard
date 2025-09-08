import { NextResponse } from 'next/server'
import { getProject } from '../../../../lib/store'
import { Agent } from '../../../../lib/agent'

export async function POST(_: Request, { params }: { params: { id: string }}) {
  const p = getProject(params.id)
  if (!p) return NextResponse.json({ error: 'not found' }, { status: 404 })
  await Agent.composeDown(p.id)
  await Agent.freePort(p.id)
  return NextResponse.json({ ok: true })
}
