import { api } from "./client";
import type { Comment } from "@/lib/types";

export interface CreateCommentRequest {
  body: string;
}

export interface UpdateCommentRequest {
  body: string;
}

export const commentsApi = {
  listByTask: (taskId: string) =>
    api.get<Comment[]>(`/api/v1/tasks/${taskId}/comments/`),
  getById: (taskId: string, commentId: string) =>
    api.get<Comment>(`/api/v1/tasks/${taskId}/comments/${commentId}`),
  create: (taskId: string, data: CreateCommentRequest) =>
    api.post<Comment>(`/api/v1/tasks/${taskId}/comments/`, data),
  update: (taskId: string, commentId: string, data: UpdateCommentRequest) =>
    api.put<Comment>(`/api/v1/tasks/${taskId}/comments/${commentId}`, data),
  delete: (taskId: string, commentId: string) =>
    api.delete<void>(`/api/v1/tasks/${taskId}/comments/${commentId}`),
};
