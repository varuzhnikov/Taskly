import { test, expect, type Page } from "@playwright/test";
import { createProject, loginUser, registerUser, loginViaUI } from "./helpers";

const user = {
  email: `e2e_projects_${Date.now()}@test.local`,
  password: "testpassword123",
};

test.describe("Delete project", () => {
  test.beforeAll(async ({ request }) => {
    await registerUser(request, user.email, user.password);
  });

  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, user.email, user.password);
  });

  // Returns the row div that wraps the project link + options button.
  // Scoping all sub-locators to this row prevents confusion when there are
  // multiple projects in the sidebar.
  function projectRow(page: Page, projectId: string) {
    return page.locator(`a[href="/projects/${projectId}"]`).locator("..");
  }

  async function openDeleteDialog(page: Page, projectId: string) {
    const row = projectRow(page, projectId);
    await row.hover();
    await row.getByRole("button", { name: /project options/i }).click();
    await page.getByRole("menuitem", { name: /delete project/i }).click();
    await expect(page.getByRole("dialog")).toBeVisible();
  }

  test("three dots button appears on project hover", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const project = await createProject(request, token.access_token, `Hover Test ${Date.now()}`);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${project.id}"]`);

    // screenshot before hover — dots are opacity-0
    await page.screenshot({ path: "test-screenshots/40-before-hover.png", fullPage: true });

    const row = projectRow(page, project.id);
    await row.hover();

    await page.screenshot({ path: "test-screenshots/40-project-hover.png", fullPage: true });

    // after hover, the button is in viewport and clickable
    await expect(row.getByRole("button", { name: /project options/i })).toBeInViewport();
  });

  test("dropdown opens on options button click", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const project = await createProject(request, token.access_token, `Dropdown Test ${Date.now()}`);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${project.id}"]`);

    const row = projectRow(page, project.id);
    await row.hover();
    await row.getByRole("button", { name: /project options/i }).click();

    await expect(page.getByRole("menuitem", { name: /delete project/i })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/41-project-dropdown.png", fullPage: true });
  });

  test("confirmation dialog appears on Delete project click", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const projectName = `Dialog Test ${Date.now()}`;
    const project = await createProject(request, token.access_token, projectName);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${project.id}"]`);

    await openDeleteDialog(page, project.id);

    const dialog = page.getByRole("dialog");
    await expect(dialog.getByRole("heading", { name: /delete/i })).toBeVisible();
    await expect(dialog.getByRole("heading", { name: new RegExp(projectName, "i") })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/42-delete-confirm-dialog.png", fullPage: true });
  });

  test("cancel dismisses dialog and project stays in sidebar", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const projectName = `Cancel Test ${Date.now()}`;
    const project = await createProject(request, token.access_token, projectName);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${project.id}"]`);

    await openDeleteDialog(page, project.id);

    await page.getByRole("dialog").getByRole("button", { name: /cancel/i }).click();

    await expect(page.getByRole("dialog")).not.toBeVisible();
    await expect(page.locator(`a[href="/projects/${project.id}"]`)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/43-delete-cancelled.png", fullPage: true });
  });

  test("confirming delete removes project from sidebar", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const projectName = `Delete Me ${Date.now()}`;
    const project = await createProject(request, token.access_token, projectName);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${project.id}"]`);

    await openDeleteDialog(page, project.id);

    await page.getByRole("dialog").getByRole("button", { name: /^delete$/i }).click();

    await expect(page.getByRole("dialog")).not.toBeVisible({ timeout: 8_000 });
    await expect(page.locator(`a[href="/projects/${project.id}"]`)).not.toBeVisible({ timeout: 8_000 });
    await page.screenshot({ path: "test-screenshots/44-project-deleted.png", fullPage: true });
  });

  test("deleting current project redirects to /inbox", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const projectName = `Active Delete ${Date.now()}`;
    const project = await createProject(request, token.access_token, projectName);

    await page.goto(`/projects/${project.id}`);
    await expect(page.getByRole("heading", { name: projectName })).toBeVisible();

    await openDeleteDialog(page, project.id);

    await page.getByRole("dialog").getByRole("button", { name: /^delete$/i }).click();

    await expect(page).toHaveURL(/\/inbox/, { timeout: 10_000 });
    await page.screenshot({ path: "test-screenshots/45-redirect-after-delete.png", fullPage: true });
  });

  test("deleting one project leaves other projects intact", async ({ page, request }) => {
    const token = await loginUser(request, user.email, user.password);
    const keepName = `Keep Me ${Date.now()}`;
    const deleteName = `Delete Me ${Date.now() + 1}`;
    const kept = await createProject(request, token.access_token, keepName);
    const toDelete = await createProject(request, token.access_token, deleteName);

    await page.goto("/inbox");
    await page.waitForSelector(`a[href="/projects/${kept.id}"]`);
    await page.waitForSelector(`a[href="/projects/${toDelete.id}"]`);

    await openDeleteDialog(page, toDelete.id);

    await page.getByRole("dialog").getByRole("button", { name: /^delete$/i }).click();

    await expect(page.locator(`a[href="/projects/${toDelete.id}"]`)).not.toBeVisible({ timeout: 8_000 });
    await expect(page.locator(`a[href="/projects/${kept.id}"]`)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/46-other-project-intact.png", fullPage: true });
  });
});
