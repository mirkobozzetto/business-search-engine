"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export function SearchBarHero() {
  const [query, setQuery] = useState("");
  const router = useRouter();

  function handleSearch() {
    const trimmed = query.trim();
    if (!trimmed) return;

    if (/^\d{9}$/.test(trimmed) || /^\d{14}$/.test(trimmed)) {
      router.push(`/company/${trimmed}`);
      return;
    }

    if (/^\d{2}\.\d{2}[A-Z]?$/.test(trimmed)) {
      router.push(`/naf/${trimmed}`);
      return;
    }

    router.push(`/search?denomination=${encodeURIComponent(trimmed)}`);
  }

  return (
    <div className="flex w-full max-w-2xl gap-2">
      <Input
        placeholder="Nom d'entreprise, SIREN (9 chiffres), SIRET (14 chiffres) ou code NAF (ex: 62.01Z)"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleSearch()}
        className="h-12 text-base"
      />
      <Button onClick={handleSearch} size="lg" className="h-12 px-8">
        Rechercher
      </Button>
    </div>
  );
}
