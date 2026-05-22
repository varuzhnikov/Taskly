"use client";

import { useState } from "react";
import { Flag, Calendar } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import { cn, formatDate, isOverdue } from "@/lib/utils";
import { useCompleteTask, useUncompleteTask } from "@/hooks/useTasks";
import type { Task } from "@/lib/types";
import { PriorityLabel } from "@/lib/types";

const priorityVariant = {
  0: "priority0",
  1: "priority1",
  2: "priority2",
  3: "priority3",
  4: "priority4",
} as const;

interface TaskItemProps {
  task: Task;
  onSelect?: (task: Task) => void;
}

export function TaskItem({ task, onSelect }: TaskItemProps) {
  const [optimisticDone, setOptimisticDone] = useState(!!task.completedAt);
  const complete = useCompleteTask();
  const uncomplete = useUncompleteTask();

  async function handleCheck(checked: boolean) {
    setOptimisticDone(checked);
    try {
      if (checked) {
        await complete.mutateAsync(task.id);
      } else {
        await uncomplete.mutateAsync(task.id);
      }
    } catch {
      setOptimisticDone(!checked);
    }
  }

  const overdue = isOverdue(task.dueAt) && !optimisticDone;

  return (
    <div
      className={cn(
        "group flex items-start gap-3 rounded-md px-2 py-2 hover:bg-gray-50 transition-colors",
        optimisticDone && "opacity-50",
      )}
    >
      <Checkbox
        checked={optimisticDone}
        onCheckedChange={handleCheck}
        aria-label={`Mark "${task.title}" as ${optimisticDone ? "incomplete" : "complete"}`}
        className="mt-0.5 shrink-0"
      />

      <button
        type="button"
        className="flex-1 text-left"
        onClick={() => onSelect?.(task)}
      >
        <p
          className={cn(
            "text-sm text-gray-900",
            optimisticDone && "line-through text-gray-400",
          )}
        >
          {task.title}
        </p>

        <div className="mt-1 flex items-center gap-2 flex-wrap">
          {task.priority > 0 && (
            <span className="flex items-center gap-1">
              <Flag className="h-3 w-3" />
              <Badge variant={priorityVariant[task.priority]} className="text-xs px-1.5 py-0">
                {PriorityLabel[task.priority]}
              </Badge>
            </span>
          )}

          {task.dueAt && (
            <span
              className={cn(
                "flex items-center gap-1 text-xs",
                overdue ? "text-red-500" : "text-gray-400",
              )}
            >
              <Calendar className="h-3 w-3" />
              {formatDate(task.dueAt)}
            </span>
          )}
        </div>
      </button>
    </div>
  );
}
