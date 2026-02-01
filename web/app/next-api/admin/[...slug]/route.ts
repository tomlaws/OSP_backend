import { NextRequest, NextResponse } from 'next/server';

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  return handleRequest(request, await params, 'GET');
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  return handleRequest(request, await params, 'POST');
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  return handleRequest(request, await params, 'DELETE');
}

async function handleRequest(request: NextRequest, params: { slug: string[] }, method: string) {
  const slug = params.slug.join('/');
  const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL;
  const token = process.env.ROOT_TOKEN || '';

  const url = `${backendUrl}/admin/${slug}${request.nextUrl.search}`;

  const headers = new Headers(request.headers);
  headers.set('Authorization', `Bearer ${token}`);
  // Remove host header to avoid issues
  headers.delete('host');

  try {
    const options: RequestInit = {
      method,
      headers,
    };

    if (method !== 'GET' && method !== 'HEAD') {
      const text = await request.text();
      if (text) {
        options.body = text;
      }
    }

    const response = await fetch(url, options);
    
    // Pass through the response
    const data = await response.text();
    
    // Forward status and headers
    return new NextResponse(data, {
      status: response.status,
      headers: {
        'Content-Type': response.headers.get('Content-Type') || 'application/json',
      },
    });
  } catch (error) {
    console.error('Proxy Error:', error);
    return NextResponse.json(
      { error: 'Internal Server Error' },
      { status: 500 }
    );
  }
}
