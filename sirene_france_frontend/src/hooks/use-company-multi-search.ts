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
  date_creation_from?: string;
  date_creation_to?: string;
  categorie_juridique?: string;
  tranche_effectifs?: string;
  limit?: number;
  offset?: number;
}

export function useCompanyMultiSearch(params: MultiSearchParams) {
  const hasAnyCriteria = !!(
    params.denomination ||
    params.naf_code ||
    params.code_postal ||
    params.commune ||
    params.etat_administratif ||
    params.date_creation_from ||
    params.date_creation_to ||
    params.categorie_juridique ||
    params.tranche_effectifs
  );

  const searchParams: Record<string, string> = {};
  if (params.denomination) searchParams.denomination = params.denomination;
  if (params.naf_code) searchParams.naf = params.naf_code;
  if (params.code_postal) searchParams.codepostal = params.code_postal;
  if (params.commune) searchParams.commune = params.commune;
  if (params.etat_administratif) searchParams.etat = params.etat_administratif;
  if (params.date_creation_from) searchParams.from = params.date_creation_from;
  if (params.date_creation_to) searchParams.to = params.date_creation_to;
  if (params.categorie_juridique) searchParams.categorie_juridique = params.categorie_juridique;
  if (params.tranche_effectifs) searchParams.tranche_effectifs = params.tranche_effectifs;
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
