"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { NafCode } from "@/types/naf";

export function useNafCode(code: string) {
  return useQuery({
    queryKey: ["naf-code", code],
    queryFn: () => fetchAPI<NafCode>(`/api/naf/code/${code}`),
    enabled: !!code,
  });
}
