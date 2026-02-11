"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { CompanyResult } from "@/types/company";

interface CompanyLookupResponse {
  criteria: Record<string, string>;
  results: CompanyResult[];
  meta: Record<string, number>;
}

export function useCompanyLookup(identifier: string) {
  return useQuery({
    queryKey: ["company-lookup", identifier],
    queryFn: () => fetchAPI<CompanyLookupResponse>(`/api/companies/lookup/${identifier}`),
    enabled: !!identifier,
  });
}
