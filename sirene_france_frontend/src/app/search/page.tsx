"use client";

import { useState, useEffect, useCallback } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CompanySearchForm } from "@/components/company/company-search-form";
import { CompanyResultsTable } from "@/components/company/company-results-table";
import { ResultsMeta } from "@/components/shared/results-meta";
import { LoadingSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { PaginationControls } from "@/components/shared/pagination-controls";
import { useCompanyMultiSearch } from "@/hooks/use-company-multi-search";
import { DEFAULT_LIMIT } from "@/lib/constants";

export default function SearchPage() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const [denomination, setDenomination] = useState(searchParams.get("denomination") || "");
  const [nafCode, setNafCode] = useState(searchParams.get("naf_code") || "");
  const [codePostal, setCodePostal] = useState(searchParams.get("code_postal") || "");
  const [commune, setCommune] = useState(searchParams.get("commune") || "");
  const [etatAdministratif, setEtatAdministratif] = useState(searchParams.get("etat_administratif") || "");
  const [dateCreationFrom, setDateCreationFrom] = useState(searchParams.get("date_creation_from") || "");
  const [dateCreationTo, setDateCreationTo] = useState(searchParams.get("date_creation_to") || "");
  const [categorieJuridique, setCategorieJuridique] = useState(searchParams.get("categorie_juridique") || "");
  const [trancheEffectifs, setTrancheEffectifs] = useState(searchParams.get("tranche_effectifs") || "");
  const [page, setPage] = useState(1);
  const [searchTriggered, setSearchTriggered] = useState(
    !!(searchParams.get("denomination") || searchParams.get("naf_code") || searchParams.get("code_postal") || searchParams.get("commune") || searchParams.get("date_creation_from") || searchParams.get("date_creation_to") || searchParams.get("categorie_juridique") || searchParams.get("tranche_effectifs"))
  );

  const offset = (page - 1) * DEFAULT_LIMIT;

  const { data, isLoading } = useCompanyMultiSearch(
    searchTriggered
      ? {
          denomination: denomination || undefined,
          naf_code: nafCode || undefined,
          code_postal: codePostal || undefined,
          commune: commune || undefined,
          etat_administratif: etatAdministratif && etatAdministratif !== "all" ? etatAdministratif : undefined,
          date_creation_from: dateCreationFrom || undefined,
          date_creation_to: dateCreationTo || undefined,
          categorie_juridique: categorieJuridique && categorieJuridique !== "all" ? categorieJuridique : undefined,
          tranche_effectifs: trancheEffectifs && trancheEffectifs !== "all" ? trancheEffectifs : undefined,
          limit: DEFAULT_LIMIT,
          offset,
        }
      : {}
  );

  const companiesRaw = data?.data;
  const companies = companiesRaw?.results || [];
  const meta = companiesRaw?.meta;
  const totalPages = meta ? Math.ceil((meta.total || 0) / (meta.limit || DEFAULT_LIMIT)) : 1;

  const updateURL = useCallback(() => {
    const params = new URLSearchParams();
    if (denomination) params.set("denomination", denomination);
    if (nafCode) params.set("naf_code", nafCode);
    if (codePostal) params.set("code_postal", codePostal);
    if (commune) params.set("commune", commune);
    if (etatAdministratif && etatAdministratif !== "all") params.set("etat_administratif", etatAdministratif);
    if (dateCreationFrom) params.set("date_creation_from", dateCreationFrom);
    if (dateCreationTo) params.set("date_creation_to", dateCreationTo);
    if (categorieJuridique && categorieJuridique !== "all") params.set("categorie_juridique", categorieJuridique);
    if (trancheEffectifs && trancheEffectifs !== "all") params.set("tranche_effectifs", trancheEffectifs);
    router.push(`/search?${params.toString()}`, { scroll: false });
  }, [denomination, nafCode, codePostal, commune, etatAdministratif, dateCreationFrom, dateCreationTo, categorieJuridique, trancheEffectifs, router]);

  function handleSearch() {
    setPage(1);
    setSearchTriggered(true);
    updateURL();
  }

  function handleReset() {
    setDenomination("");
    setNafCode("");
    setCodePostal("");
    setCommune("");
    setEtatAdministratif("");
    setDateCreationFrom("");
    setDateCreationTo("");
    setCategorieJuridique("");
    setTrancheEffectifs("");
    setPage(1);
    setSearchTriggered(false);
    router.push("/search");
  }

  useEffect(() => {
    const hasCriteria = !!(searchParams.get("denomination") || searchParams.get("naf_code") || searchParams.get("code_postal") || searchParams.get("commune") || searchParams.get("date_creation_from") || searchParams.get("date_creation_to") || searchParams.get("categorie_juridique") || searchParams.get("tranche_effectifs"));
    if (hasCriteria) {
      setSearchTriggered(true);
    }
  }, [searchParams]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Recherche d'entreprises</h1>
        <p className="text-muted-foreground">
          Recherche multi-critères dans la base SIRENE
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Critères</CardTitle>
        </CardHeader>
        <CardContent>
          <CompanySearchForm
            denomination={denomination}
            nafCode={nafCode}
            codePostal={codePostal}
            commune={commune}
            etatAdministratif={etatAdministratif}
            dateCreationFrom={dateCreationFrom}
            dateCreationTo={dateCreationTo}
            categorieJuridique={categorieJuridique}
            trancheEffectifs={trancheEffectifs}
            onDenominationChange={setDenomination}
            onNafCodeChange={setNafCode}
            onCodePostalChange={setCodePostal}
            onCommuneChange={setCommune}
            onEtatAdministratifChange={setEtatAdministratif}
            onDateCreationFromChange={setDateCreationFrom}
            onDateCreationToChange={setDateCreationTo}
            onCategorieJuridiqueChange={setCategorieJuridique}
            onTrancheEffectifsChange={setTrancheEffectifs}
            onSearch={handleSearch}
            onReset={handleReset}
          />
        </CardContent>
      </Card>

      {searchTriggered && (
        <Card>
          <CardHeader>
            <CardTitle>Résultats</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {isLoading ? (
              <LoadingSkeleton rows={5} />
            ) : companies.length > 0 ? (
              <>
                <ResultsMeta meta={meta ? { total: meta.total, duration_ms: meta.duration_ms, page, pages: totalPages } : undefined} />
                <CompanyResultsTable companies={companies} />
                <PaginationControls
                  currentPage={page}
                  totalPages={totalPages}
                  onPageChange={setPage}
                />
              </>
            ) : (
              <EmptyState
                title="Aucune entreprise trouvée"
                description="Essayez de modifier vos critères de recherche"
              />
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
