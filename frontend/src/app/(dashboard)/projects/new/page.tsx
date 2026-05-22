import type { Metadata } from "next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CreateProjectForm } from "@/components/projects/CreateProjectForm";

export const metadata: Metadata = { title: "New Project — Todoist Clone" };

export default function NewProjectPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Create project</CardTitle>
      </CardHeader>
      <CardContent>
        <CreateProjectForm />
      </CardContent>
    </Card>
  );
}
