"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { NafCode } from "@/types/naf";
import { DEFAULT_LIMIT } from "@/lib/constants";

export function useNafSearch(query: string, limit = DEFAULT_LIMIT, offset = 0) {
  return useQuery({
    queryKey: ["naf-search", query, limit, offset],
    queryFn: () =>
      fetchAPI<NafCode[]>("/api/naf/search", {
        q: query,
        limit: String(limit),
        offset: String(offset),
      }),
    enabled: query.length >= 2,
  });
}
