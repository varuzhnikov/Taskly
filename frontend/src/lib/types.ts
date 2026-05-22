export interface SessionUser {
  id: string;
  email: string;
}

export interface Project {
  id: string;
  name: string;
  color: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface Label {
  id: string;
  name: string;
  color: string;
  createdAt: string;
  updatedAt: string;
}

export type Priority = 0 | 1 | 2 | 3 | 4;

export const PriorityLabel: Record<Priority, string> = {
  0: "None",
  1: "Low",
  2: "Medium",
  3: "High",
  4: "Urgent",
};

export const PriorityColor: Record<Priority, string> = {
  0: "text-gray-400",
  1: "text-blue-500",
  2: "text-yellow-500",
  3: "text-orange-500",
  4: "text-red-500",
};

export interface TaskLink {
  url: string;
  title: string;
}

export interface Task {
  id: string;
  projectId: string | null;
  parentTaskId: string | null;
  title: string;
  descriptionMd: string;
  descriptionHtml: string;
  links: TaskLink[];
  priority: Priority;
  dueAt: string | null;
  completedAt: string | null;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface Comment {
  id: string;
  taskId: string;
  userId: string;
  body: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthTokens {
  accessToken: string;
  user: SessionUser;
}

export interface ApiError {
  error: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}
