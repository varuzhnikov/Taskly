"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Sidebar } from "@/components/layout/Sidebar";
import { Header } from "@/components/layout/Header";
import { useAuthStore } from "@/store/auth";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { accessToken, setAuth, clearAuth } = useAuthStore();
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    if (accessToken) {
      setHydrated(true);
      return;
    }

    // Page was refreshed — Zustand is empty but the httpOnly refresh cookie
    // may still be valid. Try a silent token refresh before rendering.
    fetch("/api/auth/refresh", { method: "POST" })
      .then(async (res) => {
        if (!res.ok) {
          clearAuth();
          router.replace("/login");
          return;
        }
        const data = (await res.json()) as {
          accessToken: string;
          user: { id: string; email: string };
        };
        setAuth(data.accessToken, data.user);
        setHydrated(true);
      })
      .catch(() => {
        clearAuth();
        router.replace("/login");
      });
  }, []);  // eslint-disable-line react-hooks/exhaustive-deps

  if (!hydrated) {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-[#db4035] border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header />
        <main className="flex-1 overflow-y-auto p-8 max-w-3xl mx-auto w-full">
          {children}
        </main>
      </div>
    </div>
  );
}
