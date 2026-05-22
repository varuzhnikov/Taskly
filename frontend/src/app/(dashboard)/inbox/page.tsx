import type { Metadata } from "next";
import { TaskList } from "@/components/tasks/TaskList";

export const metadata: Metadata = { title: "Inbox — Todoist Clone" };

export default function InboxPage() {
  return <TaskList title="Inbox" params={{ inbox: true }} />;
}
