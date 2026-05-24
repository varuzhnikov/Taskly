"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { FolderOpen, MoreHorizontal, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { useDeleteProject } from "@/hooks/useProjects";
import type { Project } from "@/lib/types";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogClose,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

export function ProjectNavItem({ project }: { project: Project }) {
  const pathname = usePathname();
  const router = useRouter();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const deleteProject = useDeleteProject();

  const href = `/projects/${project.id}`;
  const active = pathname === href || pathname.startsWith(href + "/");

  async function handleDelete() {
    await deleteProject.mutateAsync(project.id);
    setConfirmOpen(false);
    if (pathname.startsWith(`/projects/${project.id}`)) {
      router.push("/inbox");
    }
  }

  return (
    <>
      <div className="group relative flex items-center rounded-md transition-colors hover:bg-white/5">
        <Link
          href={href}
          className={cn(
            "flex flex-1 items-center gap-2.5 px-3 py-2 text-sm transition-colors",
            active ? "text-white font-medium" : "text-gray-400 hover:text-white",
          )}
        >
          <FolderOpen className="h-4 w-4 shrink-0" />
          {project.name}
        </Link>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              aria-label="Project options"
              className="mr-2 rounded p-0.5 opacity-0 transition-opacity group-hover:opacity-100 data-[state=open]:opacity-100 focus:opacity-100 focus:outline-none text-gray-400 hover:text-white"
            >
              <MoreHorizontal className="h-3.5 w-3.5" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="start" side="right">
            <DropdownMenuItem
              className="text-red-500 focus:text-red-500 focus:bg-red-50"
              onSelect={() => setConfirmOpen(true)}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete project
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete &ldquo;{project.name}&rdquo;?</DialogTitle>
          </DialogHeader>
          <p className="mt-2 text-sm text-gray-500">
            All tasks inside this project will also be permanently deleted.
          </p>
          <div className="mt-4 flex justify-end gap-2">
            <DialogClose asChild>
              <Button variant="outline" size="sm">
                Cancel
              </Button>
            </DialogClose>
            <Button
              variant="destructive"
              size="sm"
              onClick={handleDelete}
              disabled={deleteProject.isPending}
            >
              {deleteProject.isPending ? "Deleting…" : "Delete"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
