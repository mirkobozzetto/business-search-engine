"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { NafCode } from "@/types/naf";

export function useNafSectionCodes(sectionCode: string) {
  return useQuery({
    queryKey: ["naf-section-codes", sectionCode],
    queryFn: () => fetchAPI<NafCode[]>(`/api/naf/section/${sectionCode}`),
    enabled: !!sectionCode,
  });
}
