import { describe, it, expect } from "vitest";
import { loginSchema, registerSchema } from "@/schemas/auth";
import { createTaskSchema } from "@/schemas/task";

describe("loginSchema", () => {
  it("accepts valid credentials", () => {
    expect(loginSchema.safeParse({ email: "a@b.com", password: "12345678" }).success).toBe(true);
  });

  it("rejects invalid email", () => {
    const r = loginSchema.safeParse({ email: "notanemail", password: "12345678" });
    expect(r.success).toBe(false);
  });

  it("rejects short password", () => {
    const r = loginSchema.safeParse({ email: "a@b.com", password: "short" });
    expect(r.success).toBe(false);
  });
});

describe("registerSchema", () => {
  it("accepts matching passwords", () => {
    expect(
      registerSchema.safeParse({
        email: "a@b.com",
        password: "securepassword",
        confirmPassword: "securepassword",
      }).success,
    ).toBe(true);
  });

  it("rejects mismatched passwords", () => {
    const r = registerSchema.safeParse({
      email: "a@b.com",
      password: "securepassword",
      confirmPassword: "different",
    });
    expect(r.success).toBe(false);
  });
});

describe("createTaskSchema", () => {
  it("accepts a minimal task", () => {
    expect(createTaskSchema.safeParse({ title: "Buy groceries" }).success).toBe(true);
  });

  it("rejects empty title", () => {
    expect(createTaskSchema.safeParse({ title: "" }).success).toBe(false);
  });

  it("rejects invalid priority", () => {
    expect(createTaskSchema.safeParse({ title: "Task", priority: 5 }).success).toBe(false);
  });

  it("accepts all valid priorities", () => {
    for (const p of [0, 1, 2, 3, 4] as const) {
      expect(
        createTaskSchema.safeParse({ title: "Task", priority: p }).success,
      ).toBe(true);
    }
  });

  it("rejects invalid link URL", () => {
    const r = createTaskSchema.safeParse({
      title: "Task",
      links: [{ url: "not-a-url", title: "Bad link" }],
    });
    expect(r.success).toBe(false);
  });
});
