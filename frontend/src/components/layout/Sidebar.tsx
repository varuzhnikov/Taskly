"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { Inbox, Hash, LogOut, Plus } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/store/auth";
import { logoutViaRouteHandler } from "@/lib/api/auth";
import { useProjects } from "@/hooks/useProjects";
import { Button } from "@/components/ui/button";
import { ProjectNavItem } from "@/components/layout/ProjectNavItem";

function NavItem({
  href,
  icon: Icon,
  label,
}: {
  href: string;
  icon: React.ElementType;
  label: string;
}) {
  const pathname = usePathname();
  const active = pathname === href || pathname.startsWith(href + "/");
  return (
    <Link
      href={href}
      className={cn(
        "flex items-center gap-2.5 rounded-md px-3 py-2 text-sm transition-colors",
        active
          ? "bg-white/10 text-white font-medium"
          : "text-gray-400 hover:bg-white/5 hover:text-white",
      )}
    >
      <Icon className="h-4 w-4 shrink-0" />
      {label}
    </Link>
  );
}

export function Sidebar() {
  const router = useRouter();
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const { data: projects } = useProjects();

  async function handleLogout() {
    await logoutViaRouteHandler();
    clearAuth();
    router.push("/login");
  }

  return (
    <aside className="flex h-full w-60 flex-col bg-[#1f2937] px-3 py-4 gap-1">
      <div className="px-3 py-2 mb-2">
        <span className="text-white font-bold text-lg tracking-tight">Todoist</span>
      </div>

      <NavItem href="/inbox" icon={Inbox} label="Inbox" />

      <div className="mt-4">
        <div className="flex items-center justify-between px-3 mb-1">
          <span className="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Projects
          </span>
          <Link href="/projects/new">
            <Plus className="h-3.5 w-3.5 text-gray-500 hover:text-white transition-colors" />
          </Link>
        </div>
        {projects?.map((p) => (
          <ProjectNavItem key={p.id} project={p} />
        ))}
      </div>

      <div className="mt-4">
        <div className="px-3 mb-1">
          <span className="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Labels
          </span>
        </div>
        <NavItem href="/labels" icon={Hash} label="All labels" />
      </div>

      <div className="mt-auto">
        <Button
          variant="ghost"
          size="sm"
          className="w-full justify-start text-gray-400 hover:text-white hover:bg-white/5"
          onClick={handleLogout}
        >
          <LogOut className="h-4 w-4 mr-2" />
          Sign out
        </Button>
      </div>
    </aside>
  );
}
