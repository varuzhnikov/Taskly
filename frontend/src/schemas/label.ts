import { z } from "zod";

export const createLabelSchema = z.object({
  name: z.string().min(1, "Name required").max(100),
  color: z.string().regex(/^#[0-9a-fA-F]{6}$/, "Must be a hex color").optional(),
});

export type CreateLabelFormValues = z.infer<typeof createLabelSchema>;
