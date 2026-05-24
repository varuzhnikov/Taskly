import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ProjectNavItem } from "@/components/layout/ProjectNavItem";
import type { Project } from "@/lib/types";

// --- mocks ---

const mockMutateAsync = vi.fn().mockResolvedValue(undefined);
const mockPush = vi.fn();
const mockPathname = vi.fn(() => "/inbox");

vi.mock("@/hooks/useProjects", () => ({
  useDeleteProject: () => ({
    mutateAsync: mockMutateAsync,
    isPending: false,
  }),
}));

vi.mock("next/navigation", () => ({
  usePathname: () => mockPathname(),
  useRouter: () => ({ push: mockPush }),
}));

vi.mock("next/link", () => ({
  default: ({ href, children, className }: { href: string; children: React.ReactNode; className?: string }) => (
    <a href={href} className={className}>{children}</a>
  ),
}));

// --- helpers ---

function wrapper({ children }: { children: React.ReactNode }) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

const project: Project = {
  id: "proj-abc-123",
  name: "My Project",
  color: "#db4035",
  sortOrder: 0,
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

// --- tests ---

describe("ProjectNavItem", () => {
  beforeEach(() => {
    mockMutateAsync.mockClear();
    mockPush.mockClear();
    mockPathname.mockReturnValue("/inbox");
  });

  it("renders the project name", () => {
    render(<ProjectNavItem project={project} />, { wrapper });
    expect(screen.getByText("My Project")).toBeInTheDocument();
  });

  it("renders a link pointing to the project route", () => {
    render(<ProjectNavItem project={project} />, { wrapper });
    const link = screen.getByRole("link", { name: /my project/i });
    expect(link).toHaveAttribute("href", "/projects/proj-abc-123");
  });

  it("renders the options button in the DOM", () => {
    render(<ProjectNavItem project={project} />, { wrapper });
    expect(screen.getByRole("button", { name: /project options/i })).toBeInTheDocument();
  });

  it("options button is initially visually hidden via opacity-0", () => {
    render(<ProjectNavItem project={project} />, { wrapper });
    const btn = screen.getByRole("button", { name: /project options/i });
    expect(btn.className).toMatch(/opacity-0/);
  });

  it("opens dropdown with Delete project item on options button click", async () => {
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));

    await waitFor(() =>
      expect(screen.getByRole("menuitem", { name: /delete project/i })).toBeInTheDocument(),
    );
  });

  it("opens confirmation dialog when Delete project is clicked", async () => {
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));
    await waitFor(() => screen.getByRole("menuitem", { name: /delete project/i }));
    await user.click(screen.getByRole("menuitem", { name: /delete project/i }));

    await waitFor(() =>
      expect(screen.getByRole("dialog")).toBeInTheDocument(),
    );
    expect(screen.getByRole("heading", { name: /delete.*my project/i })).toBeInTheDocument();
  });

  it("cancel closes the dialog without calling delete", async () => {
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));
    await waitFor(() => screen.getByRole("menuitem", { name: /delete project/i }));
    await user.click(screen.getByRole("menuitem", { name: /delete project/i }));
    await waitFor(() => screen.getByRole("dialog"));

    await user.click(screen.getByRole("button", { name: /cancel/i }));

    await waitFor(() =>
      expect(screen.queryByRole("dialog")).not.toBeInTheDocument(),
    );
    expect(mockMutateAsync).not.toHaveBeenCalled();
  });

  it("confirm calls mutateAsync with the project id", async () => {
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));
    await waitFor(() => screen.getByRole("menuitem", { name: /delete project/i }));
    await user.click(screen.getByRole("menuitem", { name: /delete project/i }));
    await waitFor(() => screen.getByRole("dialog"));

    await user.click(screen.getByRole("button", { name: /^delete$/i }));

    await waitFor(() =>
      expect(mockMutateAsync).toHaveBeenCalledWith("proj-abc-123"),
    );
  });

  it("redirects to /inbox when deleting the currently active project", async () => {
    mockPathname.mockReturnValue("/projects/proj-abc-123");
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));
    await waitFor(() => screen.getByRole("menuitem", { name: /delete project/i }));
    await user.click(screen.getByRole("menuitem", { name: /delete project/i }));
    await waitFor(() => screen.getByRole("dialog"));

    await user.click(screen.getByRole("button", { name: /^delete$/i }));

    await waitFor(() => expect(mockPush).toHaveBeenCalledWith("/inbox"));
  });

  it("does not redirect when deleting a project that is not currently active", async () => {
    mockPathname.mockReturnValue("/projects/other-project");
    const user = userEvent.setup();
    render(<ProjectNavItem project={project} />, { wrapper });

    await user.click(screen.getByRole("button", { name: /project options/i }));
    await waitFor(() => screen.getByRole("menuitem", { name: /delete project/i }));
    await user.click(screen.getByRole("menuitem", { name: /delete project/i }));
    await waitFor(() => screen.getByRole("dialog"));

    await user.click(screen.getByRole("button", { name: /^delete$/i }));

    await waitFor(() => expect(mockMutateAsync).toHaveBeenCalled());
    expect(mockPush).not.toHaveBeenCalled();
  });
});
