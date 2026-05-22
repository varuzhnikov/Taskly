"use client";

import { TaskList } from "@/components/tasks/TaskList";
import { useProject } from "@/hooks/useProjects";

interface ProjectTaskListProps {
  projectId: string;
}

export function ProjectTaskList({ projectId }: ProjectTaskListProps) {
  const { data: project, isLoading, isError } = useProject(projectId);

  if (isLoading) {
    return <div className="p-4 text-sm text-gray-400 animate-pulse">Loading project…</div>;
  }

  if (isError || !project) {
    return <div className="p-4 text-sm text-red-500">Failed to load project.</div>;
  }

  return <TaskList title={project.name} params={{ projectId }} />;
}
