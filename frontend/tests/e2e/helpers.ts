import { type Page, type APIRequestContext, expect } from "@playwright/test";

const BACKEND = "http://localhost:8080";

export const TEST_USER = {
  email: `e2e_${Date.now()}@test.local`,
  password: "testpassword123",
};

export async function registerUser(
  request: APIRequestContext,
  email: string,
  password: string,
) {
  const res = await request.post(`${BACKEND}/api/v1/auth/register`, {
    data: { email, password },
  });
  if (!res.ok()) {
    throw new Error(`Registration failed: ${await res.text()}`);
  }
  return res.json();
}

export async function loginUser(
  request: APIRequestContext,
  email: string,
  password: string,
) {
  const res = await request.post(`${BACKEND}/api/v1/auth/login`, {
    data: { email, password },
  });
  if (!res.ok()) {
    throw new Error(`Login failed: ${await res.text()}`);
  }
  return res.json();
}

export async function createProject(
  request: APIRequestContext,
  accessToken: string,
  name: string,
) {
  const res = await request.post(`${BACKEND}/api/v1/projects`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
    data: { name, color: "#db4035" },
  });
  if (!res.ok()) {
    throw new Error(`Project creation failed: ${await res.text()}`);
  }
  return res.json() as Promise<{ id: string; name: string }>;
}

export async function loginViaUI(page: Page, email: string, password: string) {
  await page.goto("/login");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /sign in/i }).click();
  await expect(page).toHaveURL(/\/inbox/, { timeout: 30_000 });
}
