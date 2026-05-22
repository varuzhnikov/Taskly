"use client";

import { useQuery } from "@tanstack/react-query";
import { labelsApi } from "@/lib/api/labels";
import { Badge } from "@/components/ui/badge";

export default function LabelsPage() {
  const { data: labels, isLoading } = useQuery({
    queryKey: ["labels"],
    queryFn: labelsApi.list,
  });

  return (
    <div>
      <h1 className="text-xl font-bold text-gray-900 mb-4">Labels</h1>
      {isLoading && <p className="text-sm text-gray-400">Loading…</p>}
      <div className="flex flex-wrap gap-2">
        {labels?.map((label) => (
          <Badge
            key={label.id}
            variant="outline"
            style={{ borderColor: label.color, color: label.color }}
          >
            {label.name}
          </Badge>
        ))}
      </div>
    </div>
  );
}
