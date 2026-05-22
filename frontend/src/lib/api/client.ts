import { useAuthStore } from "@/store/auth";

function resolveApiBaseUrl(): string {
  const configuredUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!configuredUrl) {
    return "/api/backend";
  }

  try {
    const url = new URL(configuredUrl);
    if (url.hostname === "localhost" || url.hostname === "127.0.0.1") {
      return "/api/backend";
    }
  } catch {
    return configuredUrl;
  }

  return configuredUrl;
}

const API_BASE_URL = resolveApiBaseUrl();

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

function toCamelCase(s: string): string {
  return s.replace(/_([a-z])/g, (_, c: string) => c.toUpperCase());
}

function transformKeys(value: unknown): unknown {
  if (Array.isArray(value)) return value.map(transformKeys);
  if (value !== null && typeof value === "object") {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>).map(([k, v]) => [
        toCamelCase(k),
        transformKeys(v),
      ]),
    );
  }
  return value;
}

let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  const res = await fetch("/api/auth/refresh", { method: "POST" });
  if (!res.ok) {
    useAuthStore.getState().clearAuth();
    throw new ApiError(res.status, "Session expired");
  }
  const data = (await res.json()) as {
    accessToken: string;
    user: { id: string; email: string };
  };
  useAuthStore.getState().setAuth(data.accessToken, data.user);
  return data.accessToken;
}

async function getRefreshedToken(): Promise<string> {
  if (!refreshPromise) {
    refreshPromise = refreshAccessToken().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

async function request<T>(
  path: string,
  method: HttpMethod,
  body?: unknown,
  retry = true,
): Promise<T> {
  const token = useAuthStore.getState().accessToken;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (res.status === 401 && retry) {
    const newToken = await getRefreshedToken();
    return request<T>(path, method, body, false);
    void newToken;
  }

  if (!res.ok) {
    let message = res.statusText;
    try {
      const err = (await res.json()) as { error?: string };
      if (err.error) message = err.error;
    } catch {
      // ignore parse failure
    }
    throw new ApiError(res.status, message);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return transformKeys(await res.json()) as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path, "GET"),
  post: <T>(path: string, body?: unknown) => request<T>(path, "POST", body),
  put: <T>(path: string, body: unknown) => request<T>(path, "PUT", body),
  patch: <T>(path: string, body?: unknown) => request<T>(path, "PATCH", body),
  delete: <T>(path: string) => request<T>(path, "DELETE"),
};
