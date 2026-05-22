import { api } from "./client";
import type { Task, TaskLink, Priority } from "@/lib/types";

export interface CreateTaskRequest {
  title: string;
  descriptionMd?: string;
  links?: TaskLink[];
  priority?: Priority;
  dueAt?: string | null;
  projectId?: string | null;
  parentTaskId?: string | null;
  labelIds?: string[];
}

export interface UpdateTaskRequest {
  title: string;
  descriptionMd?: string;
  links?: TaskLink[];
  priority?: Priority;
  dueAt?: string | null;
  projectId?: string | null;
  parentTaskId?: string | null;
  labelIds?: string[];
}

export interface TaskListParams {
  inbox?: boolean;
  projectId?: string;
  labelId?: string;
  priority?: Priority;
  completed?: boolean;
  search?: string;
  limit?: number;
  offset?: number;
}

interface TaskBackendPayload {
  title: string;
  description?: string;
  links?: TaskLink[];
  priority?: Priority;
  due_at?: string | null;
  project_id?: string | null;
  parent_task_id?: string | null;
  label_ids?: string[];
}

function buildQuery(params: TaskListParams): string {
  const q = new URLSearchParams();
  if (params.inbox) q.set("inbox", "true");
  if (params.projectId) q.set("project_id", params.projectId);
  if (params.labelId) q.set("label_id", params.labelId);
  if (params.priority !== undefined) q.set("priority", String(params.priority));
  if (params.completed !== undefined)
    q.set("completed", String(params.completed));
  if (params.search) q.set("search", params.search);
  if (params.limit !== undefined) q.set("limit", String(params.limit));
  if (params.offset !== undefined) q.set("offset", String(params.offset));
  const s = q.toString();
  return s ? `?${s}` : "";
}

function toTaskBackendPayload(
  data: CreateTaskRequest | UpdateTaskRequest,
): TaskBackendPayload {
  return {
    title: data.title,
    description: data.descriptionMd,
    links: data.links,
    priority: data.priority,
    due_at: data.dueAt,
    project_id: data.projectId,
    parent_task_id: data.parentTaskId,
    label_ids: data.labelIds,
  };
}

export const tasksApi = {
  list: (params: TaskListParams = {}) =>
    api.get<Task[]>(`/api/v1/tasks/${buildQuery(params)}`),
  getById: (id: string) => api.get<Task>(`/api/v1/tasks/${id}`),
  create: (data: CreateTaskRequest) =>
    api.post<Task>("/api/v1/tasks/", toTaskBackendPayload(data)),
  update: (id: string, data: UpdateTaskRequest) =>
    api.put<Task>(`/api/v1/tasks/${id}`, toTaskBackendPayload(data)),
  delete: (id: string) => api.delete<void>(`/api/v1/tasks/${id}`),
  complete: (id: string) => api.patch<Task>(`/api/v1/tasks/${id}/complete`),
  uncomplete: (id: string) => api.patch<Task>(`/api/v1/tasks/${id}/uncomplete`),
};
