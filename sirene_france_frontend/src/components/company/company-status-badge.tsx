import { Badge } from "@/components/ui/badge";

interface CompanyStatusBadgeProps {
  status?: string;
}

export function CompanyStatusBadge({ status }: CompanyStatusBadgeProps) {
  if (!status) return null;

  const isActive = status === "A";

  return (
    <Badge variant={isActive ? "default" : "destructive"}>
      {isActive ? "Active" : "Fermée"}
    </Badge>
  );
}
