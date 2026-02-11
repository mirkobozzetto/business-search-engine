import type { Meta } from "@/types/api";

interface ResultsMetaProps {
  meta?: Meta;
}

export function ResultsMeta({ meta }: ResultsMetaProps) {
  if (!meta) return null;

  return (
    <div className="flex items-center gap-4 text-sm text-muted-foreground">
      {meta.total !== undefined && (
        <span>{meta.total.toLocaleString("fr-FR")} résultats</span>
      )}
      {meta.duration_ms !== undefined && (
        <span>{meta.duration_ms} ms</span>
      )}
      {meta.page !== undefined && meta.pages !== undefined && (
        <span>
          Page {meta.page} / {meta.pages}
        </span>
      )}
    </div>
  );
}
