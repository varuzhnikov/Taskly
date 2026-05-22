"use client";

import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { createTaskSchema, type CreateTaskFormValues } from "@/schemas/task";
import { useCreateTask } from "@/hooks/useTasks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from "@/components/ui/form";

interface AddTaskFormProps {
  defaultProjectId?: string;
  onCancel: () => void;
  onSuccess: () => void;
}

export function AddTaskForm({
  defaultProjectId,
  onCancel,
  onSuccess,
}: AddTaskFormProps) {
  const createTask = useCreateTask();

  const form = useForm<CreateTaskFormValues>({
    resolver: zodResolver(createTaskSchema),
    defaultValues: {
      title: "",
      descriptionMd: "",
      priority: 0,
      projectId: defaultProjectId ?? null,
    },
  });

  const { handleSubmit, formState: { isSubmitting } } = form;

  async function onSubmit(values: CreateTaskFormValues) {
    await createTask.mutateAsync(values);
    form.reset();
    onSuccess();
  }

  return (
    <FormProvider {...form}>
      <form
        onSubmit={handleSubmit(onSubmit)}
        className="rounded-md border border-gray-200 p-3 space-y-3 mt-1"
      >
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  placeholder="Task name"
                  autoFocus
                  className="border-0 shadow-none focus-visible:ring-0 p-0 text-sm font-medium"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="descriptionMd"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Textarea
                  placeholder="Description (markdown supported)"
                  className="border-0 shadow-none focus-visible:ring-0 p-0 text-sm text-gray-500 min-h-[40px] resize-none"
                  {...field}
                />
              </FormControl>
            </FormItem>
          )}
        />

        <div className="flex items-center justify-end gap-2 pt-1 border-t border-gray-100">
          <Button type="button" variant="ghost" size="sm" onClick={onCancel}>
            Cancel
          </Button>
          <Button type="submit" size="sm" disabled={isSubmitting}>
            {isSubmitting ? "Adding…" : "Add task"}
          </Button>
        </div>
      </form>
    </FormProvider>
  );
}
