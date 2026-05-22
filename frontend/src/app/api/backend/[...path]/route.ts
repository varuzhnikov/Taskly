import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:8080";
const PROXY_PREFIX = "/api/backend";

function buildBackendUrl(req: NextRequest): string {
  const url = new URL(req.url);
  const backendPath = url.pathname.startsWith(PROXY_PREFIX)
    ? url.pathname.slice(PROXY_PREFIX.length)
    : url.pathname;
  return `${BACKEND_URL}${backendPath}${url.search}`;
}

function copyRequestHeaders(req: NextRequest): Headers {
  const headers = new Headers(req.headers);
  headers.delete("host");
  headers.delete("connection");
  headers.delete("content-length");
  return headers;
}

async function readRequestBody(req: NextRequest): Promise<string | undefined> {
  if (req.method === "GET" || req.method === "HEAD") {
    return undefined;
  }

  const body = await req.text();
  return body === "" ? undefined : body;
}

async function proxy(req: NextRequest): Promise<Response> {
  const backendRes = await fetch(buildBackendUrl(req), {
    method: req.method,
    headers: copyRequestHeaders(req),
    body: await readRequestBody(req),
    cache: "no-store",
  }).catch(() => null);

  if (!backendRes) {
    return NextResponse.json({ error: "Backend unavailable" }, { status: 503 });
  }

  const headers = new Headers(backendRes.headers);
  headers.delete("content-length");

  return new Response(backendRes.body, {
    status: backendRes.status,
    statusText: backendRes.statusText,
    headers,
  });
}

export async function GET(req: NextRequest) {
  return proxy(req);
}

export async function POST(req: NextRequest) {
  return proxy(req);
}

export async function PUT(req: NextRequest) {
  return proxy(req);
}

export async function PATCH(req: NextRequest) {
  return proxy(req);
}

export async function DELETE(req: NextRequest) {
  return proxy(req);
}

export async function OPTIONS(req: NextRequest) {
  return proxy(req);
}
