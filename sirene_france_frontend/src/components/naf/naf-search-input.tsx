"use client";

import { Input } from "@/components/ui/input";

interface NafSearchInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export function NafSearchInput({ value, onChange, placeholder }: NafSearchInputProps) {
  return (
    <Input
      type="text"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder || "Rechercher un code NAF par mot-clé..."}
      className="h-11"
    />
  );
}
