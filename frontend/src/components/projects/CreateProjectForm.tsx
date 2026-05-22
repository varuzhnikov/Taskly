"use client";

import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { createProjectSchema, type CreateProjectFormValues } from "@/schemas/project";
import { useCreateProject } from "@/hooks/useProjects";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

export function CreateProjectForm() {
  const router = useRouter();
  const createProject = useCreateProject();

  const form = useForm<CreateProjectFormValues>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: "",
      color: "#db4035",
    },
  });

  const {
    handleSubmit,
    formState: { isSubmitting, errors },
    setError,
  } = form;

  async function onSubmit(values: CreateProjectFormValues) {
    try {
      const project = await createProject.mutateAsync(values);
      router.push(`/projects/${project.id}`);
    } catch (err) {
      setError("root", {
        message: err instanceof Error ? err.message : "Project creation failed",
      });
    }
  }

  return (
    <FormProvider {...form}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Project name" autoFocus {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="color"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Color</FormLabel>
              <FormControl>
                <Input type="color" className="h-11 w-20 p-1" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {errors.root && <p className="text-sm text-red-600">{errors.root.message}</p>}

        <div className="flex items-center gap-2">
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Creating…" : "Create project"}
          </Button>
          <Button type="button" variant="ghost" onClick={() => router.push("/projects")}>
            Cancel
          </Button>
        </div>
      </form>
    </FormProvider>
  );
}
