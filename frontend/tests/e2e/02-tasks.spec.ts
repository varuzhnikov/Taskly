import { test, expect } from "@playwright/test";
import { createProject, loginUser, registerUser, loginViaUI } from "./helpers";

const user = {
  email: `e2e_tasks_${Date.now()}@test.local`,
  password: "testpassword123",
};

test.describe("Tasks", () => {
  test.beforeAll(async ({ request }) => {
    await registerUser(request, user.email, user.password);
  });

  test.beforeEach(async ({ page }) => {
    await loginViaUI(page, user.email, user.password);
  });

  test("inbox shows empty state with Add task button", async ({ page }) => {
    await expect(page.getByRole("button", { name: /add task/i })).toBeVisible();
    await page.screenshot({ path: "test-screenshots/20-inbox-empty.png", fullPage: true });
  });

  test("clicking Add task opens inline form", async ({ page }) => {
    await page.getByRole("button", { name: /add task/i }).click();
    await expect(page.getByPlaceholder("Task name")).toBeVisible();
    await expect(page.getByPlaceholder(/description/i)).toBeVisible();
    await page.screenshot({ path: "test-screenshots/21-add-task-form-open.png", fullPage: true });
  });

  test("Cancel closes the add task form", async ({ page }) => {
    await page.getByRole("button", { name: /add task/i }).click();
    await expect(page.getByPlaceholder("Task name")).toBeVisible();
    await page.getByRole("button", { name: /cancel/i }).click();
    await expect(page.getByPlaceholder("Task name")).not.toBeVisible();
    await page.screenshot({ path: "test-screenshots/22-add-task-cancelled.png", fullPage: true });
  });

  test("can create a task and it appears in the list", async ({ page }) => {
    const title = `Test task ${Date.now()}`;
    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByPlaceholder("Task name").fill(title);
    await page.screenshot({ path: "test-screenshots/23-task-title-filled.png", fullPage: true });
    await page.getByRole("button", { name: /^add task$/i }).click();
    await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });
    await page.screenshot({ path: "test-screenshots/24-task-created.png", fullPage: true });
  });

  test("can create multiple tasks", async ({ page }) => {
    const titles = [`Task A ${Date.now()}`, `Task B ${Date.now()}`];
    for (const title of titles) {
      await page.getByRole("button", { name: /add task/i }).click();
      await page.getByPlaceholder("Task name").fill(title);
      await page.getByRole("button", { name: /^add task$/i }).click();
      await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });
    }
    await page.screenshot({ path: "test-screenshots/25-multiple-tasks.png", fullPage: true });
  });

  test("can complete a task", async ({ page }) => {
    const title = `Complete me ${Date.now()}`;

    // Create the task first
    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByPlaceholder("Task name").fill(title);
    await page.getByRole("button", { name: /^add task$/i }).click();
    await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });

    // Click the exact checkbox for this task via its aria-label
    await page
      .getByRole("checkbox", { name: `Mark "${title}" as complete` })
      .click();

    await page.screenshot({ path: "test-screenshots/26-task-completed.png", fullPage: true });
    await expect(page.getByText(/\d+ completed/i)).toBeVisible({ timeout: 8_000 });
  });

  test("completed tasks are listed in the collapsed section", async ({ page }) => {
    const title = `Collapse me ${Date.now()}`;

    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByPlaceholder("Task name").fill(title);
    await page.getByRole("button", { name: /^add task$/i }).click();
    await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });

    await page
      .getByRole("checkbox", { name: `Mark "${title}" as complete` })
      .click();

    const summary = page.getByText(/\d+ completed/i);
    await expect(summary).toBeVisible({ timeout: 8_000 });

    // Expand the completed section
    await summary.click();
    await page.screenshot({ path: "test-screenshots/27-completed-expanded.png", fullPage: true });
    await expect(page.getByText(title)).toBeVisible();
  });

  test("empty task title is not submitted", async ({ page }) => {
    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByRole("button", { name: /^add task$/i }).click();
    await expect(page.getByPlaceholder("Task name")).toBeVisible();
    await page.screenshot({ path: "test-screenshots/28-empty-task-validation.png", fullPage: true });
  });

  test("task created on a project page stays in that project and not inbox", async ({ page, request }) => {
    const auth = await loginUser(request, user.email, user.password);
    const projectName = `Tasks Project ${Date.now()}`;
    const project = await createProject(request, auth.access_token, projectName);
    const title = `Project task ${Date.now()}`;

    await page.goto(`/projects/${project.id}`);
    await expect(
      page.getByRole("heading", { name: projectName, exact: true }),
    ).toBeVisible();

    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByPlaceholder("Task name").fill(title);
    await page.getByRole("button", { name: /^add task$/i }).click();

    await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });

    await page.goto("/inbox");
    await expect(page.getByText(title)).not.toBeVisible();
  });

  test("task created in inbox stays visible in inbox", async ({ page }) => {
    const title = `Inbox task ${Date.now()}`;

    await page.goto("/inbox");
    await page.getByRole("button", { name: /add task/i }).click();
    await page.getByPlaceholder("Task name").fill(title);
    await page.getByRole("button", { name: /^add task$/i }).click();

    await expect(page.getByText(title)).toBeVisible({ timeout: 8_000 });
  });
});
