import { z } from "zod";

export const createProjectSchema = z.object({
  name: z.string().min(1, "Name required").max(200),
  color: z.string().regex(/^#[0-9a-fA-F]{6}$/, "Must be a hex color").optional(),
});

export type CreateProjectFormValues = z.infer<typeof createProjectSchema>;
