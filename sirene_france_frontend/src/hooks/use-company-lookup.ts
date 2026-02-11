"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { CompanyResult } from "@/types/company";

export function useCompanyLookup(identifier: string) {
  return useQuery({
    queryKey: ["company-lookup", identifier],
    queryFn: () => fetchAPI<CompanyResult>(`/api/companies/lookup/${identifier}`),
    enabled: !!identifier,
  });
}
