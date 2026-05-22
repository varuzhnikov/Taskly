import type { Metadata } from "next";
import { ProjectTaskList } from "@/components/projects/ProjectTaskList";

interface Props {
  params: Promise<{ id: string }>;
}

export const metadata: Metadata = { title: "Project — Todoist Clone" };

export default async function ProjectPage({ params }: Props) {
  const { id } = await params;
  return <ProjectTaskList projectId={id} />;
}
