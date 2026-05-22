import type { Metadata } from "next";

export const metadata: Metadata = { title: "Projects — Todoist Clone" };

export default function ProjectsPage() {
  return (
    <div>
      <h1 className="text-xl font-bold text-gray-900 mb-4">Projects</h1>
      <p className="text-sm text-gray-500">
        Select a project from the sidebar to view its tasks.
      </p>
    </div>
  );
}
