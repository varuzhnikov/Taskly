import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL ?? "http://localhost:8080";
const REFRESH_TOKEN_MAX_AGE = 7 * 24 * 60 * 60;

export async function POST(req: NextRequest) {
  const refreshToken = req.cookies.get("refresh_token")?.value;
  if (!refreshToken) {
    return NextResponse.json({ error: "No refresh token" }, { status: 401 });
  }

  const backendRes = await fetch(`${BACKEND_URL}/api/v1/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  }).catch(() => null);

  if (!backendRes) {
    return NextResponse.json({ error: "Backend unavailable" }, { status: 503 });
  }

  const data = await backendRes.json().catch(() => null);

  if (!backendRes.ok) {
    const res = NextResponse.json(data ?? { error: "Token refresh failed" }, {
      status: backendRes.status,
    });
    res.cookies.delete("refresh_token");
    return res;
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
