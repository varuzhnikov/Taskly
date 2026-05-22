import { z } from "zod";

export const taskLinkSchema = z.object({
  url: z.string().url("Must be a valid URL"),
  title: z.string().min(1, "Title required"),
});

export const createTaskSchema = z.object({
  title: z.string().min(1, "Title required").max(500),
  descriptionMd: z.string().max(10000).optional(),
  links: z.array(taskLinkSchema).max(10).optional(),
  priority: z.union([
    z.literal(0),
    z.literal(1),
    z.literal(2),
    z.literal(3),
    z.literal(4),
  ]).optional(),
  dueAt: z.string().datetime({ offset: true }).nullable().optional(),
  projectId: z.string().uuid().nullable().optional(),
  labelIds: z.array(z.string().uuid()).optional(),
});

export type CreateTaskFormValues = z.infer<typeof createTaskSchema>;
