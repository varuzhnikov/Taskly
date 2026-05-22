import { test, expect } from "@playwright/test";
import { TEST_USER, registerUser } from "./helpers";

const user = {
  email: `e2e_auth_${Date.now()}@test.local`,
  password: "testpassword123",
};

test.describe("Auth — redirect & validation", () => {
  test("unauthenticated: / redirects to /login", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL(/\/login/);
    await page.screenshot({ path: "test-screenshots/01-redirect-to-login.png", fullPage: true });
  });

  test("unauthenticated: /inbox redirects to /login", async ({ page }) => {
    await page.goto("/inbox");
    await expect(page).toHaveURL(/\/login/);
    await page.screenshot({ path: "test-screenshots/02-inbox-redirects-to-login.png", fullPage: true });
  });

  test("login page renders correctly", async ({ page }) => {
    await page.goto("/login");
    await expect(page.getByLabel("Email")).toBeVisible();
    await expect(page.getByLabel("Password")).toBeVisible();
    await expect(page.getByRole("button", { name: /sign in/i })).toBeVisible();
    await expect(page.getByRole("link", { name: /create one/i })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/03-login-page.png", fullPage: true });
  });

  test("register page renders correctly", async ({ page }) => {
    await page.goto("/register");
    await expect(page.getByLabel("Email")).toBeVisible();
    await expect(page.getByLabel(/^password$/i)).toBeVisible();
    await expect(page.getByLabel(/confirm password/i)).toBeVisible();
    await expect(page.getByRole("button", { name: /create account/i })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/04-register-page.png", fullPage: true });
  });

  test("login: invalid email shows validation error", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("Email").fill("notanemail");
    await page.getByLabel("Password").fill("password123");
    await page.getByRole("button", { name: /sign in/i }).click();
    await expect(page.getByText(/invalid email/i)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/05-login-invalid-email.png", fullPage: true });
  });

  test("login: short password shows validation error", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("Email").fill("test@example.com");
    await page.getByLabel("Password").fill("short");
    await page.getByRole("button", { name: /sign in/i }).click();
    await expect(page.getByText(/at least 8 characters/i)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/06-login-short-password.png", fullPage: true });
  });

  test("register: mismatched passwords shows error", async ({ page }) => {
    await page.goto("/register");
    await page.getByLabel("Email").fill("test@example.com");
    await page.getByLabel(/^password$/i).fill("securepassword");
    await page.getByLabel(/confirm password/i).fill("different");
    await page.getByRole("button", { name: /create account/i }).click();
    await expect(page.getByText(/passwords do not match/i)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/07-register-password-mismatch.png", fullPage: true });
  });

  test("register: wrong credentials show server error", async ({ page }) => {
    await page.goto("/login");
    await page.getByLabel("Email").fill("nobody@example.com");
    await page.getByLabel("Password").fill("wrongpassword");
    await page.getByRole("button", { name: /sign in/i }).click();
    await expect(page.getByText(/invalid|unauthorized|incorrect/i)).toBeVisible({ timeout: 8000 });
    await page.screenshot({ path: "test-screenshots/08-login-wrong-credentials.png", fullPage: true });
  });

  test("navigation: login → register → login", async ({ page }) => {
    await page.goto("/login");
    await page.getByRole("link", { name: /create one/i }).click();
    await expect(page).toHaveURL(/\/register/);
    await page.getByRole("link", { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/login/);
    await page.screenshot({ path: "test-screenshots/09-login-register-nav.png", fullPage: true });
  });
});

test.describe("Auth — full registration flow", () => {
  test("can register and land on inbox", async ({ page }) => {
    await page.goto("/register");
    await page.getByLabel("Email").fill(user.email);
    await page.getByLabel(/^password$/i).fill(user.password);
    await page.getByLabel(/confirm password/i).fill(user.password);
    await page.screenshot({ path: "test-screenshots/10-register-filled.png", fullPage: true });
    await page.getByRole("button", { name: /create account/i }).click();
    await expect(page).toHaveURL(/\/inbox/, { timeout: 30_000 });
    await page.screenshot({ path: "test-screenshots/11-after-register-inbox.png", fullPage: true });
  });

  test("registered user can log out and back in", async ({ page, request }) => {
    // Register via API so we have a clean user
    await registerUser(request, user.email, user.password).catch(() => {
      // user might already exist from previous test run, that's fine
    });

    // Login via UI
    await page.goto("/login");
    await page.getByLabel("Email").fill(user.email);
    await page.getByLabel("Password").fill(user.password);
    await page.getByRole("button", { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/inbox/, { timeout: 30_000 });
    await page.screenshot({ path: "test-screenshots/12-logged-in-inbox.png", fullPage: true });

    // Logout
    await page.getByRole("button", { name: /sign out/i }).click();
    await expect(page).toHaveURL(/\/login/, { timeout: 8_000 });
    await page.screenshot({ path: "test-screenshots/13-after-logout.png", fullPage: true });

    // Login again
    await page.getByLabel("Email").fill(user.email);
    await page.getByLabel("Password").fill(user.password);
    await page.getByRole("button", { name: /sign in/i }).click();
    await expect(page).toHaveURL(/\/inbox/, { timeout: 30_000 });
    await page.screenshot({ path: "test-screenshots/14-login-again.png", fullPage: true });
  });
});
