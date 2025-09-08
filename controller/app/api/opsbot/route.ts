import { NextRequest, NextResponse } from 'next/server'

// Placeholder endpoint receiving OpsBot tool calls. // { SPECULATION }
export async function POST(req: NextRequest) {
  const body = await req.json()
  return NextResponse.json({ received: body })
}
