"use client";

import { useState } from "react";
import { useTasks } from "@/hooks/useTasks";
import { TaskItem } from "./TaskItem";
import { AddTaskForm } from "./AddTaskForm";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import type { Task } from "@/lib/types";
import type { TaskListParams } from "@/lib/api/tasks";

interface TaskListProps {
  params?: TaskListParams;
  title: string;
  onSelectTask?: (task: Task) => void;
}

export function TaskList({ params, title, onSelectTask }: TaskListProps) {
  const [showAddForm, setShowAddForm] = useState(false);
  const { data: tasks, isLoading, isError } = useTasks(params ?? {});

  if (isLoading) {
    return (
      <div className="p-4 text-sm text-gray-400 animate-pulse">Loading…</div>
    );
  }

  if (isError) {
    return (
      <div className="p-4 text-sm text-red-500">Failed to load tasks.</div>
    );
  }

  const active = tasks?.filter((t) => !t.completedAt) ?? [];
  const done = tasks?.filter((t) => t.completedAt) ?? [];

  return (
    <div className="flex flex-col gap-1">
      <h1 className="text-xl font-bold text-gray-900 mb-4">{title}</h1>

      {active.map((task) => (
        <TaskItem key={task.id} task={task} onSelect={onSelectTask} />
      ))}

      {showAddForm ? (
        <AddTaskForm
          defaultProjectId={params?.projectId}
          onCancel={() => setShowAddForm(false)}
          onSuccess={() => setShowAddForm(false)}
        />
      ) : (
        <Button
          variant="ghost"
          size="sm"
          className="justify-start text-gray-400 hover:text-[#db4035] mt-1"
          onClick={() => setShowAddForm(true)}
        >
          <Plus className="h-4 w-4 mr-1" />
          Add task
        </Button>
      )}

      {done.length > 0 && (
        <details className="mt-6 group">
          <summary className="cursor-pointer text-sm text-gray-400 hover:text-gray-600 select-none list-none flex items-center gap-1">
            <span className="transition-transform group-open:rotate-90">▶</span>
            {done.length} completed
          </summary>
          <div className="mt-2 flex flex-col gap-1">
            {done.map((task) => (
              <TaskItem key={task.id} task={task} onSelect={onSelectTask} />
            ))}
          </div>
        </details>
      )}
    </div>
  );
}
