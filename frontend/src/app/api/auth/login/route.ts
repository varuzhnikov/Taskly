import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:8080";
const REFRESH_TOKEN_MAX_AGE = 7 * 24 * 60 * 60; // 7 days in seconds

export async function POST(req: NextRequest) {
  const body = await req.json().catch(() => null);
  if (!body) {
    return NextResponse.json({ error: "Invalid request body" }, { status: 400 });
  }

  const backendRes = await fetch(`${BACKEND_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  }).catch(() => null);

  if (!backendRes) {
    return NextResponse.json({ error: "Backend unavailable" }, { status: 503 });
  }

  const data = await backendRes.json().catch(() => null);

  if (!backendRes.ok) {
    return NextResponse.json(data ?? { error: "Login failed" }, {
      status: backendRes.status,
    });
  }

  const { refresh_token, access_token, user } = data as {
    refresh_token: string;
    access_token: string;
    user: { id: string; email: string };
  };

  const res = NextResponse.json({ accessToken: access_token, user });
  res.cookies.set("refresh_token", refresh_token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: REFRESH_TOKEN_MAX_AGE,
  });
  return res;
}
