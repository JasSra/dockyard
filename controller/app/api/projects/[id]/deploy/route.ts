import { NextResponse } from 'next/server'
import { getProject } from '../../../../lib/store'
import { Agent } from '../../../../lib/agent'

export async function POST(req: Request, { params }: { params: { id: string }}) {
  const p = getProject(params.id)
  if (!p) return NextResponse.json({ error: 'not found' }, { status: 404 })
  const { archiveUrl } = await req.json()
  if (!archiveUrl) return NextResponse.json({ error: 'archiveUrl required' }, { status: 400 })
  const data = await Agent.deploy(p.id, archiveUrl)
  return NextResponse.json(data)
}
