import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors",
  {
    variants: {
      variant: {
        default: "border-transparent bg-[#db4035] text-white",
        secondary: "border-transparent bg-gray-100 text-gray-900",
        outline: "text-gray-700 border-gray-300",
        priority0: "border-transparent bg-gray-100 text-gray-500",
        priority1: "border-transparent bg-blue-100 text-blue-700",
        priority2: "border-transparent bg-yellow-100 text-yellow-700",
        priority3: "border-transparent bg-orange-100 text-orange-700",
        priority4: "border-transparent bg-red-100 text-red-700",
      },
    },
    defaultVariants: { variant: "default" },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <span className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
