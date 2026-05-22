import { test, expect } from "@playwright/test";
import { createProject, loginUser, registerUser, loginViaUI } from "./helpers";

const user = {
  email: `e2e_nav_${Date.now()}@test.local`,
  password: "testpassword123",
};

test.describe("Navigation & layout", () => {
  test.beforeAll(async ({ request }) => {
    await registerUser(request, user.email, user.password);
  });

  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, user.email, user.password);
  });

  test("sidebar is visible on inbox", async ({ page }) => {
    await expect(page.getByRole("complementary").getByText("Todoist", { exact: true })).toBeVisible();
    await expect(page.getByRole("link", { name: /inbox/i })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/30-sidebar-visible.png", fullPage: true });
  });

  test("header shows the logged-in user email", async ({ page }) => {
    await expect(page.getByText(user.email)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/31-header-user-email.png", fullPage: true });
  });

  test("labels page navigates correctly", async ({ page }) => {
    await page.getByRole("link", { name: /all labels/i }).click();
    await expect(page).toHaveURL(/\/labels/);
    await expect(page.getByRole("heading", { name: "Labels" })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/32-labels-page.png", fullPage: true });
  });

  test("back to inbox from labels", async ({ page }) => {
    await page.goto("/labels");
    await page.getByRole("link", { name: /inbox/i }).click();
    await expect(page).toHaveURL(/\/inbox/);
    await page.screenshot({ path: "test-screenshots/33-back-to-inbox.png", fullPage: true });
  });

  test("page refresh preserves session (token hydration)", async ({ page }) => {
    await page.screenshot({ path: "test-screenshots/34-before-refresh.png", fullPage: true });
    await page.reload();
    // Should stay on inbox after reload (token hydrated from cookie)
    await expect(page).toHaveURL(/\/inbox/, { timeout: 10_000 });
    await page.screenshot({ path: "test-screenshots/35-after-refresh-still-inbox.png", fullPage: true });
  });

  test("authenticated user visiting /login is redirected to /inbox", async ({ page }) => {
    await page.goto("/login");
    // AuthLayout redirects authenticated users away
    await expect(page).toHaveURL(/\/inbox/, { timeout: 8_000 });
    await page.screenshot({ path: "test-screenshots/36-auth-login-redirect.png", fullPage: true });
  });

  test("authenticated user visiting /register is redirected to /inbox", async ({ page }) => {
    await page.goto("/register");
    await expect(page).toHaveURL(/\/inbox/, { timeout: 8_000 });
    await page.screenshot({ path: "test-screenshots/37-auth-register-redirect.png", fullPage: true });
  });

  test("project detail page renders the saved project name", async ({ page, request }) => {
    const projectName = `Project ${Date.now()}`;
    const auth = await loginUser(request, user.email, user.password);
    const project = await createProject(request, auth.access_token, projectName);

    await page.goto(`/projects/${project.id}`);

    await expect(
      page.getByRole("heading", { name: projectName, exact: true }),
    ).toBeVisible();
  });
});
