import { api } from "./client";
import type { Label } from "@/lib/types";

export interface CreateLabelRequest {
  name: string;
  color?: string;
}

export interface UpdateLabelRequest {
  name: string;
  color?: string;
}

export const labelsApi = {
  list: () => api.get<Label[]>("/api/v1/labels/"),
  getById: (id: string) => api.get<Label>(`/api/v1/labels/${id}`),
  create: (data: CreateLabelRequest) =>
    api.post<Label>("/api/v1/labels/", data),
  update: (id: string, data: UpdateLabelRequest) =>
    api.put<Label>(`/api/v1/labels/${id}`, data),
  delete: (id: string) => api.delete<void>(`/api/v1/labels/${id}`),
};
