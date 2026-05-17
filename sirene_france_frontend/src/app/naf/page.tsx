"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { NafSearchInput } from "@/components/naf/naf-search-input";
import { NafResultsTable } from "@/components/naf/naf-results-table";
import { NafSectionList } from "@/components/naf/naf-section-list";
import { ResultsMeta } from "@/components/shared/results-meta";
import { LoadingSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { PaginationControls } from "@/components/shared/pagination-controls";
import { useNafSearch } from "@/hooks/use-naf-search";
import { useNafSections } from "@/hooks/use-naf-sections";
import { useDebounce } from "@/hooks/use-debounce";
import { DEFAULT_LIMIT } from "@/lib/constants";

export default function NafPage() {
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const debouncedQuery = useDebounce(query, 300);
  const offset = (page - 1) * DEFAULT_LIMIT;

  const { data: searchData, isLoading: searchLoading } = useNafSearch(
    debouncedQuery,
    DEFAULT_LIMIT,
    offset
  );
  const { data: sectionsData, isLoading: sectionsLoading } = useNafSections();

  const codes = searchData?.data || [];
  const sections = sectionsData?.data || [];
  const meta = searchData?.meta;
  const totalPages = meta?.pages || 1;

  function handlePageChange(newPage: number) {
    setPage(newPage);
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Codes NAF</h1>
        <p className="text-muted-foreground">
          Recherchez parmi les 732 codes NAF ou parcourez par section
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recherche</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <NafSearchInput value={query} onChange={(v) => { setQuery(v); setPage(1); }} />
          {searchLoading && <LoadingSkeleton rows={3} />}
          {debouncedQuery.length >= 2 && !searchLoading && codes.length === 0 && (
            <EmptyState
              title="Aucun code NAF trouvé"
              description={`Aucun résultat pour "${debouncedQuery}"`}
            />
          )}
          {codes.length > 0 && (
            <>
              <ResultsMeta meta={meta} />
              <NafResultsTable codes={codes} />
              <PaginationControls
                currentPage={page}
                totalPages={totalPages}
                onPageChange={handlePageChange}
              />
            </>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Sections NAF</CardTitle>
        </CardHeader>
        <CardContent>
          {sectionsLoading ? (
            <LoadingSkeleton rows={5} />
          ) : (
            <NafSectionList sections={sections} />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
