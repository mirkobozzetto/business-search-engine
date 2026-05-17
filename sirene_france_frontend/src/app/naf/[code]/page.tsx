"use client";

import { useState, use } from "react";
import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { CompanyResultsTable } from "@/components/company/company-results-table";
import { ResultsMeta } from "@/components/shared/results-meta";
import { LoadingSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { PaginationControls } from "@/components/shared/pagination-controls";
import { useNafCode } from "@/hooks/use-naf-code";
import { useCompanyMultiSearch } from "@/hooks/use-company-multi-search";
import { DEFAULT_LIMIT } from "@/lib/constants";

export default function NafDetailPage({ params }: { params: Promise<{ code: string }> }) {
  const { code } = use(params);
  const [page, setPage] = useState(1);
  const offset = (page - 1) * DEFAULT_LIMIT;

  const { data: nafData, isLoading: nafLoading } = useNafCode(code);
  const { data: companiesData, isLoading: companiesLoading } = useCompanyMultiSearch({
    naf_code: code,
    limit: DEFAULT_LIMIT,
    offset,
  });

  const naf = nafData?.data;
  const companiesRaw = companiesData?.data;
  const companies = companiesRaw?.results || [];
  const meta = companiesRaw?.meta;
  const totalPages = meta ? Math.ceil((meta.total || 0) / (meta.limit || DEFAULT_LIMIT)) : 1;

  return (
    <div className="space-y-6">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/naf">Codes NAF</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{code}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {nafLoading ? (
        <LoadingSkeleton rows={2} />
      ) : naf ? (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <Badge variant="outline" className="font-mono text-base">
                {naf.code}
              </Badge>
              <CardTitle>{naf.label}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Section{" "}
              <Link
                href={`/naf/section/${naf.section_code}`}
                className="text-primary hover:underline"
              >
                {naf.section_code} - {naf.section_label}
              </Link>
            </p>
          </CardContent>
        </Card>
      ) : (
        <EmptyState title="Code NAF introuvable" />
      )}

      <Card>
        <CardHeader>
          <CardTitle>Entreprises avec le code {code}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {companiesLoading ? (
            <LoadingSkeleton rows={5} />
          ) : companies.length > 0 ? (
            <>
              <ResultsMeta meta={meta ? { total: meta.total, duration_ms: meta.duration_ms, page: page, pages: totalPages } : undefined} />
              <CompanyResultsTable companies={companies} />
              <PaginationControls
                currentPage={page}
                totalPages={totalPages}
                onPageChange={setPage}
              />
            </>
          ) : (
            <EmptyState title="Aucune entreprise trouvée pour ce code NAF" />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
