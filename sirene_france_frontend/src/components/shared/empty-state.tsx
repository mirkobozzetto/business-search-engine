interface EmptyStateProps {
  title: string;
  description?: string;
}

export function EmptyState({ title, description }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <h3 className="text-lg font-medium text-muted-foreground">{title}</h3>
      {description && (
        <p className="mt-2 text-sm text-muted-foreground/70">{description}</p>
      )}
    </div>
  );
}
