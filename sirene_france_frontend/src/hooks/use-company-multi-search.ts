"use client";

import { useQuery } from "@tanstack/react-query";
import { fetchAPI } from "@/lib/api-client";
import type { CompanyResult } from "@/types/company";
import { DEFAULT_LIMIT } from "@/lib/constants";

interface MultiSearchParams {
  denomination?: string;
  naf_code?: string;
  code_postal?: string;
  commune?: string;
  etat_administratif?: string;
  limit?: number;
  offset?: number;
}

export function useCompanyMultiSearch(params: MultiSearchParams) {
  const hasAnyCriteria = !!(
    params.denomination ||
    params.naf_code ||
    params.code_postal ||
    params.commune ||
    params.etat_administratif
  );

  const searchParams: Record<string, string> = {};
  if (params.denomination) searchParams.denomination = params.denomination;
  if (params.naf_code) searchParams.naf_code = params.naf_code;
  if (params.code_postal) searchParams.code_postal = params.code_postal;
  if (params.commune) searchParams.commune = params.commune;
  if (params.etat_administratif) searchParams.etat_administratif = params.etat_administratif;
  searchParams.limit = String(params.limit || DEFAULT_LIMIT);
  searchParams.offset = String(params.offset || 0);

  return useQuery({
    queryKey: ["company-multi-search", searchParams],
    queryFn: () =>
      fetchAPI<{ criteria: Record<string, string>; results: CompanyResult[]; meta: Record<string, number> }>(
        "/api/companies/search/multi",
        searchParams
      ),
    enabled: hasAnyCriteria,
  });
}
