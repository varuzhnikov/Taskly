"use client";

import { useAuthStore } from "@/store/auth";
import { User } from "lucide-react";

export function Header() {
  const user = useAuthStore((s) => s.user);

  return (
    <header className="flex h-12 items-center justify-end border-b border-gray-200 bg-white px-6">
      {user && (
        <div className="flex items-center gap-2 text-sm text-gray-600">
          <User className="h-4 w-4" />
          <span>{user.email}</span>
        </div>
      )}
    </header>
  );
}
