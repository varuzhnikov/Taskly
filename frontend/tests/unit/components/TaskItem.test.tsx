import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TaskItem } from "@/components/tasks/TaskItem";
import type { Task } from "@/lib/types";

vi.mock("@/hooks/useTasks", () => ({
  useCompleteTask: () => ({
    mutateAsync: vi.fn().mockResolvedValue(undefined),
  }),
  useUncompleteTask: () => ({
    mutateAsync: vi.fn().mockResolvedValue(undefined),
  }),
}));

function wrapper({ children }: { children: React.ReactNode }) {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

const baseTask: Task = {
  id: "task-1",
  projectId: null,
  parentTaskId: null,
  title: "Write unit tests",
  descriptionMd: "",
  descriptionHtml: "",
  links: [],
  priority: 0,
  dueAt: null,
  completedAt: null,
  sortOrder: 0,
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

describe("TaskItem", () => {
  it("renders the task title", () => {
    render(<TaskItem task={baseTask} />, { wrapper });
    expect(screen.getByText("Write unit tests")).toBeInTheDocument();
  });

  it("renders unchecked checkbox for incomplete task", () => {
    render(<TaskItem task={baseTask} />, { wrapper });
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).not.toBeChecked();
  });

  it("renders checked checkbox for completed task", () => {
    const completed = { ...baseTask, completedAt: "2026-01-02T00:00:00Z" };
    render(<TaskItem task={completed} />, { wrapper });
    const checkbox = screen.getByRole("checkbox");
    expect(checkbox).toBeChecked();
  });

  it("shows priority badge when priority > 0", () => {
    const urgent = { ...baseTask, priority: 4 as const };
    render(<TaskItem task={urgent} />, { wrapper });
    expect(screen.getByText("Urgent")).toBeInTheDocument();
  });

  it("does not show priority badge for priority 0", () => {
    render(<TaskItem task={baseTask} />, { wrapper });
    expect(screen.queryByText("None")).not.toBeInTheDocument();
  });

  it("shows formatted due date when dueAt is set", () => {
    const withDue = { ...baseTask, dueAt: "2026-06-15T00:00:00Z" };
    render(<TaskItem task={withDue} />, { wrapper });
    expect(screen.getByText(/jun/i)).toBeInTheDocument();
  });

  it("calls onSelect when title button is clicked", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(<TaskItem task={baseTask} onSelect={onSelect} />, { wrapper });
    await user.click(screen.getByRole("button", { name: /write unit tests/i }));
    expect(onSelect).toHaveBeenCalledWith(baseTask);
  });

  it("optimistically toggles checkbox on check", async () => {
    const user = userEvent.setup();
    render(<TaskItem task={baseTask} />, { wrapper });
    const checkbox = screen.getByRole("checkbox");
    await user.click(checkbox);
    expect(checkbox).toBeChecked();
  });
});
