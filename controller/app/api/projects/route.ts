import { NextRequest, NextResponse } from 'next/server'
import { createProject } from '../../../lib/store'

export async function POST(req: NextRequest) {
  const body = await req.json()
  const p = createProject(body.name)
  return NextResponse.json(p)
}
