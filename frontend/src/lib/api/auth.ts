import type { SessionUser } from "@/lib/types";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  accessToken: string;
  user: SessionUser;
}

export async function loginViaRouteHandler(
  data: LoginRequest,
): Promise<AuthResponse> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Login failed");
  }

  return res.json() as Promise<AuthResponse>;
}

export async function registerViaRouteHandler(
  data: RegisterRequest,
): Promise<AuthResponse> {
  const res = await fetch("/api/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Registration failed");
  }

  return res.json() as Promise<AuthResponse>;
}

export async function logoutViaRouteHandler(): Promise<void> {
  await fetch("/api/auth/logout", { method: "POST" });
}
