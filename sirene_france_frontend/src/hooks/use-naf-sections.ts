"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { NafSection } from "@/types/naf";

export function useNafSections() {
  return useQuery({
    queryKey: ["naf-sections"],
    queryFn: () => fetchAPI<NafSection[]>("/api/naf/sections"),
  });
}
