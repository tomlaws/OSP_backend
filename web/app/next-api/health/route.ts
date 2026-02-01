import { NextResponse } from 'next/server';

export const dynamic = 'force-dynamic';

export async function GET() {
  return NextResponse.json({ 
    status: 'ok',
    message: 'Next.js frontend is healthy',
    timestamp: new Date().toISOString()
  });
}
