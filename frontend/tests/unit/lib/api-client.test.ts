import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

const mockClearAuth = vi.fn();
const mockSetAccessToken = vi.fn();

vi.mock("@/store/auth", () => ({
  useAuthStore: {
    getState: vi.fn().mockReturnValue({
      accessToken: "test-token",
      user: null,
      isAuthenticated: true,
      setAuth: vi.fn(),
      setAccessToken: mockSetAccessToken,
      clearAuth: mockClearAuth,
    }),
  },
}));

describe("api client", () => {
  const originalFetch = globalThis.fetch;

  afterEach(() => {
    globalThis.fetch = originalFetch;
    vi.clearAllMocks();
  });

  it("attaches Authorization header when token is present", async () => {
    let capturedInit: RequestInit | undefined;
    globalThis.fetch = vi.fn().mockImplementation((_url: string, init: RequestInit) => {
      capturedInit = init;
      return Promise.resolve(
        new Response(JSON.stringify({ id: "1" }), { status: 200 }),
      );
    });

    const { api } = await import("@/lib/api/client");
    await api.get("/api/v1/tasks/");

    const headers = capturedInit?.headers as Record<string, string>;
    expect(headers["Authorization"]).toBe("Bearer test-token");
  });

  it("does not attach Authorization header when token is null", async () => {
    const { useAuthStore } = await import("@/store/auth");
    vi.mocked(useAuthStore.getState).mockReturnValue({
      accessToken: null,
      user: null,
      isAuthenticated: false,
      setAuth: vi.fn(),
      setAccessToken: vi.fn(),
      clearAuth: mockClearAuth,
    });

    let capturedInit: RequestInit | undefined;
    globalThis.fetch = vi.fn().mockImplementation((_url: string, init: RequestInit) => {
      capturedInit = init;
      return Promise.resolve(
        new Response(JSON.stringify([]), { status: 200 }),
      );
    });

    const { api } = await import("@/lib/api/client");
    await api.get("/api/v1/tasks/");

    const headers = capturedInit?.headers as Record<string, string>;
    expect(headers["Authorization"]).toBeUndefined();
  });

  it("throws ApiError with status on non-ok response", async () => {
    globalThis.fetch = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ error: "Not found" }), { status: 404 }),
    );

    const { api, ApiError } = await import("@/lib/api/client");

    await expect(api.get("/api/v1/tasks/nonexistent")).rejects.toBeInstanceOf(
      ApiError,
    );
  });
});
