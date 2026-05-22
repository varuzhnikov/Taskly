import { api } from "./client";
import type { Project } from "@/lib/types";

export interface CreateProjectRequest {
  name: string;
  color?: string;
}

export interface UpdateProjectRequest {
  name: string;
  color?: string;
  sortOrder?: number;
}

export const projectsApi = {
  list: () => api.get<Project[]>("/api/v1/projects/"),
  getById: (id: string) => api.get<Project>(`/api/v1/projects/${id}`),
  create: (data: CreateProjectRequest) =>
    api.post<Project>("/api/v1/projects/", data),
  update: (id: string, data: UpdateProjectRequest) =>
    api.put<Project>(`/api/v1/projects/${id}`, data),
  delete: (id: string) => api.delete<void>(`/api/v1/projects/${id}`),
};
