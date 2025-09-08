import { NextRequest, NextResponse } from 'next/server'
import { getProject, setKV } from '../../../../lib/store'

export async function GET(req: NextRequest, { params }: { params: { id: string }}) {
  const p = getProject(params.id)
  return NextResponse.json(p?.kv || {})
}

export async function POST(req: NextRequest, { params }: { params: { id: string }}) {
  const body = await req.json()
  setKV(params.id, body || {})
  return NextResponse.json({ ok: true })
}
